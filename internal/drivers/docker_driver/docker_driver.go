package dockerDriver

import (
	"bufio"
	"citadel/internal/models"
	"citadel/internal/repositories"
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	caesar "github.com/caesar-rocks/core"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type DockerDriver struct {
	Client       *client.Client
	RegistryAuth string
	AppsRepo     *repositories.ApplicationsRepository
	DeplsRepo    *repositories.DeploymentsRepository
	ipv4         string
	ipv6         string
	minioClient  *minio.Client
	minioAdmin   *madmin.AdminClient
}

func NewDockerDriver(appsRepo *repositories.ApplicationsRepository, deplsRepo *repositories.DeploymentsRepository) *DockerDriver {
	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal(err)
	}

	registryAuth := getRegistryAuth()
	res, err := client.RegistryLogin(
		context.Background(),
		registry.AuthConfig{Auth: registryAuth},
	)
	if err != nil {
		log.Fatal(err)
	}
	if res.Status != "Login Succeeded" {
		log.Fatal("Failed to login to registry")
	}

	// Set up Minio-related stuff
	minioHost := strings.Replace(os.Getenv("MINIO_HOST"), "http://", "", 1)
	minioHost = strings.Replace(minioHost, "https://", "", 1)

	minioClient, err := minio.New(minioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	minioAdmin, err := madmin.New(minioHost, os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), false)
	if err != nil {
		log.Fatalln(err)
	}

	return &DockerDriver{
		Client:       client,
		RegistryAuth: registryAuth,
		DeplsRepo:    deplsRepo,
		minioClient:  minioClient,
		minioAdmin:   minioAdmin,
	}
}

func (d *DockerDriver) Init() error {
	if err := d.initializeSwarm(); err != nil {
		return err
	}
	d.watchEvents()
	d.setIPs()
	return nil
}

func (d *DockerDriver) CreateApplication(app models.Application) error {
	return nil
}

func (d *DockerDriver) DeleteApplication(app models.Application) error {
	if d.ContainerExists(app.ID) {
		if err := d.Client.ContainerRemove(context.Background(), app.ID, container.RemoveOptions{Force: true}); err != nil {
			return err
		}
	}

	return nil
}

func (d *DockerDriver) CreateCertificate(app models.Application, cert models.Certificate) ([]models.DnsEntry, error) {
	return []models.DnsEntry{
		{Hostname: cert.Domain, Type: "A", Value: d.ipv4},
		{Hostname: cert.Domain, Type: "AAAA", Value: d.ipv6},
	}, nil
}

func (d *DockerDriver) CheckDnsConfig(app models.Application, cert models.Certificate) (bool, error) {
	records, err := net.LookupIP(cert.Domain)
	if err != nil {
		return false, err
	}

	ipv4Found := false
	ipv6Found := false

	for _, ip := range records {
		if ip.To4() != nil {
			if ip.String() == d.ipv4 {
				ipv4Found = true
			}
		} else {
			if ip.String() == d.ipv6 {
				ipv6Found = true
			}
		}
	}

	if !ipv4Found || !ipv6Found {
		return false, errors.New("no IPv4 or IPv6 address found for the domain")
	}

	return true, nil
}

func (d *DockerDriver) DeleteCertificate(app models.Application, cert models.Certificate) error {
	return nil
}

func (d *DockerDriver) IgniteApplication(app models.Application, depl models.Deployment) error {
	if _, err := d.Client.ImagePull(
		context.Background(),
		os.Getenv("REGISTRY_HOST")+"/"+app.ID,
		image.PullOptions{RegistryAuth: d.RegistryAuth},
	); err != nil {
		return err
	}

	if d.ContainerExists(app.ID) {
		if err := d.Client.ContainerRemove(context.Background(), app.ID, container.RemoveOptions{Force: true}); err != nil {
			return err
		}
	}

	traefikRule := "Host(`" + app.Slug + "." + os.Getenv("WILDCARD_TRAEFIK_DOMAIN")
	for _, cert := range depl.Application.Certificates {
		traefikRule += "`, `www." + cert.Domain
	}
	traefikRule += "`)"

	if _, err := d.Client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: os.Getenv("REGISTRY_HOST") + "/" + app.ID,
			Labels: map[string]string{
				"deployment_id": depl.ID,

				"traefik.enable": "true",
				"traefik.http.routers." + app.ID + ".rule":                      "Host(`" + app.Slug + "." + os.Getenv("WILDCARD_TRAEFIK_DOMAIN") + "`)",
				"traefik.http.routers." + app.ID + ".entrypoints":               "websecure",
				"traefik.http.routers." + app.ID + ".tls.certresolver":          "myresolver",
				"traefik.http.services." + app.ID + ".loadbalancer.server.port": app.GetEnvVar("PORT", "3000"),
				"traefik.http.routers." + app.ID + ".tls":                       "true",
			},
		},
		&container.HostConfig{AutoRemove: true},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"traefik": {NetworkID: "traefik"},
			},
		},
		nil,
		app.ID,
	); err != nil {
		return err
	}

	if err := d.Client.ContainerStart(context.Background(), app.ID, container.StartOptions{}); err != nil {
		return err
	}

	return nil
}

func (d *DockerDriver) StreamLogs(ctx *caesar.Context, app models.Application) error {
	if !d.ContainerExists(app.ID) {
		return nil
	}

	logStream, err := d.Client.ContainerLogs(ctx.Context(), app.ID, container.LogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       "50",
	})
	if err != nil {
		return err
	}

	closed := ctx.Context().Done()

	scanner := bufio.NewScanner(logStream)
	for scanner.Scan() {
		select {
		case <-closed:
			return nil
		default:
			logLine := strings.TrimLeft(scanner.Text(), " ")
			if len(logLine) > 2 {
				logLine = logLine[8:]
			}
			if !ctx.WantsJSON() {
				logLine = "<pre>" + logLine + "</pre>"
			}

			if err := ctx.SendSSE("log", logLine); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

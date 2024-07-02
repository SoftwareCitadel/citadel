package tui

import (
	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/util"
	"fmt"
	"os"
	"strings"

	"github.com/alevinval/sse/pkg/eventsource"
	"github.com/charmbracelet/huh/spinner"
)

func MonitorHealtcheck(orgId string, appSlug string) {
	healthCheckStatus := "pending"

	_ = spinner.New().Title("Waiting for healthcheck...").Action(func() {
		url := api.RetrieveApiBaseUrl() + "/orgs/" + orgId + "/apps/" + appSlug + "/logs/stream?previous=false"

		token, err := util.RetrieveTokenFromConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		es, err := eventsource.New(url, eventsource.WithBearerTokenAuth(token))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for {
			select {
			case event := <-es.MessageEvents():
				if strings.Contains(event.Data, "Health check") {
					if strings.Contains(event.Data, "is now passing") {
						healthCheckStatus = "passing"
						return
					}

					if strings.Contains(event.Data, "has failed") {
						healthCheckStatus = "failed"
						return
					}
				}
			}
		}
	}).Run()

	if healthCheckStatus == "pending" {
		fmt.Println("🟡 Health check is still pending")
		os.Exit(1)
	}

	if healthCheckStatus == "passing" {
		fmt.Println("🟢 Health check is now passing.")
		os.Exit(0)
	}

	if healthCheckStatus == "failed" {
		fmt.Println("🔴 Health check has failed.")
		os.Exit(1)
	}
}

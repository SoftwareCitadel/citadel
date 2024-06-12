package main

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/auth"
	"citadel/cmd/citadel/tui"
	"citadel/cmd/citadel/util"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Run:   runDeploy,
	Short: "Deploy your project",
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to deploy a project.")
		fmt.Println("Please run `citadel auth login` to log in.")
		return
	}

	if !util.IsAlreadyInitialized() {
		fmt.Println("Software Citadel is not initialized. Please run `citadel init` to initialize it.")
		return
	}

	fmt.Println("Uploading...")

	tarball, err := util.MakeTarball()
	if err != nil {
		fmt.Println(err)
		return
	}

	appSlug, err := util.RetrieveAppSlugFromConfig()
	if err != nil {
		fmt.Println("Failed to retrieve application id")
		os.Exit(1)
	}

	releaseCmd, err := util.RetrieveReleaseCommandFromProjectConfig()
	if err != nil {
		fmt.Println("Failed to retrieve release command")
		os.Exit(1)
	}

	shouldMonitorHealtcheck, err := api.DeployFromTarball(tarball, appSlug, releaseCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tui.StreamBuildLogs(appSlug)

	if shouldMonitorHealtcheck {
		tui.MonitorHealtcheck(appSlug)
	}
}

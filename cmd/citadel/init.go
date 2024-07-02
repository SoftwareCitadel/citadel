package main

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/auth"
	"citadel/cmd/citadel/cli"
	"citadel/cmd/citadel/tui"
	"citadel/cmd/citadel/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Run:   runInit,
	Short: "Initialize a Software Citadel project",
}

func init() {
	initCmd.Flags().StringP("org-id", "o", "", "Organization ID to use for initialization (optional)")
	initCmd.Flags().StringP("app-slug", "a", "", "Application slug to use for initialization (optional)")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	orgId, _ := cmd.Flags().GetString("org-id")
	appSlug, _ := cmd.Flags().GetString("app-slug")

	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to initialize a project.\nPlease run `citadel auth login` to log in.")
		os.Exit(1)
	}

	if util.IsAlreadyInitialized() {
		if !cli.AskYesOrNo("Software Citadel is already initialized. Do you want to reinitialize it?") {
			return
		}
	}

	if orgId == "" {
		orgId = tui.SelectOrganization()
		if orgId == "" {
			orgId = tui.CreateOrganization()
		}
	}

	if appSlug == "" {
		appSlug = tui.SelectApplication(orgId)
		if appSlug == "" {
			appSlug = tui.CreateApplication(orgId)
		}
	}

	err := util.InitializeConfigFile(orgId, appSlug)
	if err != nil {
		fmt.Println("Failed to initialize Software Citadel project.")
		return
	}

	fmt.Println("Congratulations! Your app is now initialized.")
}

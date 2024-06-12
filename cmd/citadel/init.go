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
	initCmd.Flags().StringP("application-slug", "a", "", "Application slug to use for initialization (optional)")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	appSlug, _ := cmd.Flags().GetString("application-slug")

	if !auth.IsLoggedIn() {
		fmt.Println("You must be logged in to initialize a project.\nPlease run `citadel auth login` to log in.")
		os.Exit(1)
	}

	if util.IsAlreadyInitialized() {
		if !cli.AskYesOrNo("Software Citadel is already initialized. Do you want to reinitialize it?") {
			return
		}
	}

	if appSlug == "" {
		appSlug = tui.SelectApplication()
		if appSlug == "" {
			appSlug = tui.CreateApplication()
		}
	}

	err := util.InitializeConfigFile(appSlug)
	if err != nil {
		fmt.Println("Failed to initialize Software Citadel project.")
		return
	}

	fmt.Println("Congratulations! Your app is now initialized.")
}

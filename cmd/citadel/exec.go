package main

import (
	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/auth"
	"citadel/cmd/citadel/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Run:   runExec,
	Short: "Execute a command on the running container",
}

func runExec(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You are not logged in to Software Citadel.")
		os.Exit(1)
	}

	appSlug, err := util.RetrieveAppSlugFromConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) != 1 {
		fmt.Println("Please provide a command to execute.")
		os.Exit(1)
	}

	command := args[0]

	err = api.ExecuteCommand(appSlug, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

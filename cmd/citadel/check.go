package main

import (
	"fmt"

	"citadel/cmd/citadel/auth"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Run:   runCheck,
	Short: "Check if you are logged in to Software Citadel",
}

func runCheck(cmd *cobra.Command, args []string) {
	if auth.IsLoggedIn() {
		fmt.Println("You are logged in to Software Citadel!")
	} else {
		fmt.Println("You are not logged in to Software Citadel.")
	}
}

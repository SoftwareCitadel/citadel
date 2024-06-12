package main

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/auth"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Run:   runLogout,
	Short: "Logout locally from Software Citadel",
}

func runLogout(cmd *cobra.Command, args []string) {
	if !auth.IsLoggedIn() {
		fmt.Println("You are not logged in to Software Citadel.")
		os.Exit(1)
	}

	if err := auth.Logout(); err != nil {
		fmt.Println("Failed to logout")
		os.Exit(1)
	}

	fmt.Println("You have successfully logged out from Software Citadel!")
}

package main

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/auth"
	"citadel/cmd/citadel/tui"
	"citadel/cmd/citadel/util"

	bspinner "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Run:   runLogin,
	Short: "Login to Software Citadel",
}

func runLogin(cmd *cobra.Command, args []string) {
	token, _ := cmd.Flags().GetString("token")
	if token != "" {
		err := util.StoreJWTToken(token)
		if err != nil {
			fmt.Println("Whoops. There was an error while trying to store your authentication token.")
			os.Exit(1)
		}

		return
	}

	s := bspinner.New()
	s.Spinner = bspinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#0099FF"))
	interrupt, err := tui.Run(s, runLoginTUI)
	if err != nil {
		fmt.Println("Whoops. There was an error while trying to log you in.")
	} else {
		if !interrupt {
			fmt.Println("\n🚀 You are now logged in to Software Citadel!")
		}
	}
}

func runLoginTUI(msg tui.BasicSpinnerMessager) error {
	msg.SetStatus("Requesting authentication session")
	sessionId, err := auth.GetAuthenticationSessionId()
	if err != nil {
		return err
	}

	url := api.RetrieveApiBaseUrl() + "/auth/cli/" + sessionId
	fmt.Printf("\nOpening browser to %s\n", url)

	msg.SetStatus("Opening browser...")
	util.OpenInBrowser(api.RetrieveApiBaseUrl() + "/auth/cli/" + sessionId)

	msg.SetStatus("Waiting for authentication...")
	token, err := auth.WaitForLogin(sessionId)
	if err != nil {
		return err
	}

	msg.SetStatus("Storing authentication token...")
	if err = util.StoreJWTToken(token); err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/auth"
	"citadel/cmd/citadel/util"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "citadel",
	Short: "citadel - a simple CLI to interact with the Software Citadel PaaS",
	Long:  "citadel - a simple CLI to interact with the Software Citadel PaaS",
}

func execute(version string) {
	loginCmd.Flags().StringP("token", "t", "", "Authentication token")

	authCmd := &cobra.Command{
		Use: "auth",
	}
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(checkCmd)

	envCmd := &cobra.Command{
		Use: "env",
		PreRun: func(cmd *cobra.Command, args []string) {
			logged := auth.IsLoggedIn()
			if !logged {
				fmt.Println("You are not logged in. Please type `citadel auth login` to authenticate to the API.")
				os.Exit(1)
			}

			initialized := util.IsAlreadyInitialized()
			if !initialized {
				fmt.Println("This project is not initialized. Please type `citadel init` to set up your project locally.")
			}
		},
	}
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envLoadCmd)

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(MakeVersionCmd(version))

	rootCmd.AddCommand(execCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

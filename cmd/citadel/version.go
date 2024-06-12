package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeVersionCmd(version string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Software Citadel v-" + version)
		},
	}
	return versionCmd
}

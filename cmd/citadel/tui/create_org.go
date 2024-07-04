package tui

import (
	"citadel/cmd/citadel/api"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sveltinio/prompti/input"
)

func CreateOrganization() string {
	questionPrompt := &input.Config{
		Message:     "What's the name of your organization?",
		Placeholder: "acme",
	}

	applicationName, err := input.Run(questionPrompt)
	if err != nil {
		fmt.Println("An error occurred while creating the application.")
		os.Exit(1)
	}

	m := newChooseComputingSpecs()
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	choice := res.(SelectModel).Choice
	splittedChoice := strings.Split(choice.ID, "x-")
	cpu := splittedChoice[0] + "x"
	memory := splittedChoice[1]

	application, err := api.CreateOrganization(applicationName, cpu, memory)
	if err != nil {
		fmt.Println("\n🔴 " + err.Error())
		os.Exit(1)
	}

	return application.Slug
}

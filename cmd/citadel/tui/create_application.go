package tui

import (
	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/util"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sveltinio/prompti/input"
)

func CreateApplication(orgId string) string {
	questionPrompt := &input.Config{
		Message:      "What's the name of your application?",
		Placeholder:  "webapp",
		ErrorMsg:     "Please enter a valid application name",
		ValidateFunc: util.SlugValidateFunc,
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

	application, err := api.CreateApplication(applicationName, cpu, memory)
	if err != nil {
		fmt.Println("\n🔴 " + err.Error())
		os.Exit(1)
	}

	return application.Slug
}

func newChooseComputingSpecs() SelectModel {
	choices := []SelectChoice{
		{
			Name: "shared-cpu-1x 256MB",
			ID:   "shared-cpu-1x-256MB",
		},
		{
			Name: "shared-cpu-1x 512MB",
			ID:   "shared-cpu-1x-512MB",
		},
		{
			Name: "shared-cpu-1x 1GB",
			ID:   "shared-cpu-1x-1GB",
		},
		{
			Name: "shared-cpu-1x 2GB",
			ID:   "shared-cpu-1x-2GB",
		},
		{
			Name: "shared-cpu-2x 512MB",
			ID:   "shared-cpu-2x-512MB",
		},
		{
			Name: "shared-cpu-2x 1GB",
			ID:   "shared-cpu-2x-1GB",
		},
		{
			Name: "shared-cpu-2x 2GB",
			ID:   "shared-cpu-2x-2GB",
		},
		{
			Name: "shared-cpu-2x 4GB",
			ID:   "shared-cpu-2x-4GB",
		},
		{
			Name: "shared-cpu-4x 1GB",
			ID:   "shared-cpu-4x-1GB",
		},
		{
			Name: "shared-cpu-4x 2GB",
			ID:   "shared-cpu-4x-2GB",
		},
		{
			Name: "shared-cpu-4x 4GB",
			ID:   "shared-cpu-4x-4GB",
		},
		{
			Name: "shared-cpu-4x 8GB",
			ID:   "shared-cpu-4x-8GB",
		},
		{
			Name: "shared-cpu-8x 2GB",
			ID:   "shared-cpu-8x-2GB",
		},
		{
			Name: "shared-cpu-8x 4GB",
			ID:   "shared-cpu-8x-4GB",
		},
		{
			Name: "performance-cpu-1x 2GB",
			ID:   "performance-cpu-1x-2GB",
		},
		{
			Name: "performance-cpu-1x 4GB",
			ID:   "performance-cpu-1x-4GB",
		},
		{
			Name: "performance-cpu-1x 8GB",
			ID:   "performance-cpu-1x-8GB",
		},
		{
			Name: "performance-cpu-2x 4GB",
			ID:   "performance-cpu-2x-4GB",
		},
		{
			Name: "performance-cpu-2x 8GB",
			ID:   "performance-cpu-2x-8GB",
		},
		{
			Name: "performance-cpu-2x 16GB",
			ID:   "performance-cpu-2x-16GB",
		},
		{
			Name: "performance-cpu-4x 8GB",
			ID:   "performance-cpu-4x-8GB",
		},
		{
			Name: "performance-cpu-4x 16GB",
			ID:   "performance-cpu-4x-16GB",
		},
		{
			Name: "performance-cpu-4x 32GB",
			ID:   "performance-cpu-4x-32GB",
		},
		{
			Name: "performance-cpu-8x 16GB",
			ID:   "performance-cpu-8x-16GB",
		},
		{
			Name: "performance-cpu-8x 32GB",
			ID:   "performance-cpu-8x-32GB",
		},
		{
			Name: "performance-cpu-8x 64GB",
			ID:   "performance-cpu-8x-64GB",
		},
		{
			Name: "performance-cpu-16x 32GB",
			ID:   "performance-cpu-16x-32GB",
		},
		{
			Name: "performance-cpu-16x 64GB",
			ID:   "performance-cpu-16x-64GB",
		},
		{
			Name: "performance-cpu-16x 128GB",
			ID:   "performance-cpu-16x-128GB",
		},
	}

	return NewSelectModel("Which computing specs would you like to use?", choices)
}

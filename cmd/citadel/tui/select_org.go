package tui

import (
	"fmt"
	"os"

	"citadel/cmd/citadel/api"

	tea "github.com/charmbracelet/bubbletea"
)

func newChooseOrgPromptModel() SelectModel {
	orgs, err := api.RetrieveOrgs()
	if err != nil {
		fmt.Println("Failed to retrieve applications")
		os.Exit(1)
	}

	choices := []SelectChoice{}
	for _, org := range orgs {
		choices = append(choices, SelectChoice{
			Name: org.Name,
			ID:   org.ID,
			Slug: org.Slug,
		})
	}
	choices = append(choices, SelectChoice{
		Name: "Create a new organization",
		ID:   "",
		Slug: "",
	})

	return NewSelectModel("Which organization would you like to deploy to?", choices)
}

func SelectOrganization() string {
	m := newChooseOrgPromptModel()
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return res.(SelectModel).Choice.ID
}

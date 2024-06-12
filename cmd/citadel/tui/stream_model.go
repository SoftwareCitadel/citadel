package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	durationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	appStyle      = lipgloss.NewStyle().Margin(1, 0, 1, 0)
)

type StreamModelResultMsg struct {
	Title    string
	Subtitle string
}

func (r StreamModelResultMsg) String() string {
	return fmt.Sprintf("%s %s", durationStyle.Render(r.Subtitle), r.Title)
}

type StreamModel struct {
	spinner  spinner.Model
	results  []StreamModelResultMsg
	quitting bool
	Title    string
}

func NewStreamModel(Title string) StreamModel {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
	return StreamModel{
		spinner: s,
		results: make([]StreamModelResultMsg, numLastResults),
		Title:   Title,
	}
}

func (m StreamModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m StreamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case StreamModelResultMsg:
		m.results = append(m.results[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m StreamModel) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
	} else {
		s += m.spinner.View() + " " + m.Title
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

func (m StreamModel) Run(streamFn func(p *tea.Program)) {
	p := tea.NewProgram(m)

	// Simulate activity
	go func() {
		streamFn(p)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

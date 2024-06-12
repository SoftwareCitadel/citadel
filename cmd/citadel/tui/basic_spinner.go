package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type BasicSpinnerMessager struct {
	detailCh chan string
	statusCh chan string
}

func (m *BasicSpinnerMessager) SetStatus(v string) {
	m.statusCh <- v
}
func (m *BasicSpinnerMessager) SetDetail(v string) {
	m.detailCh <- v
}

// UserFunc executes with spinner running.
type UserFunc func(msg BasicSpinnerMessager) error

type BasicSpinnerModel struct {
	m spinner.Model

	detail string
	status string

	statusCh <-chan string
	detailCh <-chan string

	interrupt bool

	doneCh <-chan any
}

func (m BasicSpinnerModel) Init() tea.Cmd {
	return m.m.Tick
}

func (m BasicSpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.doneCh:
		return m, tea.Quit
	default:
	}

	select {
	case v := <-m.statusCh:
		m.status = v
	default:
	}

	select {
	case v := <-m.detailCh:
		m.detail = v
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.interrupt = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.m, cmd = m.m.Update(msg)
		return m, cmd
	}
}

func (m BasicSpinnerModel) View() string {
	select {
	case <-m.doneCh:
		return ""
	default:
	}

	if m.interrupt {
		return fmt.Sprintf("%s %s\n", m.m.View(), "Stoppping on interrupt")
	}

	output := fmt.Sprintf("%s %s\n", m.m.View(), m.status)
	if m.detail != "" {
		output += "\n" + m.detail
	}
	return output
}

func InitializeBasicSpinnerModel(sp spinner.Model) BasicSpinnerModel {
	return BasicSpinnerModel{
		m: sp,
	}
}

func Run(sp spinner.Model, f UserFunc, opts ...tea.ProgramOption) (bool, error) {
	statusCh := make(chan string, 1)
	detailCh := make(chan string, 1)
	msg := BasicSpinnerMessager{
		statusCh: statusCh,
		detailCh: detailCh,
	}

	doneCh := make(chan any)
	var result error
	go func() {
		result = f(msg)
		close(statusCh)
		close(detailCh)
		close(doneCh)
	}()
	s := BasicSpinnerModel{
		m:        sp,
		statusCh: statusCh,
		detailCh: detailCh,
		doneCh:   doneCh,
	}
	p := tea.NewProgram(s, opts...)
	m, err := p.Run()
	if err != nil {
		return false, err
	}
	return m.(BasicSpinnerModel).interrupt, result
}

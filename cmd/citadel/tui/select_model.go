package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight   = 14
	defaultWidth = 20
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type SelectChoice struct {
	ID   string
	Name string
	Slug string
}

type SelectModel struct {
	ListModel list.Model
	Choice    SelectChoice
	Quitting  bool
}

func (i SelectChoice) FilterValue() string {
	return i.Name
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectChoice)
	if !ok {
		return
	}

	str := i.Name
	if i.Slug != "" {
		str = fmt.Sprintf("%d. %s", index+1, i.Name)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(strs ...string) string {
			return selectedItemStyle.Render("> " + strs[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ListModel.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			c, ok := m.ListModel.SelectedItem().(SelectChoice)

			if ok {
				m.Choice = SelectChoice{ID: c.ID, Name: c.Name, Slug: c.Slug}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.ListModel, cmd = m.ListModel.Update(msg)
	return m, cmd
}

func (m SelectModel) View() string {
	if m.Choice != (SelectChoice{}) {
		return fmt.Sprintf("You chose %s", m.Choice)
	}
	if m.Quitting {
		return quitTextStyle.Render("Quitting...")
	}
	return "\n" + m.ListModel.View()
}

func NewSelectModel(title string, choices []SelectChoice) SelectModel {
	items := []list.Item{}
	for _, c := range choices {
		items = append(items, c)
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SelectModel{
		ListModel: l,
	}
}

func (m SelectModel) Run() (SelectChoice, error) {
	p := tea.NewProgram(m)
	res, err := p.Run()
	if err != nil {
		return SelectChoice{}, err
	}
	return res.(SelectModel).Choice, nil
}

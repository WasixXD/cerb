package main

import (
	"log"
	"os/user"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var manager Manager

type keymap struct {
	up      key.Binding
	down    key.Binding
	refresh key.Binding
	quit    key.Binding
}

type model struct {
	table  table.Model
	keymap keymap
	help   help.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.refresh):
			manager.cacheProc()
			manager.ParseTCP()
			rows := []table.Row{}
			for _, k := range manager.tcp {
				rows = append(rows, []string{k.pid, k.port, k.process_name, string(k.inode)})
			}

			m.table.SetRows(rows)
			m.table.Focus()
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd

}

func (m model) View() string {
	s := m.table.View()
	s += "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.refresh,
		m.keymap.quit,
	})

	return s
}

func isSudo() bool {
	current, err := user.Current()

	if err != nil {
		log.Fatalln("ERROR: getting the user", err)
	}

	return current.Username == "root"
}

func main() {
	if !isSudo() {
		log.Fatalf("[!] Run this program using sudo")
		return
	}

	manager = Manager{processes: make(map[string][]string), tcp: map[string]Process{}}
	manager.cacheProc()

	columns := []table.Column{
		{Title: "Pid", Width: 6},
		{Title: "Port", Width: 6},
		{Title: "Name", Width: 32},
		{Title: "Inode", Width: 7},
	}

	rows := []table.Row{}
	manager.ParseTCP()
	for _, k := range manager.tcp {
		rows = append(rows, []string{k.pid, k.port, k.process_name, string(k.inode)})
	}
	mod := model{
		table: table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(35),
		),
		keymap: keymap{
			up: key.NewBinding(
				key.WithKeys("↑"),
				key.WithHelp("↑", "up"),
			),
			down: key.NewBinding(
				key.WithKeys("↓"),
				key.WithHelp("↓", "down"),
			),
			refresh: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh"),
			),
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	p := tea.NewProgram(mod)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

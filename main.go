package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.table.Focus()
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd

}

func (m model) View() string {
	s := m.table.View()

	return s
}

func main() {
	columns := []table.Column{
		{Title: "Pid", Width: 6},
		{Title: "Port", Width: 6},
		{Title: "Name", Width: 32},
		{Title: "Inode", Width: 7},
	}
	rows := []table.Row{}
	for _, k := range ParseTCP() {
		port, _ := strconv.ParseInt(string(k.port), 16, 64)
		rows = append(rows, []string{fmt.Sprintf("%d", k.pid), fmt.Sprintf("%d", port), string(k.process_name), string(k.inode)})
	}
	mod := model{
		table: table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		),
	}
	p := tea.NewProgram(mod)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

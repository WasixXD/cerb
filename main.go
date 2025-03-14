package main

// type model struct {
// 	timer  stopwatch.Model
// 	cursor int
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "q":
// 			return m, tea.Quit
// 		case "a":
// 			return m, m.timer.Toggle()

// 		}
// 	}

// 	var cmd tea.Cmd
// 	m.timer, cmd = m.timer.Update(msg)

// 	return m, cmd

// }

// func (m model) View() string {
// 	s := m.timer.View()

// 	return s
// }

func main() {
	// mod := model{
	// 	timer: stopwatch.New(),
	// }
	// p := tea.NewProgram(mod)

	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	ParseTCP()

}

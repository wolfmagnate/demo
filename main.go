package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(model{slide: Init()}, tea.WithFPS(100)).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	slide *SlideModel
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tickMsg:
		m.slide = m.slide.Update()
		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.slide.View()
}

func tickCmd() tea.Cmd {
	return tea.Tick(10*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

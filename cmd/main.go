package main

import (
	"fmt"
	"github/iegpeppino/syspulse/logger"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize system stats error logger
	logger.SysDataLogger()

	// Initialize bubbletea model
	m := modelInit()

	// Run TUI in clean alternate terminal
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cpuTotalPercent float64
	// cpuLoads        cpu.TimesStat
	// processes       []systeminfo.ProcessInfo
	// memory          mem.VirtualMemoryStat
	// disk            []systeminfo.DiskInfo
	err error
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tickMsg:
		cpu, err := systeminfo.GetCPUPercent()
		m.cpuTotalPercent = cpu
		m.err = err
		return m, tick()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress 'q' to quit.", m.err)
	}
	return fmt.Sprintf("CPU Usage: %2f%%\nPress 'q' to quit.", m.cpuTotalPercent)
}

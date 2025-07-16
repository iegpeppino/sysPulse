package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	cpuColumns := []table.Column{
		{Title: "Load", Width: 30},
		{Title: "Value (%)", Width: 30},
		{Title: "Delta", Width: 20},
	}

	cpuTable := initTable(cpuColumns)

	memCols := []table.Column{
		{Title: "Type", Width: 40},
		{Title: "Value", Width: 40},
	}

	memTable := initTable(memCols)

	procCols := []table.Column{
		{Title: "PID", Width: 5},
		{Title: "Name", Width: 25},
		{Title: "Status", Width: 15},
		{Title: "Runtime", Width: 20},
		{Title: "Memory", Width: 10},
		{Title: "CPU", Width: 10},
	}

	procTable := initTable(procCols)

	diskCols := []table.Column{
		{Title: "Partition", Width: 25},
		{Title: "FsType", Width: 20},
		{Title: "Total", Width: 15},
		{Title: "Used", Width: 15},
		{Title: "Free", Width: 15},
	}

	diskTable := initTable(diskCols)

	m := model{
		tabs:      []string{"CPU", "MEMORY", "PROCESSES", "DISK"},
		ActiveTab: 0,
		cpuTable:  cpuTable,
		memTable:  memTable,
		procTable: procTable,
		diskTable: diskTable,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program", err)
		os.Exit(1)
	}
}

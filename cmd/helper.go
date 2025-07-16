package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Helper functions and structs

// Compares previous and actual number and returns symbol
func delta(now, prev float64) string {
	d := now - prev
	if d > 0 {
		return "↑"
	} else if d < 0 {
		return "↓"
	} else {
		return "="
	}
}

// Returns min of two int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Returns max of two int
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Used to construct the CPU and Memory load gauge bars
func loadGauge(loadPercent float64, width int) string {
	full := int(loadPercent / 100 * float64(width))
	empty := width - full

	bar := strings.Builder{}
	bar.WriteString(amber.Render(" | "))
	for i := 0; i < full; i++ {
		bar.WriteString(gaugeLoadStyle(loadPercent).Render("█ "))
	}
	for i := 0; i < empty; i++ {
		bar.WriteString(gray.Render("░ "))
	}
	bar.WriteString(amber.Render(" | "))
	return bar.String()
}

func (m model) renderTab(activeTab int) string {
	switch {
	case activeTab == 0:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf(
				"CPU: %.2f%%\n%s\n",
				m.cpuTotalPercent,
				loadGauge(m.cpuTotalPercent, 45)),
			baseStyle.Render(m.cpuTable.View()),
		)
	case activeTab == 1:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf(
				"RAM: %.2f%%\n%s\n",
				m.memory.UsedPercent,
				loadGauge(m.memory.UsedPercent, 45)),
			baseStyle.Render(m.memTable.View()),
		)
	case activeTab == 2:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			baseStyle.Render(m.procTable.View()),
		)
	case activeTab == 3:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			baseStyle.Render(m.diskTable.View()),
		)
	default:
		return fmt.Sprint(m.tabs)
	}
}

func convertBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func initTable(cols []table.Column) table.Model {
	t := table.New(
		table.WithFocused(false),
		table.WithHeight(20),
		table.WithColumns(cols),
		table.WithRows([]table.Row{}),
		table.WithStyles(TableStyle()),
	)
	return t
}

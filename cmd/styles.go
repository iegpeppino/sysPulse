package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func DefaulStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#FFBF00")
	return s
}

func TableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		Margin(2, 0, 0, 0).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFBF00")).
		BorderBottom(true).
		AlignVertical(lipgloss.Center).
		Bold(false)
	s.Cell = s.Cell.
		AlignHorizontal(lipgloss.Left).
		Padding(1, 0, 0, 1)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Bold(false)
	return s
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.Border{
		Top:    "---",
		Bottom: "---"},
	).
	BorderForeground(lipgloss.Color("#FFBF00")).
	Bold(true).
	Padding(1, 1, 1, 2).
	Margin(1, 1, 1, 2).
	AlignHorizontal(lipgloss.Center)

var tabStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FFBF00")).
	Bold(true).
	Padding(1, 1, 1, 1).
	Margin(0, 1, 1, 1).
	AlignHorizontal(lipgloss.Center).
	Height(30)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#139213ff")).Bold(false)
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1f155ff")).Bold(false)
	orange = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(false)
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#d62222ff")).Bold(false)
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	white  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FBFBFB")).Bold(true)
	amber  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBF00")).Bold(true)
)

func gaugeLoadStyle(cpuPercent float64) lipgloss.Style {
	switch {
	case cpuPercent < 50:
		return green
	case cpuPercent < 75:
		return yellow
	case cpuPercent < 90:
		return orange
	default:
		return red
	}
}

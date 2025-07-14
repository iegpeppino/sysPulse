package main

import (
	"encoding/json"
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/mem"
)

type model struct {
	height          int
	width           int
	cpuTotalPercent float64
	cpuLoads        map[string]interface{}
	// processes       []systeminfo.ProcessInfo
	memory mem.VirtualMemoryStat
	// disk            []systeminfo.DiskInfo
	styles *Styles
	err    error
}

type Styles struct {
	BorderColor lipgloss.Color
}

func DefaulStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("44")
	return s
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		cpuPercent, err := systeminfo.GetCPUPercent()
		m.cpuTotalPercent = cpuPercent
		m.err = err

		mem, err := systeminfo.GetMEMLoad()
		m.memory = mem
		m.err = err

		cpuStats, err := systeminfo.GetCPULoad()

		var cpuLoads map[string]interface{}
		err = json.Unmarshal([]byte(cpuStats.String()), &cpuLoads)
		if err != nil {
			fmt.Printf("Error unmarshalling CPU Loads: %v\n", err)
		}

		m.cpuLoads = cpuLoads
		m.err = err
		return m, tick()

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#139213ff")).Bold(false)
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1f155ff")).Bold(false)
	orange = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(false)
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#d62222ff")).Bold(false)
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	white  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FBFBFB")).Bold(true)
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

func loadGauge(loadPercent float64, width int) string {
	full := int(loadPercent / 100 * float64(width))
	empty := width - full

	bar := strings.Builder{}
	bar.WriteString(white.Render(" | "))
	for i := 0; i < full; i++ {
		bar.WriteString(gaugeLoadStyle(loadPercent).Render("█ "))
	}
	for i := 0; i < empty; i++ {
		bar.WriteString(gray.Render("░ "))
	}
	bar.WriteString(white.Render(" | "))
	return bar.String()
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading...\nPress 'q' to quit."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress 'q' to quit.", m.err)
	}

	cpuLoads := make([][]string, 0, len(m.cpuLoads))
	for key, value := range m.cpuLoads {
		if key != "cpu" {
			cpuLoads = append(cpuLoads, []string{key, fmt.Sprintf("%.2f%%", value.(float64))})
		}
	}
	// t := table.New().
	// 	Width(40).
	// 	Border(lipgloss.NormalBorder()).
	// 	BorderStyle(white).
	// 	Headers("Cpu Stats").
	// 	Rows(cpuLoads...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf(
			"CPU: %.2f%%\n%s\n",
			m.cpuTotalPercent,
			loadGauge(m.cpuTotalPercent, 40)),
		fmt.Sprintf(
			"RAM: %.2f%%\n%s\n",
			m.memory.UsedPercent,
			loadGauge(m.memory.UsedPercent, 40)),
		"Press 'q' to quit.",
	)
}

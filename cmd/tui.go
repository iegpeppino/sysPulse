package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type model struct {
	height          int
	width           int
	cpuTotalPercent float64
	cpuStats        cpu.TimesStat
	cpuPrevStats    cpu.TimesStat
	cpuTable        table.Model
	// processes       []systeminfo.ProcessInfo
	memory mem.VirtualMemoryStat
	// disk            []systeminfo.DiskInfo
	styles      *Styles
	err         error
	initialized bool
}

type cpuItem struct {
	title, desc string
}

type Styles struct {
	BorderColor lipgloss.Color
}

func DefaulStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("44")
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
		m.memory = *mem
		m.err = err

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("#FFBF00")).
			Bold(false)
		m.cpuTable.SetStyles(s)

		cpuTimes, _ := cpu.Times(false)
		if len(cpuTimes) > 0 {
			m.cpuPrevStats = m.cpuStats
			m.cpuStats = cpuTimes[0]
		}

		t := table.New(
			table.WithFocused(true),
			table.WithHeight(7),
			table.WithWidth(80),
		)
		m.cpuTable = t

		delta := func(now, prev float64) string {
			d := now - prev
			if d > 0 {
				return "↑"
			} else if d < 0 {
				return "↓"
			} else {
				return "="
			}
		}

		rows := []table.Row{
			{"User", fmt.Sprintf("%.2f%%", m.cpuStats.User), delta(m.cpuStats.User, m.cpuPrevStats.User)},
			{"System", fmt.Sprintf("%.2f%%", m.cpuStats.System), delta(m.cpuStats.System, m.cpuPrevStats.System)},
			{"Idle", fmt.Sprintf("%.2f%%", m.cpuStats.Idle), delta(m.cpuStats.Idle, m.cpuPrevStats.Idle)},
			{"Nice", fmt.Sprintf("%.2f%%", m.cpuStats.Nice), delta(m.cpuStats.Nice, m.cpuPrevStats.Nice)},
			{"Guest", fmt.Sprintf("%.2f%%", m.cpuStats.Guest), delta(m.cpuStats.Guest, m.cpuPrevStats.Guest)},
			{"IRQ", fmt.Sprintf("%.2f%%", m.cpuStats.Irq), delta(m.cpuStats.Irq, m.cpuPrevStats.Irq)},
			{"SoftIRQ", fmt.Sprintf("%.2f%%", m.cpuStats.Softirq), delta(m.cpuStats.Softirq, m.cpuPrevStats.Softirq)},
		}
		columns := []table.Column{
			{Title: "Load", Width: 30},
			{Title: "Value (%)", Width: 30},
			{Title: "Delta", Width: 20},
		}

		m.cpuTable.SetStyles(s)
		m.cpuTable.SetColumns(columns)
		m.cpuTable.SetRows(rows)
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

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf(
				"CPU: %.2f%%\n%s\n",
				m.cpuTotalPercent,
				loadGauge(m.cpuTotalPercent, 45)),
			// fmt.Sprintf(
			// 	"RAM: %.2f%%\n%s\n",
			// 	m.memory.UsedPercent,
			// 	loadGauge(m.memory.UsedPercent, 40)),
			baseStyle.Render(m.cpuTable.View()),
			"Press 'q' to quit.",
		),
	)
}

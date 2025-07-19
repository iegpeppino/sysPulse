package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type model struct {
	tabs            []string
	ActiveTab       int
	height          int
	width           int
	keys            keyMap
	help            help.Model
	cpuTotalPercent float64
	cpuStats        cpu.TimesStat
	cpuPrevStats    cpu.TimesStat
	cpuTable        table.Model
	processes       []systeminfo.ProcessInfo
	procTable       table.Model
	memory          mem.VirtualMemoryStat
	memTable        table.Model
	disk            []systeminfo.DiskInfo
	diskTable       table.Model
	err             error
}

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "a"),
		key.WithHelp("←/a", "switch tab to left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d"),
		key.WithHelp("→/d", "switch tab to right"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type tickMsg struct{}

// Setting the ticker for 500 milliseconds intervals
func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		// Setting width using terminal size
		fd := uintptr(os.Stdout.Fd())
		width, _, _ := term.GetSize(fd)
		m.width = width
		m.height = msg.Height
		m.help.Width = msg.Width

	case tickMsg:
		cpuPercent, err := systeminfo.GetCPUPercent()
		m.cpuTotalPercent = cpuPercent
		m.err = err

		mem, err := systeminfo.GetMEMLoad()
		m.memory = *mem
		m.err = err

		cpuTimes, _ := systeminfo.GetCPULoad()
		if len(cpuTimes) > 0 {
			m.cpuPrevStats = m.cpuStats
			m.cpuStats = cpuTimes[0]
		}

		processes, _ := systeminfo.GetProcessInfo(7)
		m.processes = processes

		disks, _ := systeminfo.GetDISKUse()
		m.disk = disks

		cpuRows := []table.Row{
			{"User", fmt.Sprintf("%.2f%%", m.cpuStats.User), delta(m.cpuStats.User, m.cpuPrevStats.User)},
			{"System", fmt.Sprintf("%.2f%%", m.cpuStats.System), delta(m.cpuStats.System, m.cpuPrevStats.System)},
			{"Idle", fmt.Sprintf("%.2f%%", m.cpuStats.Idle), delta(m.cpuStats.Idle, m.cpuPrevStats.Idle)},
			{"Nice", fmt.Sprintf("%.2f%%", m.cpuStats.Nice), delta(m.cpuStats.Nice, m.cpuPrevStats.Nice)},
			{"Guest", fmt.Sprintf("%.2f%%", m.cpuStats.Guest), delta(m.cpuStats.Guest, m.cpuPrevStats.Guest)},
			{"IRQ", fmt.Sprintf("%.2f%%", m.cpuStats.Irq), delta(m.cpuStats.Irq, m.cpuPrevStats.Irq)},
			{"SoftIRQ", fmt.Sprintf("%.2f%%", m.cpuStats.Softirq), delta(m.cpuStats.Softirq, m.cpuPrevStats.Softirq)},
		}

		m.cpuTable.SetRows(cpuRows)

		memRows := []table.Row{
			{"Total", fmt.Sprintf("%s", getByteMagnitude(m.memory.Total))},
			{"Used", fmt.Sprintf("%s", getByteMagnitude(m.memory.Used))},
			{"Available", fmt.Sprintf("%s", getByteMagnitude(m.memory.Available))},
			{"Free", fmt.Sprintf("%s", getByteMagnitude(m.memory.Free))},
			{"Buffers", fmt.Sprintf("%s", getByteMagnitude(m.memory.Buffers))},
			{"Cached", fmt.Sprintf("%s", getByteMagnitude(m.memory.Cached))},
		}

		m.memTable.SetRows(memRows)

		procRows := []table.Row{}
		for _, p := range processes {
			row := table.Row{
				fmt.Sprintf("%d", p.PID),
				p.Name,
				fmt.Sprint(p.Status),
				fmt.Sprint(p.Runtime),
				fmt.Sprintf("%s", getByteMagnitude(p.Memory)),
				fmt.Sprintf("%.2f%%", p.CPU),
			}
			procRows = append(procRows, row)
		}

		m.procTable.SetRows(procRows)

		diskRows := []table.Row{}
		for _, d := range m.disk {
			row := table.Row{
				fmt.Sprint(d.Partition.Mountpoint),
				d.Partition.Fstype,
				fmt.Sprintf("%s", getByteMagnitude(d.Total)),
				fmt.Sprintf("%s", getByteMagnitude(d.Used)),
				fmt.Sprintf("%s", getByteMagnitude(d.Free)),
			}
			diskRows = append(diskRows, row)
		}

		m.diskTable.SetRows(diskRows)

		return m, tick()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, cmd
		case key.Matches(msg, m.keys.Right):
			m.ActiveTab = min(m.ActiveTab+1, len(m.tabs)-1)
			return m, cmd
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading...\nPress 'q' to quit."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress 'q' to quit.", m.err)
	}

	page := strings.Builder{}

	tabRow := make([]string, len(m.tabs))
	for i, t := range m.tabs {
		if i == m.ActiveTab {
			tabRow[i] = activeTab.Render(t)
		} else {
			tabRow[i] = tab.Render(t)
		}
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabRow...,
	)

	tabGap.MaxWidth(m.width)
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	sep := tabGap.Render(strings.Repeat(" ", max(0, m.width)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	page.WriteString(row + "\n\n")

	baseStyle.MaxWidth(m.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		page.String(),
		m.renderTab(m.ActiveTab),
		fmt.Sprint(sep),
		baseStyle.Render(m.help.View(m.keys)),
	)

}

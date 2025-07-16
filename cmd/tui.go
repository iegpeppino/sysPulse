package main

import (
	"fmt"
	"github/iegpeppino/syspulse/systeminfo"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type model struct {
	tabs            []string
	ActiveTab       int
	height          int
	width           int
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
	styles          *Styles
	err             error
	initialized     bool
}

type cpuItem struct {
	title, desc string
}

type Styles struct {
	BorderColor lipgloss.Color
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
	var cmd tea.Cmd
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
			{"Total", fmt.Sprintf("%s", convertBytes(m.memory.Total))},
			{"Used", fmt.Sprintf("%s", convertBytes(m.memory.Used))},
			{"Available", fmt.Sprintf("%s", convertBytes(m.memory.Available))},
			{"Free", fmt.Sprintf("%s", convertBytes(m.memory.Free))},
			{"Buffers", fmt.Sprintf("%s", convertBytes(m.memory.Buffers))},
			{"Cached", fmt.Sprintf("%s", convertBytes(m.memory.Cached))},
		}

		m.memTable.SetRows(memRows)

		procRows := []table.Row{}
		for _, p := range processes {
			row := table.Row{
				fmt.Sprintf("%d", p.PID),
				p.Name,
				fmt.Sprint(p.Status),
				fmt.Sprint(p.Runtime),
				fmt.Sprintf("%s", convertBytes(uint64(p.Memory))),
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
				fmt.Sprintf("%s", convertBytes(d.Total)),
				fmt.Sprintf("%s", convertBytes(d.Used)),
				fmt.Sprintf("%s", convertBytes(d.Free)),
			}
			diskRows = append(diskRows, row)
		}

		m.diskTable.SetRows(diskRows)

		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "down", "s":
			m.cpuTable.MoveDown(1)
			return m, cmd
		case "right", "d":
			m.ActiveTab = min(m.ActiveTab+1, len(m.tabs)-1)
		case "left", "a":
			m.ActiveTab = max(m.ActiveTab-1, 0)
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

	return tabStyle.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderTab(m.ActiveTab),
			"Press 'q' to quit.",
		),
	))
}

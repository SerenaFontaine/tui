package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"tui"
	"tui/widget"
)

type tickMsg time.Time

func tickCmd() tui.Cmd {
	return func() tui.Msg {
		time.Sleep(100 * time.Millisecond)
		return tickMsg(time.Now())
	}
}

type app struct {
	tabs  *widget.Tabs
	focus *tui.FocusManager
	theme tui.Theme

	// Dashboard tab
	list      *widget.List
	sparkline *widget.Sparkline
	gauge     *widget.Gauge

	// Table tab
	table *widget.Table

	// Logs tab
	viewport *widget.Viewport

	// Tree tab
	tree *widget.Tree

	// Form tab
	form *widget.Form

	// State
	progress  *widget.Progress
	spinner   *widget.Spinner
	tickCount int
	cpuData   []float64
}

func newApp() *app {
	a := &app{
		tabs:  widget.NewTabs([]string{"Dashboard", "Table", "Logs", "Tree", "Form"}),
		focus: tui.NewFocusManager("tabs", "content"),
		theme: tui.NordTheme,

		list: widget.NewList([]string{
			"System Status: OK",
			"CPU Usage: 42%",
			"Memory: 8.2 GB / 16 GB",
			"Disk: 120 GB / 500 GB",
			"Network: 1.2 Gbps",
			"Uptime: 14 days",
			"Processes: 312",
			"Load Average: 1.24",
		}),

		table:     widget.NewTable([]string{"PID", "Name", "CPU%", "Mem%", "Status"}),
		progress:  widget.NewProgress(),
		gauge:     widget.NewGauge(),
		spinner:   widget.NewSpinner().SetLabel("Monitoring..."),
		sparkline: widget.NewSparkline(nil),

		tree: widget.NewTree(
			widget.NewTreeNode("System",
				widget.NewTreeNode("Hardware",
					widget.NewTreeNode("CPU: AMD Ryzen 9"),
					widget.NewTreeNode("RAM: 16 GB DDR5"),
					widget.NewTreeNode("GPU: RTX 4070"),
				),
				widget.NewTreeNode("Storage",
					widget.NewTreeNode("NVMe: 500 GB"),
					widget.NewTreeNode("HDD: 2 TB"),
				),
			),
			widget.NewTreeNode("Network",
				widget.NewTreeNode("Interfaces",
					widget.NewTreeNode("eth0: 1.2 Gbps"),
					widget.NewTreeNode("wlan0: disabled"),
					widget.NewTreeNode("lo: 127.0.0.1"),
				),
				widget.NewTreeNode("DNS",
					widget.NewTreeNode("Primary: 1.1.1.1"),
					widget.NewTreeNode("Secondary: 8.8.8.8"),
				),
			),
			widget.NewTreeNode("Services",
				widget.NewTreeNode("nginx: running"),
				widget.NewTreeNode("postgres: running"),
				widget.NewTreeNode("redis: running"),
			),
		),

		form: widget.NewForm(
			widget.NewFormField("Host", "192.168.1.1"),
			widget.NewFormField("Port", "8080"),
			widget.NewFormField("Username", "admin"),
			widget.NewFormField("Password", ""),
			widget.NewFormField("Database", "mydb"),
		),
	}

	a.table.SetRows([][]string{
		{"1234", "nginx", "2.1", "1.4", "running"},
		{"5678", "postgres", "15.3", "8.2", "running"},
		{"9012", "redis", "0.8", "2.1", "running"},
		{"3456", "node", "8.7", "4.5", "running"},
		{"7890", "python", "22.1", "12.8", "running"},
		{"1357", "go-server", "5.4", "3.2", "running"},
		{"2468", "cron", "0.1", "0.3", "sleeping"},
	})

	a.viewport = widget.NewViewport(
		"[INFO]  Application started successfully\n" +
			"[INFO]  Listening on :8080\n" +
			"[INFO]  Connected to database\n" +
			"[WARN]  High memory usage detected\n" +
			"[INFO]  Request processed in 42ms\n" +
			"[INFO]  Cache hit ratio: 94.2%\n" +
			"[ERROR] Connection timeout to upstream\n" +
			"[INFO]  Retrying connection...\n" +
			"[INFO]  Connection restored\n" +
			"[INFO]  Health check passed\n" +
			"[INFO]  New deployment detected\n" +
			"[INFO]  Rolling update in progress\n" +
			"[INFO]  Pod 1/3 ready\n" +
			"[INFO]  Pod 2/3 ready\n" +
			"[INFO]  Pod 3/3 ready\n" +
			"[INFO]  Deployment complete\n")

	a.progress.SetPercent(0.0)
	a.gauge.SetPercent(0.42)
	a.focus.Focus("content")

	return a
}

func (a *app) Init() tui.Cmd {
	return tui.Batch(tickCmd(), a.spinner.Tick())
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyCtrlC:
			return a, tui.QuitCmd()
		case tui.KeyTab:
			if msg.Alt {
				a.focus.Next()
			} else if a.focus.IsFocused("tabs") {
				a.tabs.Selected = (a.tabs.Selected + 1) % len(a.tabs.Titles)
			} else {
				// Delegate to content
				return a, a.updateContent(msg)
			}
			return a, nil
		case tui.KeyBacktab:
			if a.focus.IsFocused("tabs") {
				a.tabs.Selected--
				if a.tabs.Selected < 0 {
					a.tabs.Selected = len(a.tabs.Titles) - 1
				}
			} else {
				return a, a.updateContent(msg)
			}
			return a, nil
		}

		if a.focus.IsFocused("content") {
			return a, a.updateContent(msg)
		}

	case tickMsg:
		a.tickCount++
		p := float64(a.tickCount%200) / 200.0
		a.progress.SetPercent(p)

		// Update CPU sparkline data
		cpu := 30 + 20*math.Sin(float64(a.tickCount)*0.1) + rand.Float64()*10
		a.cpuData = append(a.cpuData, cpu)
		if len(a.cpuData) > 100 {
			a.cpuData = a.cpuData[1:]
		}
		a.sparkline.SetData(a.cpuData)
		a.gauge.SetPercent(cpu / 100.0)

		return a, tickCmd()

	case widget.SpinnerTickMsg:
		var cmd tui.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd
	}

	return a, nil
}

func (a *app) updateContent(msg tui.Msg) tui.Cmd {
	switch a.tabs.Selected {
	case 0:
		var cmd tui.Cmd
		a.list, cmd = a.list.Update(msg)
		return cmd
	case 1:
		var cmd tui.Cmd
		a.table, cmd = a.table.Update(msg)
		return cmd
	case 2:
		var cmd tui.Cmd
		a.viewport, cmd = a.viewport.Update(msg)
		return cmd
	case 3:
		var cmd tui.Cmd
		a.tree, cmd = a.tree.Update(msg)
		return cmd
	case 4:
		var cmd tui.Cmd
		a.form, cmd = a.form.Update(msg)
		return cmd
	}
	return nil
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
	regions := tui.VSplit(area,
		tui.Fixed(3), // tab bar
		tui.Flex(1),  // content
		tui.Fixed(3), // gauge
		tui.Fixed(1), // status
	)

	// Tab bar
	tabsFocused := a.focus.IsFocused("tabs")
	tabBlock := a.theme.Block("TUI Framework Demo", tabsFocused)
	a.tabs.SetBlock(tabBlock)
	a.tabs.ActiveStyle = a.theme.TitleStyle()
	a.tabs.InactiveStyle = a.theme.MutedStyle()
	a.tabs.Render(buf, regions[0])

	// Content
	switch a.tabs.Selected {
	case 0:
		a.renderDashboard(buf, regions[1])
	case 1:
		block := a.theme.Block("Processes", !tabsFocused)
		block.Border = tui.BorderRounded
		a.table.SetBlock(block)
		a.table.SetSelectedStyle(a.theme.SelectedStyle())
		a.table.HeaderStyle = a.theme.TitleStyle()
		a.table.Render(buf, regions[1])
	case 2:
		block := a.theme.Block("Logs", !tabsFocused)
		block.Border = tui.BorderRounded
		a.viewport.SetBlock(block)
		a.viewport.Render(buf, regions[1])
	case 3:
		block := a.theme.Block("System Tree", !tabsFocused)
		block.Border = tui.BorderRounded
		a.tree.SetBlock(block)
		a.tree.SetSelectedStyle(a.theme.SelectedStyle())
		a.tree.Render(buf, regions[1])
	case 4:
		block := a.theme.Block("Connection Settings", !tabsFocused)
		block.Border = tui.BorderRounded
		a.form.SetBlock(block)
		a.form.Render(buf, regions[1])
	}

	// Gauge
	gaugeBlock := a.theme.Block("CPU Usage", false)
	gaugeBlock.Border = tui.BorderRounded
	a.gauge.SetBlock(gaugeBlock)
	a.gauge.FilledStyle = tui.NewStyle().Bg(a.theme.Primary).Fg(a.theme.TextPrimary)
	a.gauge.Render(buf, regions[2])

	// Status bar
	status := fmt.Sprintf(" Alt+Tab: focus | Tab: switch | ↑↓: navigate | Ctrl+C: quit | %s | Tab %d/%d",
		a.spinner.View(), a.tabs.Selected+1, len(a.tabs.Titles))
	statusStyle := a.theme.StatusBarStyle()
	for x := regions[3].X; x < regions[3].Right(); x++ {
		buf.SetChar(x, regions[3].Y, ' ', statusStyle)
	}
	buf.SetString(regions[3].X, regions[3].Y, status, statusStyle)
}

func (a *app) renderDashboard(buf *tui.Buffer, area tui.Rect) {
	cols := tui.HSplit(area, tui.Percent(50), tui.Flex(1))

	// Left: list
	block := a.theme.Block("System Info", a.focus.IsFocused("content"))
	block.Border = tui.BorderRounded
	a.list.SetBlock(block)
	a.list.SetSelectedStyle(a.theme.SelectedStyle())
	a.list.Render(buf, cols[0])

	// Right: sparkline + details
	rightRows := tui.VSplit(cols[1], tui.Fixed(5), tui.Flex(1))

	// Sparkline
	sparkBlock := a.theme.Block("CPU History", false)
	sparkBlock.Border = tui.BorderRounded
	a.sparkline.SetBlock(sparkBlock)
	a.sparkline.SetStyle(tui.NewStyle().Fg(a.theme.Success))
	a.sparkline.Render(buf, rightRows[0])

	// Details
	detailBlock := a.theme.Block("Details", false)
	detailBlock.Border = tui.BorderRounded
	inner := detailBlock.Render(buf, rightRows[1])

	item := a.list.SelectedItem()
	if item != nil {
		// Rich styled text
		lines := tui.NewStyledText(
			tui.NewStyledLine(
				tui.BoldSpan("Selected: "),
				tui.ColorSpan(item.Text, a.theme.Secondary),
			),
			tui.NewStyledLine(),
			tui.NewStyledLine(
				tui.ColorSpan("Use ↑↓ to navigate", a.theme.TextMuted),
			),
			tui.NewStyledLine(
				tui.ColorSpan("Press Tab to switch views", a.theme.TextMuted),
			),
		)
		lines.Render(buf, inner)
	}
}

func main() {
	if err := tui.Run(newApp(), tui.WithTitle("TUI Demo")); err != nil {
		log.Fatal(err)
	}
}

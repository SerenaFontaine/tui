---
title: Demo Application
weight: 10
---

## Running the Demos

The `examples/` directory contains runnable demo applications:

```bash
# Simple list and input
go run ./examples/hello

# Full dashboard with tabs, sparklines, tree, form
go run ./examples/demo

# KGP image rendering
go run ./examples/image

# KGP animation
go run ./examples/animation
```

## Demo Application Structure

The full demo (`examples/demo/main.go`) showcases:

- **Tab navigation** with 5 tabs (Dashboard, Table, Logs, Tree, Form)
- **Focus management** with Alt+Tab to switch focus zones
- **Nord theme** with derived styles for borders, status bar, selection
- **Live sparkline** updated via periodic tick
- **Gauge widget** showing real-time CPU simulation
- **Animated spinner** in the status bar
- **Tree view** with expand/collapse
- **Form** with tab-based field navigation
- **Scrollable viewport** for log viewing

## Architecture Pattern

The demo follows a common pattern for multi-tab applications:

```go
type app struct {
    tabs   *widget.Tabs
    focus  *tui.FocusManager
    theme  tui.Theme
    // Per-tab state
    list   *widget.List
    table  *widget.Table
    // ...
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    // Global keys (tab switching, quit)
    // Then delegate to active tab's widget
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    rows := tui.VSplit(area, tui.Fixed(3), tui.Flex(1), tui.Fixed(1))
    a.tabs.Render(buf, rows[0])

    switch a.tabs.Selected {
    case 0: a.renderDashboard(buf, rows[1])
    case 1: a.table.Render(buf, rows[1])
    // ...
    }

    renderStatusBar(buf, rows[2])
}
```

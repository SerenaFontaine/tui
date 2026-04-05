---
title: Widgets
weight: 4
---

## Text

Static text display with alignment:

```go
t := widget.NewText("Hello, World!")
t.SetStyle(tui.NewStyle().Fg(tui.Green))
t.SetAlignment(widget.AlignCenter)
t.Render(buf, area)
```

## Input

Single-line text input with cursor:

```go
input := widget.NewInput("placeholder...")
input.Focused = true

// In Update:
input, cmd = input.Update(msg)

// In Render:
input.Render(buf, area)

// Read the value:
fmt.Println(input.Value)
```

Supports Home, End, Ctrl+U (clear left), Ctrl+K (clear right), Ctrl+W (delete word).

## List

Scrollable, selectable list with vim keys (j/k/g/G):

```go
list := widget.NewList([]string{"Alpha", "Bravo", "Charlie"})
list.SetSelectedStyle(tui.NewStyle().Bg(tui.Blue).Fg(tui.White))

// Navigate:
list, _ = list.Update(msg)

// Read selection:
item := list.SelectedItem()
```

## Table

Data table with headers, column widths, and selection:

```go
table := widget.NewTable([]string{"Name", "Status", "CPU%"})
table.SetRows([][]string{
    {"nginx", "running", "2.1"},
    {"postgres", "running", "15.3"},
})
table.SetColWidths([]int{20, 10, 8})
table.SetSelectedStyle(tui.NewStyle().Reverse(true))

table, _ = table.Update(msg)
table.Render(buf, area)
```

## Tabs

Tab bar with keyboard switching:

```go
tabs := widget.NewTabs([]string{"Dashboard", "Settings", "Logs"})

// Switch with Left/Right arrows or number keys (1-9):
tabs, _ = tabs.Update(msg)

// Render content based on active tab:
switch tabs.Selected {
case 0: renderDashboard(buf, area)
case 1: renderSettings(buf, area)
case 2: renderLogs(buf, area)
}
```

## Progress

Progress bar with percentage label:

```go
p := widget.NewProgress().SetPercent(0.75)
p.Render(buf, area)
```

## Gauge

Full-width gauge with centered label overlay:

```go
g := widget.NewGauge().SetPercent(0.42).SetLabel("CPU: 42%")
g.FilledStyle = tui.NewStyle().Bg(tui.Blue).Fg(tui.White)
g.Render(buf, area)
```

## Spinner

Animated loading indicator with 7 built-in styles:

```go
spinner := widget.NewSpinner().SetLabel("Loading...")
spinner.SetSpinnerStyle(widget.SpinnerDots) // or Line, Circle, Bounce, Meter, Globe, Block

// Start ticking in Init:
return spinner.Tick()

// Update on tick:
spinner, cmd = spinner.Update(msg)
```

## Viewport

Scrollable text viewer with mouse wheel support:

```go
vp := widget.NewViewport(longText)
vp, _ = vp.Update(msg) // handles Up/Down/PgUp/PgDn/mouse wheel
vp.Render(buf, area)
```

## Sparkline

Mini line chart:

```go
sl := widget.NewSparkline([]float64{0.2, 0.5, 0.8, 0.3, 0.9, 0.4})
sl.SetStyle(tui.NewStyle().Fg(tui.Green))
sl.Render(buf, area) // renders: ▂▄▇▃█▃
```

## Scrollbar

Vertical or horizontal scrollbar:

```go
sb := widget.NewScrollbar(totalLines, visibleLines, scrollOffset)
sb.Render(buf, tui.NewRect(area.Right()-1, area.Y, 1, area.Height))
```

## Tree

Navigable tree with expand/collapse:

```go
tree := widget.NewTree(
    widget.NewTreeNode("src",
        widget.NewTreeNode("main.go"),
        widget.NewTreeNode("app.go"),
    ),
    widget.NewTreeNode("docs"),
)
tree, _ = tree.Update(msg) // Left=collapse, Right=expand, Enter=toggle
```

## Dialog

Modal dialog with buttons:

```go
d := widget.NewDialog("Confirm", "Delete this item?")
d.SetButtons("Yes", "No")
d, _ = d.Update(msg) // Left/Right to switch buttons
d.Render(buf, area)  // renders centered overlay

selected := d.SelectedButton() // "Yes" or "No"
```

## Form

Multi-field labeled input form:

```go
form := widget.NewForm(
    widget.NewFormField("Host", "192.168.1.1"),
    widget.NewFormField("Port", "8080"),
    widget.NewFormField("User", "admin"),
)
form, _ = form.Update(msg) // Tab/Down to next field, Shift+Tab/Up to previous

values := form.Values() // map[string]string{"Host": "...", ...}
```

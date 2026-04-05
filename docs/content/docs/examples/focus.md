---
title: Focus and Input
weight: 8
---

## Focus Manager

Track focus across named components:

```go
focus := tui.NewFocusManager("sidebar", "content", "input")

focus.Next()               // advance to next
focus.Prev()               // go to previous
focus.Focus("input")       // focus by name
focus.IsFocused("sidebar") // check current
focus.Current()            // get current name
focus.Index()              // get current index
```

## Tab-Based Focus Switching

```go
func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    switch msg := msg.(type) {
    case tui.KeyMsg:
        if msg.Type == tui.KeyTab {
            a.focus.Next()
            return a, nil
        }
        if msg.Type == tui.KeyBacktab {
            a.focus.Prev()
            return a, nil
        }
    }

    // Delegate to focused component
    switch a.focus.Current() {
    case "sidebar":
        a.list, _ = a.list.Update(msg)
    case "content":
        a.viewport, _ = a.viewport.Update(msg)
    case "input":
        a.input, _ = a.input.Update(msg)
    }
    return a, nil
}
```

## Visual Focus Indication

Use the theme system to style focused borders:

```go
func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    block := a.theme.Block("Sidebar", a.focus.IsFocused("sidebar"))
    a.list.SetBlock(block)
    a.list.Render(buf, area)
}
```

## Form Widget

The Form widget has built-in focus management across fields:

```go
form := widget.NewForm(
    widget.NewFormField("Name", "Enter name"),
    widget.NewFormField("Email", "user@example.com"),
    widget.NewFormField("Password", ""),
)

// Tab/Down moves to next field, Shift+Tab/Up moves to previous
form, _ = form.Update(msg)

// Read values
name := form.Value("Name")
all := form.Values() // map[string]string
```

---
title: Hello World
weight: 1
---

## Minimal Application

The simplest TUI application renders text and quits on Ctrl+C:

```go
package main

import (
    "log"
    "tui"
)

type app struct{}

func (a *app) Init() tui.Cmd                              { return nil }
func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    if msg, ok := msg.(tui.KeyMsg); ok && msg.Type == tui.KeyCtrlC {
        return a, tui.QuitCmd()
    }
    return a, nil
}
func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    buf.SetString(0, 0, "Hello, World!", tui.NewStyle().Fg(tui.Green).Bold(true))
}

func main() {
    if err := tui.Run(&app{}); err != nil {
        log.Fatal(err)
    }
}
```

## With a List Widget

```go
package main

import (
    "log"
    "tui"
    "tui/widget"
)

type app struct {
    list *widget.List
}

func (a *app) Init() tui.Cmd { return nil }

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    switch msg := msg.(type) {
    case tui.KeyMsg:
        if msg.Type == tui.KeyCtrlC {
            return a, tui.QuitCmd()
        }
    }
    a.list, _ = a.list.Update(msg)
    return a, nil
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    block := tui.NewBlock()
    block.Title = "Items"
    block.Style = tui.NewStyle().Fg(tui.Cyan)
    a.list.SetBlock(block)
    a.list.SetSelectedStyle(tui.NewStyle().Bg(tui.Blue).Fg(tui.White))
    a.list.Render(buf, area)
}

func main() {
    items := []string{"First", "Second", "Third", "Fourth", "Fifth"}
    a := &app{list: widget.NewList(items)}
    if err := tui.Run(a); err != nil {
        log.Fatal(err)
    }
}
```

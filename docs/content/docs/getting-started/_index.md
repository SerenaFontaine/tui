---
title: Getting Started
weight: 10
---

## Installation

Add TUI to your Go module:

```bash
go get tui
```

## Requirements

- Go 1.23 or later
- Any terminal with ANSI escape sequence support
- For image features: a terminal with [Kitty Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/) support (Kitty, WezTerm, or Konsole)

## Quick Start

### Minimal Application

Every TUI application implements the `Component` interface with three methods:

```go
package main

import (
    "log"
    "tui"
)

type app struct {
    message string
}

func (a *app) Init() tui.Cmd { return nil }

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    switch msg := msg.(type) {
    case tui.KeyMsg:
        switch msg.Type {
        case tui.KeyCtrlC:
            return a, tui.QuitCmd()
        case tui.KeyRune:
            a.message += string(msg.Rune)
        case tui.KeyBackspace:
            if len(a.message) > 0 {
                a.message = a.message[:len(a.message)-1]
            }
        }
    }
    return a, nil
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    buf.SetString(0, 0, "Type something: "+a.message, tui.NewStyle().Fg(tui.Cyan))
    buf.SetString(0, 1, "Press Ctrl+C to quit", tui.NewStyle().Dim(true))
}

func main() {
    if err := tui.Run(&app{}); err != nil {
        log.Fatal(err)
    }
}
```

### With Widgets

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
    block.Title = "Pick one"
    a.list.SetBlock(block)
    a.list.Render(buf, area)
}

func main() {
    a := &app{list: widget.NewList([]string{"Alpha", "Bravo", "Charlie"})}
    if err := tui.Run(a); err != nil {
        log.Fatal(err)
    }
}
```

## The Elm Architecture

TUI follows the Elm Architecture pattern:

1. **Init** — Return an optional command to run on startup (fetch data, start a timer, etc.)
2. **Update** — Receive a message, update state, return an optional command
3. **Render** — Draw the current state into a buffer

Messages (`Msg`) flow in from keyboard, mouse, resize events, and async commands. The framework handles the event loop, terminal management, and efficient rendering automatically.

# TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/tui.svg)](https://pkg.go.dev/tui)

A fully featured terminal user interface framework for Go with first-class
[Kitty Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/)
support via [kgp](https://github.com/SerenaFontaine/kgp).

## Features

- **Elm Architecture** — Init/Update/Render loop with immutable message passing
- **Buffer-based rendering** — Diff-based updates for efficient terminal output
- **Kitty Graphics Protocol** — Full KGP integration: images, animation, all transmission methods
- **14 built-in widgets** — Text, Input, List, Table, Tabs, Viewport, Progress, Gauge, Spinner, Dialog, Sparkline, Scrollbar, Tree, Form
- **Flexible layout** — Constraint-based splitting: Fixed, Flex, Percent, Min, Max
- **True color** — ANSI 16, 256-color, and 24-bit RGB support
- **Theme system** — Built-in Nord, Gruvbox, Monochrome themes with derived styles
- **Rich text** — Inline styled spans for mixed formatting
- **Focus management** — Named focus tracking with tab navigation
- **Mouse support** — SGR mouse tracking with press, release, motion, and scroll
- **Suspend/Resume** — Ctrl+Z suspension with automatic state restore

## Installation

```bash
go get tui
```

## Quick Start

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
    block.Title = "My List"
    a.list.SetBlock(block)
    a.list.Render(buf, area)
}

func main() {
    a := &app{list: widget.NewList([]string{"One", "Two", "Three"})}
    if err := tui.Run(a); err != nil {
        log.Fatal(err)
    }
}
```

## Usage Examples

### Layout

Split areas into rows and columns with constraints:

```go
func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    // Vertical split: header, content, footer
    rows := tui.VSplit(area, tui.Fixed(3), tui.Flex(1), tui.Fixed(1))

    // Horizontal split: sidebar and main
    cols := tui.HSplit(rows[1], tui.Fixed(30), tui.Flex(1))
}
```

### Styling and Colors

```go
style := tui.NewStyle().
    Fg(tui.RGB(200, 100, 50)).
    Bg(tui.Hex("#1a1b26")).
    Bold(true).
    Italic(true)

buf.SetString(0, 0, "Hello", style)
```

### Borders

```go
block := tui.Block{
    Border: tui.BorderRounded,
    Title:  "Panel",
    Style:  tui.NewStyle().Fg(tui.Cyan),
}
inner := block.Render(buf, area)
// Render content inside inner
```

### Themes

```go
theme := tui.NordTheme

block := theme.Block("Title", focused)
statusStyle := theme.StatusBarStyle()
errorStyle := theme.ErrorStyle()
```

### Rich Styled Text

```go
line := tui.NewStyledLine(
    tui.BoldSpan("Status: "),
    tui.ColorSpan("OK", tui.Green),
)
line.Render(buf, 0, 0)
```

### Displaying Images (KGP)

```go
// From a Go image
imgWidget := widget.NewImage(myImage)
imgWidget.Render(buf, area)

// From a PNG file on disk
imgWidget := widget.NewImageFromFile("/path/to/image.png")

// From a pre-transmitted image ID
imgWidget := widget.NewImageFromID(imageID)

// With compression and z-index
imgWidget := widget.NewImageFromRGBA(pixels, 256, 256)
imgWidget.SetCompression(true).SetZIndex(-1)
```

### Image Management

```go
mgr := tui.NewImageManager()

// Transmit once, place many times
seq, id := mgr.Transmit(myImage)
screen.Flush(seq)

// Place at different positions and sizes
screen.Flush(mgr.Place(id, 0, 0, 20, 10))
screen.Flush(mgr.Place(id, 25, 0, 10, 5))

// Clean up
screen.Flush(mgr.Delete(id))
```

### Animation

```go
anim := tui.NewAnimation(imageID)
anim.AddImageFrame(frame1, 100) // 100ms gap
anim.AddImageFrame(frame2, 100)
anim.AddImageFrame(frame3, 100)

// Transmit frames
screen.Flush(anim.Encode())

// Playback control
screen.Flush(tui.PlayLoop(imageID))
screen.Flush(tui.StopAnimation(imageID))
screen.Flush(tui.GoToFrame(imageID, 2))
```

### Timers and Commands

```go
func (a *app) Init() tui.Cmd {
    return tui.TickCmd(100 * time.Millisecond)
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    switch msg.(type) {
    case tui.TickMsg:
        a.frame++
        return a, tui.TickCmd(100 * time.Millisecond)
    }
    return a, nil
}
```

### Focus Management

```go
focus := tui.NewFocusManager("sidebar", "content", "input")
focus.Next()                     // advance to next
focus.IsFocused("sidebar")       // check current
focus.Focus("input")             // focus by name
```

## Widgets

| Widget | Description |
|--------|-------------|
| `Text` | Static text with alignment and wrapping |
| `Input` | Single-line text input with cursor and scrolling |
| `List` | Scrollable selectable list with vim keys |
| `Table` | Data table with headers, column widths, selection |
| `Tabs` | Tab bar with keyboard switching |
| `Viewport` | Scrollable text viewer with mouse wheel |
| `Progress` | Progress bar with percentage label |
| `Gauge` | Full-width gauge with centered label overlay |
| `Spinner` | Animated indicator (7 styles: dots, line, circle, bounce, meter, globe, block) |
| `Dialog` | Modal dialog with buttons |
| `Sparkline` | Mini line chart using block characters |
| `Scrollbar` | Vertical or horizontal scrollbar |
| `Tree` | Navigable tree view with expand/collapse |
| `Form` | Multi-field labeled input form |
| `Image` | KGP image display (PNG, RGBA, RGB, file, virtual) |
| `AnimatedImage` | KGP animated image display |

## API Reference

### Core

| Type | Description |
|------|-------------|
| `Component` | Main interface: `Init()`, `Update(Msg)`, `Render(*Buffer, Rect)` |
| `App` | Application runner with event loop |
| `Buffer` | 2D cell grid with diff-based rendering |
| `Style` | Chainable text style (colors, bold, italic, etc.) |
| `Color` | Terminal color (ANSI, 256, RGB) |
| `Rect` | Rectangle with split, intersect, padding operations |
| `Block` | Border container with title |

### Layout

| Function | Description |
|----------|-------------|
| `VSplit(area, ...Constraint)` | Split vertically (top to bottom) |
| `HSplit(area, ...Constraint)` | Split horizontally (left to right) |
| `Fixed(n)` | Exact cell count |
| `Flex(weight)` | Weighted flexible space |
| `Percent(p)` | Percentage of total |
| `Min(n)` / `Max(n)` | Constrained sizing |

### KGP Image

| Function | Description |
|----------|-------------|
| `NewImagePlacement` | Create image placement with position and size |
| `TransmitImageWithID` | Transmit PNG image to terminal |
| `TransmitImageRGBAWithID` | Transmit RGBA with optional compression |
| `TransmitFileWithID` | Transmit from file path |
| `TransmitSharedMemWithID` | Transmit via POSIX shared memory |
| `DeleteImageByID` / `DeleteAllImages` | Remove images from terminal |
| `QueryKGPSupport` | Check terminal KGP support |
| `ParseImageResponse` | Parse terminal response |
| `NewImageManager` | Create image lifecycle manager |

### KGP Animation

| Function | Description |
|----------|-------------|
| `NewAnimation` | Create animation with frames |
| `PlayOnce` / `PlayLoop` / `PlayWithLoops` | Playback control |
| `StopAnimation` / `ResetAnimation` | Stop and reset |
| `GoToFrame` | Jump to specific frame |
| `ComposeFrames` | Composite frames together |

### Events

| Type | Description |
|------|-------------|
| `KeyMsg` | Keyboard event with key type, rune, alt modifier |
| `MouseMsg` | Mouse event with position, button, modifiers |
| `ResizeMsg` | Terminal resize with new dimensions |
| `TickMsg` | Timer tick |

## Terminal Support

TUI works in any terminal that supports ANSI escape sequences. KGP image features
require a terminal with Kitty Graphics Protocol support:

- **Kitty** (v0.19.0+) — Full support
- **WezTerm** — Partial support
- **Konsole** — Experimental support

Text-based widgets work in all terminals.

## Performance Tips

1. Use `TransmitImageWithID` to upload images once, then `NewImagePlacement().WithImageID()` for multiple placements
2. Enable ZLIB compression for large RGBA/RGB images with `WithCompression()`
3. Use `TransmitFromFile` or `TransmitFromSharedMem` for large images to avoid base64 overhead
4. The diff-based renderer only updates changed cells — avoid unnecessary full redraws

## License

[MIT](LICENSE)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## References

- [Kitty Graphics Protocol Specification](https://sw.kovidgoyal.net/kitty/graphics-protocol/)
- [KGP Go Bindings](https://github.com/SerenaFontaine/kgp)

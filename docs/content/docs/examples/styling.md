---
title: Styling
weight: 3
---

## Colors

TUI supports three color spaces:

```go
// ANSI 16 standard colors
tui.Red
tui.BrightBlue

// 256-color palette
tui.ANSI256(196)

// 24-bit true color
tui.RGB(255, 128, 0)
tui.Hex("#FF8000")
```

## Style Builder

Styles are value types with a chainable builder:

```go
style := tui.NewStyle().
    Fg(tui.RGB(200, 100, 50)).
    Bg(tui.Hex("#1a1b26")).
    Bold(true).
    Italic(true).
    Underline(true)

buf.SetString(0, 0, "styled text", style)
```

Available attributes: `Bold`, `Dim`, `Italic`, `Underline`, `Blink`, `Reverse`, `Strikethrough`.

## Borders

Five built-in border styles:

```go
tui.BorderSingle   // ┌─┐│└─┘
tui.BorderDouble   // ╔═╗║╚═╝
tui.BorderRounded  // ╭─╮│╰─╯
tui.BorderThick    // ┏━┓┃┗━┛
tui.BorderASCII    // +-+|+-+
```

Draw a bordered panel:

```go
block := tui.Block{
    Border: tui.BorderRounded,
    Title:  "Panel",
    Style:  tui.NewStyle().Fg(tui.Cyan),
}
inner := block.Render(buf, area)
// Render content inside `inner`
```

## Themes

Built-in themes provide consistent color schemes:

```go
theme := tui.NordTheme      // Nord palette
theme := tui.GruvboxTheme   // Gruvbox palette
theme := tui.MonochromeTheme // ANSI-only, maximum compatibility
theme := tui.DefaultTheme    // Blue accents

// Derived styles
block := theme.Block("Title", isFocused)
status := theme.StatusBarStyle()
err := theme.ErrorStyle()
selected := theme.SelectedStyle()
muted := theme.MutedStyle()
```

## Custom Theme

```go
myTheme := tui.Theme{
    Primary:          tui.RGB(100, 200, 255),
    Secondary:        tui.RGB(200, 150, 255),
    Accent:           tui.RGB(255, 100, 150),
    TextPrimary:      tui.White,
    TextMuted:        tui.BrightBlack,
    BorderColor:      tui.BrightBlack,
    BorderFocusColor: tui.RGB(100, 200, 255),
    Success:          tui.Green,
    Warning:          tui.Yellow,
    Error:            tui.Red,
    Info:             tui.Cyan,
}
```

---
title: Styled Text
weight: 7
---

## Spans

A `Span` is a run of text with a uniform style:

```go
span := tui.NewSpan("important", tui.NewStyle().Fg(tui.Red).Bold(true))
```

Convenience constructors:

```go
tui.PlainSpan("normal text")
tui.BoldSpan("bold text")
tui.ColorSpan("colored", tui.Green)
tui.StyledSpan("custom", myStyle)
```

## Styled Lines

Combine spans into a single line:

```go
line := tui.NewStyledLine(
    tui.BoldSpan("Status: "),
    tui.ColorSpan("OK", tui.Green),
    tui.PlainSpan(" | "),
    tui.ColorSpan("3 warnings", tui.Yellow),
)
line.Render(buf, 0, 0)
```

## Styled Text Blocks

Stack multiple lines:

```go
text := tui.NewStyledText(
    tui.NewStyledLine(
        tui.BoldSpan("Server: "),
        tui.ColorSpan("production", tui.Cyan),
    ),
    tui.NewStyledLine(
        tui.BoldSpan("Uptime: "),
        tui.PlainSpan("14 days"),
    ),
    tui.NewStyledLine(
        tui.BoldSpan("Load:   "),
        tui.ColorSpan("1.24", tui.Green),
    ),
)
text.Render(buf, area)
```

## Line Width

Calculate the display width of a styled line:

```go
line := tui.NewStyledLine(tui.PlainSpan("hello"), tui.BoldSpan(" world"))
width := line.Width() // 11
```

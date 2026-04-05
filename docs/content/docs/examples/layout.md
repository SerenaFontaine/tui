---
title: Layout
weight: 2
---

## Vertical Split

Split an area into rows with constraints:

```go
func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    rows := tui.VSplit(area,
        tui.Fixed(3),   // header: exactly 3 rows
        tui.Flex(1),    // content: fills remaining space
        tui.Fixed(1),   // status bar: exactly 1 row
    )

    renderHeader(buf, rows[0])
    renderContent(buf, rows[1])
    renderStatus(buf, rows[2])
}
```

## Horizontal Split

Split an area into columns:

```go
cols := tui.HSplit(area,
    tui.Fixed(30),   // sidebar: 30 columns wide
    tui.Flex(1),     // main: fills the rest
)
```

## Nested Layout

Combine splits for complex layouts:

```go
func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
    // Top-level: header, body, footer
    rows := tui.VSplit(area, tui.Fixed(3), tui.Flex(1), tui.Fixed(1))

    // Body: sidebar and content
    cols := tui.HSplit(rows[1], tui.Percent(30), tui.Flex(1))

    // Content: two equal panels
    panels := tui.HSplit(cols[1], tui.Flex(1), tui.Flex(1))
}
```

## Constraint Types

| Constraint | Description |
|------------|-------------|
| `Fixed(n)` | Exactly `n` cells |
| `Flex(weight)` | Proportional share of remaining space |
| `Percent(p)` | Percentage of total space |
| `Min(n)` | At least `n` cells, may grow |
| `Max(n)` | At most `n` cells, may shrink |

## Rect Splitting

`Rect` also provides direct splitting methods:

```go
left, right := area.SplitVertical(30)      // split at column 30
top, bottom := area.SplitHorizontal(10)    // split at row 10
inner := area.Inner(1)                      // 1-cell margin on all sides
inner := area.InnerPadding(1, 2, 1, 2)     // top, right, bottom, left
```

---
title: Buffer
weight: 2
---

`Buffer` is a 2D grid of cells representing terminal content. The framework creates buffers automatically; components write into them during `Render`.

## Creating Buffers

```go
buf := tui.NewBuffer(width, height)
```

Typically you receive a buffer from the framework rather than creating one.

## Writing Content

| Method | Description |
|--------|-------------|
| `SetChar(x, y, rune, Style)` | Set a single character |
| `SetString(x, y, string, Style) int` | Write a string, returns cells written |
| `SetStringInRect(string, Rect, Style) int` | Write with wrapping, returns lines used |
| `Fill(Rect, Cell)` | Fill area with a cell |
| `FillStyle(Rect, Style)` | Apply style to area, preserving characters |
| `DrawHLine(x, y, width, rune, Style)` | Draw horizontal line |
| `DrawVLine(x, y, height, rune, Style)` | Draw vertical line |

## Reading Content

| Method | Description |
|--------|-------------|
| `Get(x, y) Cell` | Read cell at position (emptyCell if out of bounds) |

## Buffer Operations

| Method | Description |
|--------|-------------|
| `Clear()` | Reset all cells to empty |
| `Resize(width, height)` | Change dimensions, discard content |
| `Merge(other, offsetX, offsetY)` | Overlay another buffer at offset |
| `AddImage(ImagePlacement)` | Register an image for KGP rendering |

## Diff-Based Rendering

The framework uses `Diff` internally to produce minimal terminal updates:

```go
output := current.Diff(previous)  // only changed cells
output := current.RenderFull()     // full redraw
```

## Cell

```go
type Cell struct {
    Char  rune
    Style Style
}
```

## Rect

```go
type Rect struct {
    X, Y, Width, Height int
}
```

| Method | Description |
|--------|-------------|
| `Right() int` | X + Width |
| `Bottom() int` | Y + Height |
| `Area() int` | Width * Height |
| `IsEmpty() bool` | Width or Height is zero |
| `Contains(x, y) bool` | Point inside rectangle |
| `Inner(margin) Rect` | Inset by margin on all sides |
| `InnerPadding(t, r, b, l) Rect` | Inset by individual sides |
| `SplitVertical(x) (Rect, Rect)` | Split into left and right |
| `SplitHorizontal(y) (Rect, Rect)` | Split into top and bottom |
| `Intersect(Rect) Rect` | Intersection of two rectangles |

## Block

`Block` draws a border and optional title:

```go
type Block struct {
    Border BorderStyle
    Title  string
    Style  Style
}
```

```go
block := tui.NewBlock()           // single-line border
block.Border = tui.BorderRounded  // change border style
block.Title = "Panel"
inner := block.Render(buf, area)  // returns inner area for content
```

Border styles: `BorderSingle`, `BorderDouble`, `BorderRounded`, `BorderThick`, `BorderASCII`, `BorderNone`.

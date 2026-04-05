---
title: Widgets
weight: 8
---

All widgets are in the `tui/widget` package. Interactive widgets follow the pattern:

```go
widget, cmd = widget.Update(msg)  // in Update
widget.Render(buf, area)           // in Render
```

Most widgets support `SetBlock(tui.Block)` for borders and titles.

## Text

Static text display with alignment.

```go
func NewText(content string) *Text
```

| Method | Description |
|--------|-------------|
| `SetStyle(Style)` | Text style |
| `SetAlignment(Alignment)` | `AlignLeft`, `AlignCenter`, `AlignRight` |
| `SetBlock(Block)` | Border |

## Input

Single-line text input with cursor.

```go
func NewInput(placeholder string) *Input
```

| Field | Description |
|-------|-------------|
| `Value` | Current text value |
| `Focused` | Whether input accepts keys |

| Method | Description |
|--------|-------------|
| `Update(Msg) (*Input, Cmd)` | Handle key events |
| `Focus() Cmd` | Return focus command |
| `SetBlock(Block)` | Border |
| `SetStyle(Style)` | Text style |

Keyboard: Home, End, Left, Right, Backspace, Delete, Ctrl+U, Ctrl+K, Ctrl+W.

## List

Scrollable, selectable list.

```go
func NewList(items []string) *List
func NewListFromItems(items []ListItem) *List
```

| Method | Description |
|--------|-------------|
| `Update(Msg) (*List, Cmd)` | Navigate with Up/Down/j/k/g/G/PgUp/PgDn |
| `SelectedItem() *ListItem` | Current selection |
| `SetSelectedStyle(Style)` | Selection highlight |
| `SetBlock(Block)` | Border |

## Table

Data table with headers.

```go
func NewTable(headers []string) *Table
```

| Method | Description |
|--------|-------------|
| `SetRows([][]string)` | Set table data |
| `SetColWidths([]int)` | Fixed column widths (0 = auto) |
| `Update(Msg) (*Table, Cmd)` | Navigate rows |
| `SelectedRow() []string` | Current row data |
| `SetSelectedStyle(Style)` | Selection highlight |
| `SetBlock(Block)` | Border |

## Tabs

Tab bar with keyboard switching.

```go
func NewTabs(titles []string) *Tabs
```

| Field | Description |
|-------|-------------|
| `Selected` | Active tab index |
| `ActiveStyle` | Active tab style |
| `InactiveStyle` | Inactive tab style |
| `Separator` | String between tabs |

Navigate with Left/Right arrows or number keys 1-9.

## Viewport

Scrollable text viewer.

```go
func NewViewport(content string) *Viewport
```

| Method | Description |
|--------|-------------|
| `SetContent(string)` | Update text |
| `ScrollTo(y)` | Jump to line |
| `LineCount() int` | Total wrapped lines |
| `Update(Msg)` | Up/Down/PgUp/PgDn/mouse wheel |

## Progress

Progress bar with label.

```go
func NewProgress() *Progress
```

| Field | Description |
|-------|-------------|
| `Percent` | 0.0 to 1.0 |
| `FilledChar` / `EmptyChar` | Bar characters |
| `FilledStyle` / `EmptyStyle` | Bar styles |
| `ShowLabel` | Show percentage label |

## Gauge

Full-width gauge with centered overlay label.

```go
func NewGauge() *Gauge
```

| Method | Description |
|--------|-------------|
| `SetPercent(float64)` | Set value |
| `SetLabel(string)` | Override label (empty = auto %) |

## Spinner

Animated loading indicator.

```go
func NewSpinner() *Spinner
```

| Method | Description |
|--------|-------------|
| `Tick() Cmd` | Start animation |
| `Update(Msg) (*Spinner, Cmd)` | Advance frame |
| `View() string` | Current frame string |
| `SetLabel(string)` | Text after spinner |
| `SetSpinnerStyle(SpinnerStyle)` | Animation style |

Styles: `SpinnerDots`, `SpinnerLine`, `SpinnerCircle`, `SpinnerBounce`, `SpinnerMeter`, `SpinnerGlobe`, `SpinnerBlock`.

## Dialog

Modal dialog with buttons.

```go
func NewDialog(title, message string) *Dialog
```

| Method | Description |
|--------|-------------|
| `SetButtons(...string)` | Set button labels |
| `SelectedButton() string` | Current button text |
| `Update(Msg)` | Left/Right to switch |

Renders centered over the area.

## Sparkline

Mini line chart using block characters.

```go
func NewSparkline(data []float64) *Sparkline
```

| Method | Description |
|--------|-------------|
| `SetData([]float64)` | Replace data |
| `PushData(val, maxLen)` | Append with cap |
| `SetMaxVal(float64)` | Fixed max (0 = auto) |

`SparklineGroup` stacks multiple sparklines vertically.

## Scrollbar

Vertical or horizontal scrollbar.

```go
func NewScrollbar(total, visible, offset int) *Scrollbar
func NewHScrollbar(total, visible, offset int) *Scrollbar
```

Hides automatically when `total <= visible`.

## Tree

Navigable tree with expand/collapse.

```go
func NewTree(roots ...*TreeNode) *Tree
func NewTreeNode(text string, children ...*TreeNode) *TreeNode
```

| Method | Description |
|--------|-------------|
| `SelectedNode() *TreeNode` | Current node |
| `Update(Msg)` | Up/Down/Left(collapse)/Right(expand)/Enter(toggle) |

Vim keys supported: j/k/h/l.

## Form

Multi-field labeled input form.

```go
func NewForm(fields ...FormField) *Form
func NewFormField(label, placeholder string) FormField
```

| Method | Description |
|--------|-------------|
| `Values() map[string]string` | All field values |
| `Value(label) string` | Single field value |
| `FocusedField() *FormField` | Current field |
| `Update(Msg)` | Tab/Down=next, Shift+Tab/Up=prev |

## Image

KGP image display.

```go
func NewImage(image.Image) *Image
func NewImageFromPNG([]byte) *Image
func NewImageFromRGBA([]byte, w, h) *Image
func NewImageFromRGB([]byte, w, h) *Image
func NewImageFromFile(path) *Image
func NewImageFromID(uint32) *Image
```

| Method | Description |
|--------|-------------|
| `SetCompression(bool)` | ZLIB compress |
| `SetZIndex(int)` | Layering |
| `SetCrop(x, y, w, h)` | Source rect |
| `SetVirtual(bool)` | Unicode placeholder |

## AnimatedImage

KGP animated image display.

```go
func NewAnimatedImage(anim *tui.Animation) *AnimatedImage
```

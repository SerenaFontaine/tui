---
title: Layout
weight: 4
---

The layout system splits rectangular areas according to constraints.

## Functions

### VSplit

```go
func VSplit(area Rect, constraints ...Constraint) []Rect
```

Splits vertically (top to bottom). Returns one `Rect` per constraint.

### HSplit

```go
func HSplit(area Rect, constraints ...Constraint) []Rect
```

Splits horizontally (left to right). Returns one `Rect` per constraint.

### Layout

```go
func Layout(area Rect, dir Direction, constraints []Constraint) []Rect
```

General-purpose split with explicit direction (`Horizontal` or `Vertical`).

## Constraints

| Constructor | Description |
|-------------|-------------|
| `Fixed(n)` | Exactly `n` cells |
| `Flex(weight)` | Proportional share of remaining space |
| `Percent(p)` | Percentage of total space |
| `Min(n)` | At least `n` cells |
| `Max(n)` | At most `n` cells |

## Resolution Order

1. `Fixed` and `Percent` are allocated first
2. `Max` is allocated up to its limit
3. `Min` is allocated at its minimum
4. Remaining space is distributed to `Flex` items by weight
5. Any rounding remainder goes to the last `Flex` item
6. `Min` constraints are satisfied by taking from `Flex` items if needed

## Examples

```go
// Header, content, footer
rows := tui.VSplit(area, tui.Fixed(3), tui.Flex(1), tui.Fixed(1))

// 30% sidebar, flexible main
cols := tui.HSplit(area, tui.Percent(30), tui.Flex(1))

// Three equal columns
cols := tui.HSplit(area, tui.Flex(1), tui.Flex(1), tui.Flex(1))

// 2:1 ratio
cols := tui.HSplit(area, tui.Flex(2), tui.Flex(1))
```

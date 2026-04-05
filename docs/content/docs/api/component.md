---
title: Component
weight: 1
---

`Component` is the core interface for TUI applications.

## Interface

```go
type Component interface {
    Init() Cmd
    Update(msg Msg) (Component, Cmd)
    Render(buf *Buffer, area Rect)
}
```

### Init

```go
func Init() Cmd
```

Called once when the application starts. Return a `Cmd` to run an initial command (start a timer, fetch data), or `nil` for no initial action.

### Update

```go
func Update(msg Msg) (Component, Cmd)
```

Called when a message arrives (key press, mouse event, timer tick, custom message). Update state and return the component and an optional command. Return a different `Component` to swap the root component entirely.

### Render

```go
func Render(buf *Buffer, area Rect)
```

Called after every `Update`. Draw the current state into the buffer within the given area. The framework handles diffing and flushing to the terminal.

## Running

```go
func Run(c Component, opts ...Option) error
```

Creates an `App` and runs the event loop. Blocks until quit.

```go
func NewApp(c Component, opts ...Option) *App
func (a *App) Run() error
```

For more control, create an `App` directly.

## Options

| Option | Description |
|--------|-------------|
| `WithAltScreen(bool)` | Enable alternate screen buffer (default: true) |
| `WithMouseEnabled(bool)` | Enable mouse event tracking (default: true) |
| `WithTitle(string)` | Set the terminal window title |

## Commands

```go
type Cmd func() Msg
```

Commands are functions that produce a message asynchronously. The framework runs them in goroutines and feeds the result back through `Update`.

| Function | Description |
|----------|-------------|
| `QuitCmd()` | Returns a `Cmd` that signals quit |
| `Batch(cmds ...Cmd)` | Combines commands to run concurrently |
| `ShowCursorCmd(x, y)` | Shows the terminal cursor at a position |
| `HideCursorCmd()` | Hides the terminal cursor |

## App Methods

| Method | Description |
|--------|-------------|
| `Send(msg Msg)` | Send an external message (goroutine-safe) |
| `SetCursor(x, y)` | Show cursor at position |
| `HideCursor()` | Hide cursor |
| `Images` | Built-in `ImageManager` instance |

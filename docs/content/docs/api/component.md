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
| `WithInput(io.Reader)` | Custom input reader (disables raw mode and signal handling) |
| `WithOutput(io.Writer)` | Custom output writer (disables raw mode and signal handling) |
| `WithSizeFunc(func() (int, int))` | Custom terminal size function |

### Custom I/O

By default, the app reads from stdin, writes to stdout, manages raw mode, and listens for OS signals (SIGWINCH, SIGCONT). When `WithInput` or `WithOutput` is provided, the screen becomes **unmanaged**: raw mode is skipped, signal handling is disabled, and Ctrl+Z suspension is passed through as a normal key event.

This is intended for embedding TUI applications in custom environments (e.g. an SSH server) where the caller controls the terminal lifecycle. Use `App.Send` to inject `ResizeMsg` from external sources:

```go
app := tui.NewApp(myComponent,
    tui.WithInput(session),
    tui.WithOutput(session),
    tui.WithSizeFunc(func() (int, int) {
        return ptyWidth, ptyHeight
    }),
)

// Push resize events from the SSH window-change callback
app.Send(tui.ResizeMsg{Width: newW, Height: newH})
```

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

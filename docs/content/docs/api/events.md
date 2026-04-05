---
title: Events
weight: 5
---

All events implement the `Msg` interface. The framework sends them to your `Update` method.

## Keyboard

```go
type KeyMsg struct {
    Type KeyType
    Rune rune   // only valid when Type == KeyRune
    Alt  bool
}
```

### Key Types

| Constant | Description |
|----------|-------------|
| `KeyRune` | Regular character (check `Rune` field) |
| `KeyEnter` | Enter/Return |
| `KeyTab` / `KeyBacktab` | Tab / Shift+Tab |
| `KeyBackspace` / `KeyDelete` | Backspace / Delete |
| `KeyEscape` | Escape |
| `KeySpace` | Space |
| `KeyUp` / `KeyDown` / `KeyLeft` / `KeyRight` | Arrow keys |
| `KeyHome` / `KeyEnd` | Home / End |
| `KeyPageUp` / `KeyPageDown` | Page Up / Page Down |
| `KeyInsert` | Insert |
| `KeyF1` through `KeyF12` | Function keys |
| `KeyCtrlA` through `KeyCtrlZ` | Ctrl+letter combinations |

### Pattern Matching

```go
switch msg := msg.(type) {
case tui.KeyMsg:
    switch msg.Type {
    case tui.KeyCtrlC:
        return a, tui.QuitCmd()
    case tui.KeyRune:
        switch msg.Rune {
        case 'q':
            return a, tui.QuitCmd()
        }
    }
}
```

## Mouse

```go
type MouseMsg struct {
    X, Y   int
    Button MouseButton
    Action MouseAction
    Alt, Ctrl, Shift bool
}
```

| Button | Description |
|--------|-------------|
| `MouseLeft` / `MouseMiddle` / `MouseRight` | Standard buttons |
| `MouseWheelUp` / `MouseWheelDown` | Scroll wheel |
| `MouseNone` | No button (motion event) |

| Action | Description |
|--------|-------------|
| `MousePress` | Button pressed |
| `MouseRelease` | Button released |
| `MouseMotion` | Mouse moved |

## Resize

```go
type ResizeMsg struct {
    Width, Height int
}
```

Sent when the terminal dimensions change.

## Lifecycle

| Message | Description |
|---------|-------------|
| `QuitMsg` | Application should exit |
| `FocusMsg` | Terminal gained focus |
| `BlurMsg` | Terminal lost focus |
| `SuspendMsg` | App is suspending (Ctrl+Z) |
| `ResumeMsg` | App resumed from suspension |

## Timer

| Message | Description |
|---------|-------------|
| `TickMsg` | Periodic tick (from `TickCmd`) |
| `AnimationTickMsg` | Animation tick (from `AnimateCmd`) |

## Focus

```go
type FocusChangeMsg struct {
    Previous string
    Current  string
}
```

## Cursor

```go
type CursorMsg struct {
    X, Y    int
    Visible bool
}
```

## Timer Commands

| Function | Description |
|----------|-------------|
| `TickCmd(duration)` | Send `TickMsg` after duration |
| `TickEvery(fps)` | Send `TickMsg` at FPS rate |
| `AfterCmd(duration, msg)` | Send custom msg after delay |
| `PeriodicCmd(interval, fn)` | Call fn at interval |
| `AnimateCmd(fps)` | Send `AnimationTickMsg` at FPS |

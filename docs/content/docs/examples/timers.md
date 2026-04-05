---
title: Timers
weight: 9
---

## Periodic Tick

Send a message at regular intervals:

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

## FPS-Based Tick

```go
// 30 FPS
return tui.TickEvery(30)
```

## Delayed Message

Send a message after a one-time delay:

```go
cmd := tui.AfterCmd(2*time.Second, myCustomMsg{})
```

## Custom Periodic Function

```go
cmd := tui.PeriodicCmd(5*time.Second, func(t time.Time) tui.Msg {
    return statusCheckMsg{time: t}
})
```

## Animation Tick

For KGP animation frame updates:

```go
func (a *app) Init() tui.Cmd {
    return tui.AnimateCmd(24) // 24 FPS
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
    switch msg.(type) {
    case tui.AnimationTickMsg:
        a.advanceFrame()
        return a, tui.AnimateCmd(24)
    }
    return a, nil
}
```

## Batching Commands

Run multiple commands concurrently:

```go
func (a *app) Init() tui.Cmd {
    return tui.Batch(
        tui.TickCmd(100 * time.Millisecond),
        a.spinner.Tick(),
        fetchDataCmd(),
    )
}
```

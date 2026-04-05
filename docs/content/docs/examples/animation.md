---
title: Animation
weight: 6
---

## Basic Animation Flow

1. Transmit the base image (frame 1)
2. Add animation frames with `NewAnimation`
3. Encode and send all frames
4. Control playback

## Building an Animation

```go
mgr := tui.NewImageManager()
imageID := mgr.NextID()

// Transmit base frame
seq := tui.TransmitImageWithID(frame0, imageID)
screen.Flush(seq)

// Build animation with additional frames
anim := tui.NewAnimation(imageID)
anim.AddImageFrame(frame1, 100) // 100ms gap
anim.AddImageFrame(frame2, 100)
anim.AddImageFrame(frame3, 100)
anim.DefaultGap = 100

// Transmit frames
screen.Flush(anim.Encode())
```

## Playback Control

```go
// Play once
screen.Flush(tui.PlayOnce(imageID))

// Loop infinitely
screen.Flush(tui.PlayLoop(imageID))

// Loop N times
screen.Flush(tui.PlayWithLoops(imageID, 5))

// Stop at current frame
screen.Flush(tui.StopAnimation(imageID))

// Reset to first frame
screen.Flush(tui.ResetAnimation(imageID))

// Jump to a specific frame
screen.Flush(tui.GoToFrame(imageID, 3))
```

## RGBA Frame Animation

Solid-color frames are useful for transitions and testing:

```go
anim := tui.NewAnimation(imageID)
anim.AddSolidFrame(80, 80, 255, 0, 0, 255, 200)   // red, 200ms
anim.AddSolidFrame(80, 80, 0, 255, 0, 255, 200)    // green
anim.AddSolidFrame(80, 80, 0, 0, 255, 255, 200)    // blue
```

## Frame Composition

Control how frames are blended:

```go
anim.SetFrameComposition(1, tui.CompositionBlend, 0)   // alpha blend
anim.SetFrameComposition(2, tui.CompositionReplace, 0)  // full replace
anim.SetFrameBackground(1, 0, 0, 0, 255)               // black background
```

## Compositing Frames

Composite a source frame onto a destination:

```go
screen.Flush(tui.ComposeFrames(
    imageID,
    2,                  // source frame
    1,                  // destination frame
    [4]int{0, 0, 40, 40},  // source rect
    [2]int{10, 10},         // dest offset
    tui.CompositionBlend,
))
```

## Animation Speed Control

Override the gap for all frames:

```go
screen.Flush(tui.SetAnimationGap(imageID, 50)) // 50ms between all frames
```

## Animated Image Widget

For rendering animated images within a TUI layout:

```go
anim := tui.NewAnimation(imageID)
anim.AddImageFrame(frame1, 100)
anim.AddImageFrame(frame2, 100)

animWidget := widget.NewAnimatedImage(anim)
animWidget.Render(buf, area)
```

## Cleanup

```go
screen.Flush(tui.DeleteAnimationFrames(imageID))     // remove frames only
screen.Flush(tui.DeleteAnimationFramesFree(imageID)) // remove and free memory
```

---
title: Animation
weight: 7
---

KGP animation support for building, transmitting, and controlling animated images.

## Animation Builder

```go
anim := tui.NewAnimation(imageID)
```

### Adding Frames

| Method | Description |
|--------|-------------|
| `AddFrame(data, format, w, h, gapMS)` | Raw pixel data |
| `AddPNGFrame(pngData, gapMS)` | PNG data |
| `AddImageFrame(image.Image, gapMS)` | Go image (converted to PNG) |
| `AddRGBAFrame(data, w, h, gapMS)` | Raw RGBA pixels |
| `AddRGBFrame(data, w, h, gapMS)` | Raw RGB pixels |
| `AddSolidFrame(w, h, r, g, b, a, gapMS)` | Solid color frame |

### Frame Options

| Method | Description |
|--------|-------------|
| `SetFrameComposition(idx, mode, bgFrame)` | Set composition mode |
| `SetFrameBackground(idx, r, g, b, a)` | Set background color |

### Fields

| Field | Description |
|-------|-------------|
| `ImageID` | Target image ID |
| `Frames` | Slice of `AnimationFrame` |
| `LoopCount` | 0=default, 1=infinite, N>1 loops N-1 times |
| `DefaultGap` | Default frame gap in ms |

### Encoding

```go
func (a *Animation) Encode() string
```

Returns escape sequences for all frames after the first (frame 1 is the base image, transmitted separately).

## Playback Control

| Function | Description |
|----------|-------------|
| `PlayOnce(imageID)` | Play animation once |
| `PlayLoop(imageID)` | Loop infinitely |
| `PlayWithLoops(imageID, count)` | Loop N times |
| `StopAnimation(imageID)` | Stop at current frame |
| `ResetAnimation(imageID)` | Reset to frame 1 |
| `GoToFrame(imageID, frameNum)` | Jump to specific frame |
| `SetAnimationState(imageID, state)` | Set playback state |
| `SetAnimationGap(imageID, gapMS)` | Override all frame gaps |

## Animation States

| Constant | Description |
|----------|-------------|
| `AnimationStop` | Stopped at current frame |
| `AnimationLoading` | Waiting for more frames |
| `AnimationLoop` | Playing with looping |

## Frame Composition

```go
func ComposeFrames(imageID, srcFrame, dstFrame uint32, srcRect [4]int, dstOffset [2]int, mode CompositionMode) string
```

Composite a source frame onto a destination frame.

| Mode | Description |
|------|-------------|
| `CompositionBlend` | Alpha blending (default) |
| `CompositionReplace` | Replace without blending |

## Cleanup

| Function | Description |
|----------|-------------|
| `DeleteAnimationFrames(imageID)` | Remove all frames |
| `DeleteAnimationFramesFree(imageID)` | Remove frames and free memory |

## Animation Tick

For frame-based animation loops in TUI applications:

```go
func AnimateCmd(fps int) Cmd
```

Sends `AnimationTickMsg` at the specified frame rate.

package tui

import (
	"image"
	"strings"
	"time"

	"github.com/SerenaFontaine/kgp"
)

// AnimationState controls playback.
type AnimationState = kgp.AnimationState

// Re-export animation states for convenience.
var (
	AnimationStop    = kgp.AnimationStop
	AnimationLoading = kgp.AnimationLoading
	AnimationLoop    = kgp.AnimationLoop
)

// CompositionMode controls how frames are blended.
type CompositionMode = kgp.CompositionMode

// Re-export composition modes.
var (
	CompositionBlend   = kgp.CompositionBlend
	CompositionReplace = kgp.CompositionReplace
)

// AnimationFrame represents a single frame in an animation.
type AnimationFrame struct {
	// Image data for this frame
	Data   []byte
	Format ImageFormat
	Width  int
	Height int

	// Frame timing
	GapMS uint32 // Delay in milliseconds before next frame

	// Composition
	Mode            CompositionMode
	BackgroundFrame uint32 // Frame number to use as base
	BackgroundColor uint32 // 32-bit RGBA background fill
}

// Animation manages a KGP image animation.
type Animation struct {
	ImageID    uint32
	Frames     []AnimationFrame
	LoopCount  uint32 // 0 = use default, 1 = infinite, N > 1 = loop N-1 times
	DefaultGap uint32 // Default gap in ms if frames don't specify one
}

// NewAnimation creates a new animation bound to an image ID.
func NewAnimation(imageID uint32) *Animation {
	return &Animation{ImageID: imageID}
}

// AddFrame adds a frame from raw pixel data.
func (a *Animation) AddFrame(data []byte, format ImageFormat, width, height int, gapMS uint32) *Animation {
	a.Frames = append(a.Frames, AnimationFrame{
		Data:   data,
		Format: format,
		Width:  width,
		Height: height,
		GapMS:  gapMS,
	})
	return a
}

// AddPNGFrame adds a frame from PNG data.
func (a *Animation) AddPNGFrame(pngData []byte, gapMS uint32) *Animation {
	return a.AddFrame(pngData, ImagePNG, 0, 0, gapMS)
}

// AddImageFrame adds a frame from a Go image.
func (a *Animation) AddImageFrame(img image.Image, gapMS uint32) *Animation {
	data, err := kgp.ImageToPNG(img)
	if err != nil {
		return a
	}
	return a.AddPNGFrame(data, gapMS)
}

// AddRGBAFrame adds a frame from raw RGBA data.
func (a *Animation) AddRGBAFrame(data []byte, width, height int, gapMS uint32) *Animation {
	return a.AddFrame(data, ImageRGBA, width, height, gapMS)
}

// AddRGBFrame adds a frame from raw RGB data.
func (a *Animation) AddRGBFrame(data []byte, width, height int, gapMS uint32) *Animation {
	return a.AddFrame(data, ImageRGB, width, height, gapMS)
}

// AddSolidFrame adds a solid-color frame (useful for transitions).
func (a *Animation) AddSolidFrame(width, height int, r, g, b, alpha uint8, gapMS uint32) *Animation {
	data := kgp.SolidColorImage(width, height, r, g, b, alpha)
	return a.AddRGBAFrame(data, width, height, gapMS)
}

// SetFrameComposition sets the composition mode for a specific frame.
func (a *Animation) SetFrameComposition(frameIdx int, mode CompositionMode, bgFrame uint32) *Animation {
	if frameIdx >= 0 && frameIdx < len(a.Frames) {
		a.Frames[frameIdx].Mode = mode
		a.Frames[frameIdx].BackgroundFrame = bgFrame
	}
	return a
}

// SetFrameBackground sets the background color for a specific frame.
func (a *Animation) SetFrameBackground(frameIdx int, r, g, b, alpha uint8) *Animation {
	if frameIdx >= 0 && frameIdx < len(a.Frames) {
		a.Frames[frameIdx].BackgroundColor = kgp.CreateRGBAColor(r, g, b, alpha)
	}
	return a
}

// Encode generates all the KGP escape sequences to transmit the animation.
// The first frame should be transmitted as a regular image first.
// This returns sequences for frames 2+ (frame number 2 onward).
func (a *Animation) Encode() string {
	if len(a.Frames) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, frame := range a.Frames {
		if i == 0 {
			// First frame is already transmitted as the base image
			continue
		}

		builder := kgp.NewFrame(a.ImageID).
			FrameNumber(uint32(i + 1))

		// Set frame data
		switch frame.Format {
		case ImagePNG:
			builder = builder.Format(kgp.FormatPNG).FrameData(frame.Data)
		case ImageRGBA:
			builder = builder.Format(kgp.FormatRGBA).
				Dimensions(frame.Width, frame.Height).
				FrameData(frame.Data)
		case ImageRGB:
			builder = builder.Format(kgp.FormatRGB).
				Dimensions(frame.Width, frame.Height).
				FrameData(frame.Data)
		}

		// Timing
		if frame.GapMS > 0 {
			builder = builder.Gap(frame.GapMS)
		} else if a.DefaultGap > 0 {
			builder = builder.Gap(a.DefaultGap)
		}

		// Composition
		if frame.Mode != 0 {
			builder = builder.Composition(frame.Mode)
		}
		if frame.BackgroundFrame > 0 {
			builder = builder.BackgroundFrame(frame.BackgroundFrame)
		}
		if frame.BackgroundColor > 0 {
			builder = builder.BackgroundColor(frame.BackgroundColor)
		}

		cmd := builder.Build()
		chunks := cmd.EncodeChunked(4096)
		for _, chunk := range chunks {
			sb.WriteString(chunk)
		}
	}

	return sb.String()
}

// ---------- Playback commands ----------

// PlayOnce returns the escape sequence to play the animation once.
func PlayOnce(imageID uint32) string {
	return kgp.PlayAnimation(imageID).Encode()
}

// PlayLoop returns the escape sequence to loop the animation forever.
func PlayLoop(imageID uint32) string {
	return kgp.PlayAnimationLoop(imageID).Encode()
}

// PlayWithLoops returns the escape sequence to loop the animation N times.
func PlayWithLoops(imageID uint32, count uint32) string {
	return kgp.PlayAnimationWithLoopCount(imageID, count).Encode()
}

// StopAnimation returns the escape sequence to stop playback.
func StopAnimation(imageID uint32) string {
	return kgp.StopAnimation(imageID).Encode()
}

// ResetAnimation returns the escape sequence to reset to frame 1.
func ResetAnimation(imageID uint32) string {
	return kgp.ResetAnimation(imageID).Encode()
}

// SetAnimationState sets the animation to a specific state.
func SetAnimationState(imageID uint32, state AnimationState) string {
	return kgp.NewAnimate(imageID).State(state).Build().Encode()
}

// SetAnimationGap overrides the gap for all frames.
func SetAnimationGap(imageID uint32, gapMS uint32) string {
	return kgp.NewAnimate(imageID).GapOverride(gapMS).Build().Encode()
}

// GoToFrame jumps to a specific frame number.
func GoToFrame(imageID uint32, frameNum uint32) string {
	return kgp.NewAnimate(imageID).
		State(kgp.AnimationStop).
		Frame(frameNum).
		Build().Encode()
}

// ---------- Frame composition ----------

// ComposeFrames composites a source frame onto a destination frame.
func ComposeFrames(imageID uint32, srcFrame, dstFrame uint32, srcRect [4]int, dstOffset [2]int, mode CompositionMode) string {
	builder := kgp.NewCompose(imageID).
		SourceFrame(srcFrame).
		DestFrame(dstFrame).
		Composition(mode)

	if srcRect[2] > 0 && srcRect[3] > 0 {
		builder = builder.SourceRect(srcRect[0], srcRect[1], srcRect[2], srcRect[3])
	}
	if dstOffset[0] > 0 || dstOffset[1] > 0 {
		builder = builder.DestOffset(dstOffset[0], dstOffset[1])
	}

	return builder.Build().Encode()
}

// ---------- Animation-related messages ----------

// AnimationTickMsg is sent by the animation ticker.
type AnimationTickMsg struct {
	Time time.Time
}

// AnimateCmd creates a command that ticks at the given FPS for animations.
func AnimateCmd(fps int) Cmd {
	interval := time.Second / time.Duration(fps)
	return func() Msg {
		time.Sleep(interval)
		return AnimationTickMsg{Time: time.Now()}
	}
}

// ---------- Delete animation frames ----------

// DeleteAnimationFrames removes all frames from an animation.
func DeleteAnimationFrames(imageID uint32) string {
	return kgp.NewDelete(kgp.DeleteFrames).ImageID(imageID).Build().Encode()
}

// DeleteAnimationFramesFree removes all frames and frees memory.
func DeleteAnimationFramesFree(imageID uint32) string {
	return kgp.NewDelete(kgp.DeleteFramesFree).ImageID(imageID).Build().Encode()
}

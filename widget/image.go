package widget

import (
	"image"
	"github.com/SerenaFontaine/tui"
)

// Image displays an image using the Kitty Graphics Protocol.
type Image struct {
	img        image.Image
	pngData    []byte
	rgbaData   []byte
	rgbData    []byte
	imgW, imgH int
	Block      *tui.Block

	// Options
	imageID                    uint32
	compress                   bool
	zIndex                     int
	cropX, cropY, cropW, cropH int
	filePath                   string
	virtual                    bool
}

// NewImage creates an image widget from a Go image.
func NewImage(img image.Image) *Image {
	return &Image{img: img}
}

// NewImageFromPNG creates an image widget from raw PNG data.
func NewImageFromPNG(data []byte) *Image {
	return &Image{pngData: data}
}

// NewImageFromRGBA creates an image widget from raw RGBA pixel data.
func NewImageFromRGBA(data []byte, width, height int) *Image {
	return &Image{rgbaData: data, imgW: width, imgH: height}
}

// NewImageFromRGB creates an image widget from raw RGB pixel data.
func NewImageFromRGB(data []byte, width, height int) *Image {
	return &Image{rgbData: data, imgW: width, imgH: height}
}

// NewImageFromFile creates an image widget that loads from a file path.
func NewImageFromFile(path string) *Image {
	return &Image{filePath: path}
}

// NewImageFromID creates an image widget for a previously transmitted image.
func NewImageFromID(imageID uint32) *Image {
	return &Image{imageID: imageID}
}

// SetBlock adds a border block.
func (i *Image) SetBlock(b tui.Block) *Image { i.Block = &b; return i }

// SetCompression enables ZLIB compression for RGBA/RGB data.
func (i *Image) SetCompression(enabled bool) *Image { i.compress = enabled; return i }

// SetZIndex sets the z-index for layering.
func (i *Image) SetZIndex(z int) *Image { i.zIndex = z; return i }

// SetCrop sets a source rectangle for cropping.
func (i *Image) SetCrop(x, y, w, h int) *Image {
	i.cropX = x
	i.cropY = y
	i.cropW = w
	i.cropH = h
	return i
}

// SetVirtual enables virtual (Unicode placeholder) placement.
func (i *Image) SetVirtual(v bool) *Image { i.virtual = v; return i }

// Render draws the image into the buffer area.
func (i *Image) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if i.Block != nil {
		inner = i.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	placement := tui.NewImagePlacement(inner.X, inner.Y, inner.Width, inner.Height)

	// Set image data source
	switch {
	case i.imageID > 0:
		placement = placement.WithImageID(i.imageID)
	case i.pngData != nil:
		placement = placement.WithPNG(i.pngData)
	case i.rgbaData != nil:
		placement = placement.WithRGBA(i.rgbaData, i.imgW, i.imgH)
	case i.rgbData != nil:
		placement = placement.WithRGB(i.rgbData, i.imgW, i.imgH)
	case i.filePath != "":
		placement = placement.WithFile(i.filePath, tui.ImagePNG)
	case i.img != nil:
		placement = placement.WithImage(i.img)
	default:
		return
	}

	// Apply options
	if i.compress {
		placement = placement.WithCompression()
	}
	if i.zIndex != 0 {
		placement = placement.WithZIndex(i.zIndex)
	}
	if i.cropW > 0 && i.cropH > 0 {
		placement = placement.WithCrop(i.cropX, i.cropY, i.cropW, i.cropH)
	}
	if i.virtual {
		placement = placement.WithVirtual()
	}

	buf.AddImage(placement)
}

// AnimatedImage displays an animated image using KGP animation frames.
type AnimatedImage struct {
	Animation *tui.Animation
	Block     *tui.Block
	imageID   uint32
	built     bool
	encoded   string
}

// NewAnimatedImage creates an animated image widget. The base image is
// transmitted as frame 1, and the Animation contains subsequent frames.
func NewAnimatedImage(anim *tui.Animation) *AnimatedImage {
	return &AnimatedImage{
		Animation: anim,
		imageID:   anim.ImageID,
	}
}

// SetBlock adds a border block.
func (a *AnimatedImage) SetBlock(b tui.Block) *AnimatedImage { a.Block = &b; return a }

// Render draws the animated image. On the first render, it transmits all frames.
func (a *AnimatedImage) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if a.Block != nil {
		inner = a.Block.Render(buf, area)
	}

	if inner.IsEmpty() {
		return
	}

	// Place the base image
	placement := tui.NewImagePlacement(inner.X, inner.Y, inner.Width, inner.Height).
		WithImageID(a.imageID)
	buf.AddImage(placement)
}

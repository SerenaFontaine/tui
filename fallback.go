package tui

import (
	"bytes"
	goimage "image"
	"image/color"
	"image/png"
)

// transmitPreference is the priority order for transmission method downgrade.
var transmitPreference = []TransmitMethod{
	TransmitFromSharedMem,
	TransmitFromFile,
	TransmitFromTempFile,
	TransmitDirect,
}

// processFallbacks checks each image in the buffer against capabilities
// and either downgrades, re-encodes, or replaces with a placeholder.
func processFallbacks(buf *Buffer, caps Capabilities, theme Theme) {
	if !caps.HasKGP() {
		for _, img := range buf.Images {
			area := NewRect(img.X, img.Y, img.Columns, img.Rows)
			renderPlaceholder(buf, area, img.alt, img.isAnimation, theme)
		}
		buf.Images = buf.Images[:0]
		return
	}

	// kept reuses the same backing array: safe because tryFallback is 1-in/0-or-1-out.
	kept := buf.Images[:0]
	for _, img := range buf.Images {
		// Handle animations: downgrade to static if animation not supported
		if img.isAnimation && !caps.Animation {
			img.isAnimation = false
		}

		result, ok := tryFallback(img, caps)
		if ok {
			kept = append(kept, result)
		} else {
			area := NewRect(img.X, img.Y, img.Columns, img.Rows)
			renderPlaceholder(buf, area, img.alt, img.isAnimation, theme)
		}
	}
	buf.Images = kept
}

// tryFallback attempts to make an image placement work within the given capabilities.
func tryFallback(img ImagePlacement, caps Capabilities) (ImagePlacement, bool) {
	// Step 1: Check/downgrade transmission method
	if !caps.SupportsTransmission(img.transmit) {
		downgraded := false
		if img.data != nil {
			if caps.SupportsTransmission(TransmitDirect) {
				img.transmit = TransmitDirect
				downgraded = true
			}
		}
		if img.filePath != "" && !downgraded {
			for _, m := range transmitPreference {
				if m == TransmitFromSharedMem {
					continue
				}
				if caps.SupportsTransmission(m) && (m == TransmitFromFile || m == TransmitFromTempFile) {
					img.transmit = m
					downgraded = true
					break
				}
			}
		}
		if !downgraded {
			return ImagePlacement{}, false
		}
	}

	// Step 2: Check/re-encode format
	if !caps.SupportsFormat(img.format) {
		reencoded := false
		if caps.SupportsFormat(ImagePNG) && img.data != nil {
			if img.format == ImageRGBA && img.imgWidth > 0 && img.imgHeight > 0 {
				pngData, err := reencodeRGBAToPNG(img.data, img.imgWidth, img.imgHeight)
				if err == nil {
					img.format = ImagePNG
					img.data = pngData
					reencoded = true
				}
			} else if img.format == ImageRGB && img.imgWidth > 0 && img.imgHeight > 0 {
				pngData, err := reencodeRGBToPNG(img.data, img.imgWidth, img.imgHeight)
				if err == nil {
					img.format = ImagePNG
					img.data = pngData
					reencoded = true
				}
			}
		}
		if !reencoded {
			return ImagePlacement{}, false
		}
	}

	return img, true
}

// reencodeRGBAToPNG converts raw RGBA pixel data to PNG.
func reencodeRGBAToPNG(data []byte, width, height int) ([]byte, error) {
	expectedLen := width * height * 4
	if len(data) < expectedLen {
		return nil, png.FormatError("insufficient RGBA data")
	}

	img := &rawRGBAImage{data: data, width: width, height: height}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// reencodeRGBToPNG converts raw RGB pixel data to PNG.
func reencodeRGBToPNG(data []byte, width, height int) ([]byte, error) {
	expectedLen := width * height * 3
	if len(data) < expectedLen {
		return nil, png.FormatError("insufficient RGB data")
	}

	img := &rawRGBImage{data: data, width: width, height: height}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// rawRGBAImage implements image.Image for raw RGBA data.
type rawRGBAImage struct {
	data          []byte
	width, height int
}

func (r *rawRGBAImage) ColorModel() color.Model  { return color.RGBAModel }
func (r *rawRGBAImage) Bounds() goimage.Rectangle { return goimage.Rect(0, 0, r.width, r.height) }
func (r *rawRGBAImage) At(x, y int) color.Color {
	if x < 0 || x >= r.width || y < 0 || y >= r.height {
		return color.RGBA{}
	}
	i := (y*r.width + x) * 4
	return color.RGBA{R: r.data[i], G: r.data[i+1], B: r.data[i+2], A: r.data[i+3]}
}

// rawRGBImage implements image.Image for raw RGB data.
type rawRGBImage struct {
	data          []byte
	width, height int
}

func (r *rawRGBImage) ColorModel() color.Model  { return color.RGBAModel }
func (r *rawRGBImage) Bounds() goimage.Rectangle { return goimage.Rect(0, 0, r.width, r.height) }
func (r *rawRGBImage) At(x, y int) color.Color {
	if x < 0 || x >= r.width || y < 0 || y >= r.height {
		return color.RGBA{}
	}
	i := (y*r.width + x) * 3
	return color.RGBA{R: r.data[i], G: r.data[i+1], B: r.data[i+2], A: 255}
}

package tui

import "testing"

func TestProcessFallbacksNoKGP(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(2, 2, 10, 5).WithPNG([]byte("fake"))
	buf.AddImage(img)

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 0 {
		t.Errorf("images = %d, want 0 (all replaced with placeholders)", len(buf.Images))
	}
	if buf.Get(2, 2).Char != BorderSingle.TopLeft {
		t.Errorf("placeholder border at (2,2) = %q, want %q", buf.Get(2, 2).Char, BorderSingle.TopLeft)
	}
}

func TestProcessFallbacksFullKGP(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).WithPNG([]byte("fake"))
	buf.AddImage(img)

	caps := KittyCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 1 {
		t.Errorf("images = %d, want 1 (fully supported)", len(buf.Images))
	}
}

func TestProcessFallbacksTransmissionDowngrade(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).
		WithSharedMemory("test", 1024, ImagePNG)
	buf.AddImage(img)

	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile},
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	// SharedMem with no raw data can't be downgraded
	if len(buf.Images) != 0 {
		t.Errorf("images = %d, want 0 (SharedMem can't be downgraded without data)", len(buf.Images))
	}
}

func TestProcessFallbacksDirectTransmissionDowngrade(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).WithPNG([]byte("fakepng"))
	img.transmit = TransmitFromSharedMem
	img.shmName = "test"
	img.shmSize = 7
	buf.AddImage(img)

	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 1 {
		t.Errorf("images = %d, want 1 (should downgrade to Direct)", len(buf.Images))
	}
	if buf.Images[0].transmit != TransmitDirect {
		t.Errorf("transmit = %d, want TransmitDirect (%d)", buf.Images[0].transmit, TransmitDirect)
	}
}

func TestProcessFallbacksWithAltText(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(2, 2, 16, 6).
		WithPNG([]byte("fake")).
		WithAlt("Chart")
	buf.AddImage(img)

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	midY := 2 + 6/2
	if buf.Get(6, midY).Char != '[' {
		t.Errorf("alt label at (%d,%d) = %q, want '['", 6, midY, buf.Get(6, midY).Char)
	}
}

func TestProcessFallbacksFormatUnsupported(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).
		WithRGBA([]byte("fakedata"), 2, 2)
	buf.AddImage(img)

	// Terminal only supports PNG
	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	// RGBA with invalid/insufficient data can't be re-encoded to PNG
	// The fake data "fakedata" (8 bytes) is less than 2*2*4=16 bytes needed
	if len(buf.Images) != 0 {
		t.Errorf("images = %d, want 0 (invalid data can't be re-encoded)", len(buf.Images))
	}
}

func TestProcessFallbacksFormatReencode(t *testing.T) {
	buf := NewBuffer(20, 10)
	// Valid 2x2 RGBA data: 16 bytes (4 pixels * 4 bytes each)
	rgba := make([]byte, 2*2*4)
	for i := range rgba {
		rgba[i] = 0xFF
	}
	img := NewImagePlacement(0, 0, 10, 5).
		WithRGBA(rgba, 2, 2)
	buf.AddImage(img)

	// Terminal only supports PNG, not RGBA
	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	// Valid RGBA data should be re-encoded to PNG successfully
	if len(buf.Images) != 1 {
		t.Errorf("images = %d, want 1 (RGBA should re-encode to PNG)", len(buf.Images))
	}
	if len(buf.Images) > 0 && buf.Images[0].format != ImagePNG {
		t.Errorf("format = %d, want ImagePNG (%d)", buf.Images[0].format, ImagePNG)
	}
}

func TestProcessFallbacksAnimationToPlaceholder(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).
		WithImageID(42).
		WithAlt("Spinner").
		WithAnimation()
	buf.AddImage(img)

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 0 {
		t.Error("animation should become placeholder when no KGP")
	}
}

func TestProcessFallbacksAnimationToStatic(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).
		WithImageID(42).
		WithAnimation()
	buf.AddImage(img)

	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
		Animation:     false,
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 1 {
		t.Errorf("images = %d, want 1 (animation downgraded to static)", len(buf.Images))
	}
	if buf.Images[0].isAnimation {
		t.Error("should no longer be marked as animation")
	}
}

func TestProcessFallbacksMultipleImages(t *testing.T) {
	buf := NewBuffer(40, 20)
	img1 := NewImagePlacement(0, 0, 10, 5).WithPNG([]byte("fake"))
	img2 := NewImagePlacement(15, 0, 10, 5).WithPNG([]byte("fake"))
	buf.AddImage(img1)
	buf.AddImage(img2)

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 0 {
		t.Errorf("images = %d, want 0", len(buf.Images))
	}
	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Error("first placeholder missing")
	}
	if buf.Get(15, 0).Char != BorderSingle.TopLeft {
		t.Error("second placeholder missing")
	}
}

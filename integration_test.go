package tui

import "testing"

// TestFullFallbackPipeline simulates what App.render() does:
// component renders images, processFallbacks intercepts, placeholders appear.
func TestFullFallbackPipeline(t *testing.T) {
	buf := NewBuffer(60, 20)

	// Static image
	img1 := NewImagePlacement(0, 0, 20, 8).
		WithPNG([]byte("fakepng")).
		WithAlt("Photo")

	// Animation
	img2 := NewImagePlacement(25, 0, 20, 8).
		WithImageID(42).
		WithAnimation().
		WithAlt("Loading")

	// Image using shared memory
	img3 := NewImagePlacement(0, 10, 20, 8).
		WithSharedMemory("shmtest", 1024, ImagePNG)

	buf.AddImage(img1)
	buf.AddImage(img2)
	buf.AddImage(img3)

	// Terminal: basic KGP, PNG only, Direct only, no animation
	caps := Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
		Animation:     false,
	}
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	// img1: PNG via Direct — should survive
	// img2: Animation downgraded to static, has imageID — should survive (just a placement)
	// img3: SharedMem with no raw data — should become placeholder

	if len(buf.Images) != 2 {
		t.Errorf("images = %d, want 2", len(buf.Images))
	}

	// img3's area should have placeholder
	if buf.Get(0, 10).Char != BorderSingle.TopLeft {
		t.Error("SharedMem image should have placeholder border")
	}

	// img1's area should NOT have placeholder
	if buf.Get(0, 0).Char == BorderSingle.TopLeft {
		t.Error("PNG image should not have placeholder")
	}
}

// TestFullFallbackNoKGP verifies everything becomes placeholders.
func TestFullFallbackNoKGP(t *testing.T) {
	buf := NewBuffer(40, 10)

	buf.AddImage(NewImagePlacement(0, 0, 15, 5).WithPNG([]byte("fake")))
	buf.AddImage(NewImagePlacement(20, 0, 15, 5).WithImageID(1).WithAnimation())

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 0 {
		t.Errorf("images = %d, want 0", len(buf.Images))
	}

	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Error("first placeholder missing")
	}
	if buf.Get(20, 0).Char != BorderSingle.TopLeft {
		t.Error("second placeholder missing")
	}
}

// TestFullFallbackFullKGP verifies nothing is touched.
func TestFullFallbackFullKGP(t *testing.T) {
	buf := NewBuffer(40, 10)

	buf.AddImage(NewImagePlacement(0, 0, 15, 5).WithPNG([]byte("fake")))
	buf.AddImage(NewImagePlacement(20, 0, 15, 5).WithImageID(1).WithAnimation())

	caps := KittyCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 2 {
		t.Errorf("images = %d, want 2 (all supported)", len(buf.Images))
	}
}

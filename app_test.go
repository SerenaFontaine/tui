package tui

import "testing"

func TestWithCapabilitiesOption(t *testing.T) {
	caps := KittyCapabilities()
	app := &App{}
	opt := WithCapabilities(caps)
	opt(app)

	if !app.capabilities.HasKGP() {
		t.Error("WithCapabilities should set capabilities on App")
	}
	if !app.capabilities.Animation {
		t.Error("should preserve Animation field")
	}
	if !app.capabilitiesSet {
		t.Error("capabilitiesSet should be true")
	}
}

func TestAppCapabilitiesAccessor(t *testing.T) {
	app := &App{}
	app.capabilities = WezTermCapabilities()

	c := app.Capabilities()
	if !c.HasKGP() {
		t.Error("Capabilities() should return the set capabilities")
	}
	if c.Animation {
		t.Error("WezTerm should not have animation")
	}
}

func TestWithCapabilitiesSkipsDetection(t *testing.T) {
	app := &App{}
	app.capabilitiesSet = true
	app.capabilities = KittyCapabilities()
	app.detectCapabilities()

	if !app.capabilities.Animation {
		t.Error("detectCapabilities should not overwrite when capabilitiesSet is true")
	}
}

func TestDetectCapabilitiesDefault(t *testing.T) {
	app := &App{
		screen: &Screen{
			in:     &slowReader{},
			out:    &nopWriter{},
			events: make(chan Msg, 64),
			done:   make(chan struct{}),
			fd:     -1,
		},
	}
	app.detectCapabilities()

	if app.capabilities.HasKGP() {
		t.Error("default detection should set NoKGP when terminal does not respond")
	}
}

func TestWithThemeOption(t *testing.T) {
	app := &App{}
	opt := WithTheme(NordTheme)
	opt(app)

	if app.theme.Primary != NordTheme.Primary {
		t.Error("WithTheme should set the theme")
	}
}

func TestRenderCallsFallback(t *testing.T) {
	buf := NewBuffer(20, 10)
	img := NewImagePlacement(0, 0, 10, 5).WithPNG([]byte("fake"))
	buf.AddImage(img)

	caps := NoKGPCapabilities()
	theme := DefaultTheme

	processFallbacks(buf, caps, theme)

	if len(buf.Images) != 0 {
		t.Error("render should replace images with placeholders when no KGP")
	}
	if buf.Get(0, 0).Char != BorderSingle.TopLeft {
		t.Error("placeholder should be rendered")
	}
}

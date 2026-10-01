package tui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

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

// sizeProbe records the render area and quits on the first ResizeMsg.
type sizeProbe struct{ w, h int }

func (s *sizeProbe) Init() Cmd { return nil }

func (s *sizeProbe) Update(msg Msg) (Component, Cmd) {
	if _, ok := msg.(ResizeMsg); ok {
		return s, QuitCmd()
	}
	return s, nil
}

func (s *sizeProbe) Render(_ *Buffer, area Rect) { s.w, s.h = area.Width, area.Height }

// runApp runs c with custom I/O, sends msg, and waits for the app to quit.
// Returns everything written to the output.
func runApp(t *testing.T, c Component, msg Msg, opts ...Option) string {
	t.Helper()
	in := &blockingReader{ch: make(chan struct{})}
	defer close(in.ch)
	var out bytes.Buffer
	opts = append([]Option{
		WithInput(in),
		WithOutput(&out),
		WithSizeFunc(func() (int, int) { return 80, 24 }),
		WithCapabilities(NoKGPCapabilities()),
	}, opts...)
	app := NewApp(c, opts...)

	done := make(chan error)
	go func() { done <- app.Run() }()
	app.Send(msg)

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("app did not quit")
	}
	return out.String()
}

func TestAppExternalResize(t *testing.T) {
	probe := &sizeProbe{}
	runApp(t, probe, ResizeMsg{Width: 120, Height: 40})

	if probe.w != 120 || probe.h != 40 {
		t.Errorf("rendered at %dx%d, want 120x40", probe.w, probe.h)
	}
}

func TestAppMouseEnabledByDefault(t *testing.T) {
	out := runApp(t, &sizeProbe{}, QuitMsg{})
	if !strings.Contains(out, "\x1b[?1003h") {
		t.Error("mouse tracking should be enabled by default")
	}
}

func TestAppMouseDisabled(t *testing.T) {
	out := runApp(t, &sizeProbe{}, QuitMsg{}, WithMouseEnabled(false))
	if strings.Contains(out, "\x1b[?1003h") {
		t.Error("mouse tracking should not be enabled with WithMouseEnabled(false)")
	}
}

func TestAppDoubleQuit(t *testing.T) {
	app := &App{quit: make(chan struct{})}

	// Should not panic on a second close
	app.handleMsg(QuitMsg{})
	app.handleMsg(QuitMsg{})

	select {
	case <-app.quit:
	default:
		t.Error("quit channel should be closed")
	}
}

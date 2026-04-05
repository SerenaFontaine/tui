package tui

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// App manages the lifecycle of a TUI application.
type App struct {
	component Component
	screen    *Screen
	prev      *Buffer
	curr      *Buffer
	width     int
	height    int
	mu        sync.Mutex
	cmds      chan Cmd
	msgs      chan Msg
	quit      chan struct{}

	// Cursor state
	cursorX, cursorY int
	cursorVisible    bool

	// Image manager
	Images *ImageManager

	// Options
	altScreen    bool
	mouseEnabled bool
	title        string
	input        io.Reader
	output       io.Writer
	sizeFunc     func() (int, int)
}

// Option configures the App.
type Option func(*App)

// WithAltScreen enables the alternate screen buffer (enabled by default).
func WithAltScreen(enabled bool) Option {
	return func(a *App) { a.altScreen = enabled }
}

// WithMouseEnabled enables mouse event tracking (enabled by default).
func WithMouseEnabled(enabled bool) Option {
	return func(a *App) { a.mouseEnabled = enabled }
}

// WithTitle sets the terminal window title.
func WithTitle(title string) Option {
	return func(a *App) { a.title = title }
}

// WithInput sets a custom input reader instead of os.Stdin.
// When set, the screen will not manage raw mode or signal handling.
func WithInput(r io.Reader) Option {
	return func(a *App) { a.input = r }
}

// WithOutput sets a custom output writer instead of os.Stdout.
// When set, the screen will not manage raw mode or signal handling.
func WithOutput(w io.Writer) Option {
	return func(a *App) { a.output = w }
}

// WithSizeFunc sets a custom function to retrieve terminal dimensions.
func WithSizeFunc(f func() (width, height int)) Option {
	return func(a *App) { a.sizeFunc = f }
}

// Run creates an App for the given component and runs it.
// This is the simplest way to start a TUI application.
func Run(c Component, opts ...Option) error {
	app := NewApp(c, opts...)
	return app.Run()
}

// NewApp creates a new App with the given root component and options.
func NewApp(c Component, opts ...Option) *App {
	a := &App{
		component:    c,
		cmds:         make(chan Cmd, 64),
		msgs:         make(chan Msg, 64),
		quit:         make(chan struct{}),
		altScreen:    true,
		mouseEnabled: true,
		Images:       NewImageManager(),
	}
	for _, opt := range opts {
		opt(a)
	}

	if a.input != nil || a.output != nil {
		in := a.input
		if in == nil {
			in = os.Stdin
		}
		out := a.output
		if out == nil {
			out = os.Stdout
		}
		a.screen = &Screen{
			in:       in,
			out:      out,
			managed:  false,
			fd:       -1,
			sizeFunc: a.sizeFunc,
			events:   make(chan Msg, 64),
			done:     make(chan struct{}),
		}
	} else {
		a.screen = newScreen()
		a.screen.sizeFunc = a.sizeFunc
	}

	return a
}

// Run starts the application event loop. Blocks until the app quits.
func (a *App) Run() error {
	if err := a.screen.Start(); err != nil {
		return err
	}
	defer a.screen.Stop()

	// Set window title if specified
	if a.title != "" {
		a.screen.SetTitle(a.title)
	}

	// Handle SIGWINCH for terminal resize, SIGTSTP for suspend.
	// Only register signal handlers when managing the local terminal.
	var sigCh chan os.Signal
	if a.screen.managed {
		sigCh = make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGWINCH, syscall.SIGCONT)
		defer signal.Stop(sigCh)
	}

	// Initial size and render
	a.width, a.height = a.screen.Size()
	a.curr = NewBuffer(a.width, a.height)
	a.prev = nil

	// Run Init command
	if cmd := a.component.Init(); cmd != nil {
		a.runCmd(cmd)
	}

	// Initial render
	a.render()

	for {
		select {
		case <-a.quit:
			return nil

		case sig := <-sigCh:
			switch sig {
			case syscall.SIGWINCH:
				w, h := a.screen.Size()
				if w != a.width || h != a.height {
					a.width = w
					a.height = h
					a.curr.Resize(w, h)
					a.prev = nil
					a.handleMsg(ResizeMsg{Width: w, Height: h})
					a.render()
				}
			case syscall.SIGCONT:
				// Resuming from suspend — re-enter raw mode and redraw
				a.screen.Resume()
				a.prev = nil
				a.render()
				a.handleMsg(ResumeMsg{})
			}

		case msg := <-a.screen.Events():
			// Handle suspend request (only for managed terminals)
			if a.screen.managed {
				if km, ok := msg.(KeyMsg); ok && km.Type == KeyCtrlZ {
					a.suspend()
					continue
				}
			}
			a.handleMsg(msg)
			a.render()

		case msg := <-a.msgs:
			a.handleMsg(msg)
			a.render()
		}
	}
}

// Send sends an external message to the application.
// This is safe to call from any goroutine.
func (a *App) Send(msg Msg) {
	select {
	case a.msgs <- msg:
	case <-a.quit:
	}
}

// SetCursor shows the cursor at the given position.
func (a *App) SetCursor(x, y int) {
	a.cursorX = x
	a.cursorY = y
	a.cursorVisible = true
}

// HideCursor hides the cursor.
func (a *App) HideCursor() {
	a.cursorVisible = false
}

func (a *App) handleMsg(msg Msg) {
	switch msg.(type) {
	case QuitMsg:
		close(a.quit)
		return
	}

	switch m := msg.(type) {
	case BatchMsg:
		for _, cmd := range m {
			a.runCmd(cmd)
		}
		return
	case CursorMsg:
		if m.Visible {
			a.cursorX = m.X
			a.cursorY = m.Y
			a.cursorVisible = true
		} else {
			a.cursorVisible = false
		}
		return
	default:
		_ = m
	}

	var cmd Cmd
	a.component, cmd = a.component.Update(msg)
	if cmd != nil {
		a.runCmd(cmd)
	}
}

func (a *App) runCmd(cmd Cmd) {
	go func() {
		msg := cmd()
		if msg != nil {
			a.Send(msg)
		}
	}()
}

func (a *App) render() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.curr.Clear()
	area := NewRect(0, 0, a.width, a.height)
	a.component.Render(a.curr, area)

	// Diff and flush
	output := a.curr.Diff(a.prev)
	if output != "" {
		a.screen.Flush(output)
	}

	// Render images via KGP
	for _, img := range a.curr.Images {
		a.screen.Flush(img.Encode())
	}

	// Handle cursor
	if a.cursorVisible {
		a.screen.ShowCursor(a.cursorX, a.cursorY)
	} else {
		a.screen.HideCursor()
	}

	// Swap buffers
	old := a.prev
	a.prev = a.curr
	if old != nil && old.Width == a.width && old.Height == a.height {
		a.curr = old
	} else {
		a.curr = NewBuffer(a.width, a.height)
	}
}

func (a *App) suspend() {
	a.screen.Suspend()
	a.handleMsg(SuspendMsg{})
	// Send SIGTSTP to self to actually suspend
	syscall.Kill(syscall.Getpid(), syscall.SIGTSTP)
}

// --- Cursor message ---

// CursorMsg is returned from a Cmd to set the cursor position.
type CursorMsg struct {
	X, Y    int
	Visible bool
}

// ShowCursorCmd returns a Cmd that shows the cursor at (x, y).
func ShowCursorCmd(x, y int) Cmd {
	return func() Msg {
		return CursorMsg{X: x, Y: y, Visible: true}
	}
}

// HideCursorCmd returns a Cmd that hides the cursor.
func HideCursorCmd() Cmd {
	return func() Msg {
		return CursorMsg{Visible: false}
	}
}

// --- Suspend/Resume messages ---

// SuspendMsg is sent when the app is about to be suspended (Ctrl+Z).
type SuspendMsg struct{}

// ResumeMsg is sent when the app resumes from suspension.
type ResumeMsg struct{}

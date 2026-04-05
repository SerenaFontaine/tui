package tui

// Component is the core interface for TUI applications.
// It follows the Elm Architecture: Init → Update → Render loop.
type Component interface {
	// Init returns an initial command to run when the component starts.
	// Return nil if no initial command is needed.
	Init() Cmd

	// Update handles a message and returns the updated component and
	// an optional command to run. The component may mutate itself and
	// return itself, or return a completely new component.
	Update(msg Msg) (Component, Cmd)

	// Render draws the component into the buffer within the given area.
	Render(buf *Buffer, area Rect)
}

// Widget is a simpler interface for components that only need to render.
type Widget interface {
	Render(buf *Buffer, area Rect)
}

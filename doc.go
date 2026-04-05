// Package tui provides a terminal user interface framework for Go with
// first-class Kitty Graphics Protocol support.
//
// TUI follows the Elm Architecture with an Init/Update/Render loop.
// Applications implement the [Component] interface and are run with [Run]:
//
//	type myApp struct{}
//
//	func (a *myApp) Init() tui.Cmd                              { return nil }
//	func (a *myApp) Update(msg tui.Msg) (tui.Component, tui.Cmd) { return a, nil }
//	func (a *myApp) Render(buf *tui.Buffer, area tui.Rect)       { ... }
//
//	func main() {
//	    tui.Run(&myApp{})
//	}
//
// The framework provides buffer-based rendering with automatic diffing,
// a flexible constraint-based layout system, mouse and keyboard event handling,
// and full integration with the Kitty Graphics Protocol for inline images
// and animation via the [github.com/SerenaFontaine/kgp] package.
//
// Built-in widgets are available in the [tui/widget] sub-package.
package tui

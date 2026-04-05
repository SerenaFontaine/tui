// Package widget provides built-in UI components for the tui framework.
//
// Widgets fall into two categories: passive widgets that only render
// (like [Text] and [Progress]) and interactive widgets that handle
// input events (like [Input], [List], and [Table]).
//
// Interactive widgets follow a consistent pattern:
//
//	// In your Component's Update method:
//	a.list, cmd = a.list.Update(msg)
//
//	// In your Component's Render method:
//	a.list.Render(buf, area)
//
// Most widgets support an optional [tui.Block] for borders and titles:
//
//	block := tui.NewBlock()
//	block.Title = "Items"
//	a.list.SetBlock(block)
//
// Image widgets use the Kitty Graphics Protocol via [tui.ImagePlacement]
// for rendering pixel-based graphics inline in the terminal.
package widget

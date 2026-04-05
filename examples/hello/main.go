package main

import (
	"log"

	"github.com/SerenaFontaine/tui"
	"github.com/SerenaFontaine/tui/widget"
)

type app struct {
	input *widget.Input
	items []string
	list  *widget.List
}

func newApp() *app {
	items := []string{
		"Hello, World!",
		"Welcome to TUI",
		"Press 'q' to quit",
		"Use arrow keys to navigate",
		"Type in the input below",
	}
	return &app{
		input: widget.NewInput("Type something..."),
		items: items,
		list:  widget.NewList(items),
	}
}

func (a *app) Init() tui.Cmd {
	a.input.Focused = true
	return nil
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyCtrlC, tui.KeyEscape:
			return a, tui.QuitCmd()
		case tui.KeyTab:
			a.input.Focused = !a.input.Focused
			return a, nil
		}
	}

	if a.input.Focused {
		var cmd tui.Cmd
		a.input, cmd = a.input.Update(msg)
		return a, cmd
	}

	var cmd tui.Cmd
	a.list, cmd = a.list.Update(msg)
	return a, cmd
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
	// Split: list on top, input at bottom
	regions := tui.VSplit(area, tui.Flex(1), tui.Fixed(3))

	// List with border
	block := tui.NewBlock()
	block.Title = "Items"
	block.Style = tui.NewStyle().Fg(tui.Cyan)
	a.list.SetBlock(block)
	a.list.SetSelectedStyle(tui.NewStyle().Bg(tui.Blue).Fg(tui.White))
	a.list.Render(buf, regions[0])

	// Input with border
	inputBlock := tui.NewBlock()
	inputBlock.Title = "Input"
	if a.input.Focused {
		inputBlock.Style = tui.NewStyle().Fg(tui.Green)
	}
	a.input.SetBlock(inputBlock)
	a.input.Render(buf, regions[1])
}

func main() {
	if err := tui.Run(newApp()); err != nil {
		log.Fatal(err)
	}
}

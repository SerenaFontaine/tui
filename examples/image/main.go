package main

import (
	"image"
	"image/color"
	"log"

	"tui"
	"tui/widget"
)

type app struct {
	info *widget.Text
}

func newApp() *app {
	return &app{
		info: widget.NewText("Kitty Graphics Protocol Image Demo\n\nPress 'q' or Ctrl+C to quit.\n\nThis example generates a gradient image and\ndisplays it using the Kitty Graphics Protocol.\nYour terminal must support KGP (Kitty, WezTerm, etc)."),
	}
}

func (a *app) Init() tui.Cmd { return nil }

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyCtrlC, tui.KeyEscape:
			return a, tui.QuitCmd()
		case tui.KeyRune:
			if msg.Rune == 'q' {
				return a, tui.QuitCmd()
			}
		}
	}
	return a, nil
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
	regions := tui.VSplit(area,
		tui.Fixed(3), // title
		tui.Flex(1),  // image area
		tui.Fixed(1), // status
	)

	// Title
	titleBlock := tui.NewBlock()
	titleBlock.Title = "Image Demo"
	titleBlock.Style = tui.NewStyle().Fg(tui.Magenta)
	titleBlock.Border = tui.BorderRounded
	titleText := widget.NewText("KGP Image Rendering")
	titleText.SetAlignment(widget.AlignCenter)
	titleText.SetStyle(tui.NewStyle().Fg(tui.BrightWhite).Bold(true))
	inner := titleBlock.Render(buf, regions[0])
	titleText.Render(buf, inner)

	// Image area - split into image and info
	cols := tui.HSplit(regions[1], tui.Percent(60), tui.Flex(1))

	// Generate and display gradient image
	imgBlock := tui.NewBlock()
	imgBlock.Title = "Gradient"
	imgBlock.Border = tui.BorderRounded
	imgBlock.Style = tui.NewStyle().Fg(tui.Cyan)
	imgInner := imgBlock.Render(buf, cols[0])

	gradient := generateGradient(256, 256)
	imgWidget := widget.NewImage(gradient)
	imgWidget.Render(buf, imgInner)

	// Info panel
	infoBlock := tui.NewBlock()
	infoBlock.Title = "Info"
	infoBlock.Border = tui.BorderRounded
	a.info.SetStyle(tui.NewStyle().Fg(tui.White))
	a.info.SetBlock(infoBlock)
	a.info.Render(buf, cols[1])

	// Status bar
	status := " Press 'q' to quit"
	statusStyle := tui.NewStyle().Bg(tui.Magenta).Fg(tui.White)
	for x := regions[2].X; x < regions[2].Right(); x++ {
		buf.SetChar(x, regions[2].Y, ' ', statusStyle)
	}
	buf.SetString(regions[2].X, regions[2].Y, status, statusStyle)
}

func generateGradient(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(x * 255 / width)
			g := uint8(y * 255 / height)
			b := uint8(255 - r)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func main() {
	if err := tui.Run(newApp()); err != nil {
		log.Fatal(err)
	}
}

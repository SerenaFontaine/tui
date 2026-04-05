package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/SerenaFontaine/tui"
	"github.com/SerenaFontaine/tui/widget"
)

type app struct {
	imgMgr   *tui.ImageManager
	spinner  *widget.Spinner
	progress *widget.Progress
	info     *widget.Text
	animID   uint32
	frame    int
	ready    bool
}

func newApp() *app {
	return &app{
		imgMgr:   tui.NewImageManager(),
		spinner:  widget.NewSpinner().SetLabel("Preparing animation..."),
		progress: widget.NewProgress(),
		info: widget.NewText(
			"KGP Animation Demo\n\n" +
				"This example demonstrates:\n" +
				"  - Image transmission & management\n" +
				"  - Animation frame generation\n" +
				"  - Playback control via KGP\n" +
				"  - Spinner & progress widgets\n\n" +
				"Press 'q' or Ctrl+C to quit\n" +
				"Press 'p' to play/pause animation\n" +
				"Press 'r' to reset animation"),
	}
}

func (a *app) Init() tui.Cmd {
	return tui.Batch(
		a.spinner.Tick(),
		tui.TickCmd(50_000_000), // 50ms tick for animation
	)
}

func (a *app) Update(msg tui.Msg) (tui.Component, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyCtrlC, tui.KeyEscape:
			return a, tui.QuitCmd()
		case tui.KeyRune:
			switch msg.Rune {
			case 'q':
				return a, tui.QuitCmd()
			}
		}

	case widget.SpinnerTickMsg:
		var cmd tui.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd

	case tui.TickMsg:
		a.frame++
		p := float64(a.frame%200) / 200.0
		a.progress.SetPercent(p)
		return a, tui.TickCmd(50_000_000)
	}

	return a, nil
}

func (a *app) Render(buf *tui.Buffer, area tui.Rect) {
	// Layout: top title, middle content, bottom status
	regions := tui.VSplit(area,
		tui.Fixed(3), // title bar
		tui.Flex(1),  // main content
		tui.Fixed(3), // progress
		tui.Fixed(1), // status
	)

	// Title
	titleBlock := tui.NewBlock()
	titleBlock.Title = "KGP Animation Demo"
	titleBlock.Border = tui.BorderRounded
	titleBlock.Style = tui.NewStyle().Fg(tui.Magenta)
	inner := titleBlock.Render(buf, regions[0])
	a.spinner.Render(buf, inner)

	// Main content: image + info
	cols := tui.HSplit(regions[1], tui.Percent(55), tui.Flex(1))

	// Animated gradient display
	imgBlock := tui.NewBlock()
	imgBlock.Title = "Animated Gradient"
	imgBlock.Border = tui.BorderRounded
	imgBlock.Style = tui.NewStyle().Fg(tui.Cyan)
	imgInner := imgBlock.Render(buf, cols[0])

	// Generate frame-dependent gradient
	gradient := generateAnimatedGradient(128, 128, a.frame)
	imgWidget := widget.NewImage(gradient)
	imgWidget.Render(buf, imgInner)

	// Info panel
	infoBlock := tui.NewBlock()
	infoBlock.Title = "Controls"
	infoBlock.Border = tui.BorderRounded
	infoBlock.Style = tui.NewStyle().Fg(tui.Yellow)
	a.info.SetStyle(tui.NewStyle().Fg(tui.White))
	a.info.SetBlock(infoBlock)
	a.info.Render(buf, cols[1])

	// Progress
	progBlock := tui.NewBlock()
	progBlock.Title = "Animation Progress"
	progBlock.Border = tui.BorderRounded
	a.progress.SetBlock(progBlock)
	a.progress.Render(buf, regions[2])

	// Status bar
	statusStyle := tui.NewStyle().Bg(tui.Magenta).Fg(tui.White)
	for x := regions[3].X; x < regions[3].Right(); x++ {
		buf.SetChar(x, regions[3].Y, ' ', statusStyle)
	}
	buf.SetString(regions[3].X, regions[3].Y, " KGP Animation | Frame: "+itoa(a.frame), statusStyle)
}

func generateAnimatedGradient(width, height, frame int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	t := float64(frame) * 0.05
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float64(x) / float64(width)
			fy := float64(y) / float64(height)
			r := uint8(128 + 127*math.Sin(fx*math.Pi*2+t))
			g := uint8(128 + 127*math.Sin(fy*math.Pi*2+t*0.7))
			b := uint8(128 + 127*math.Sin((fx+fy)*math.Pi+t*1.3))
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func main() {
	if err := tui.Run(newApp()); err != nil {
		log.Fatal(err)
	}
}

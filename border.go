package tui

// BorderStyle defines the characters used to draw borders.
type BorderStyle struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
}

// Predefined border styles.
var (
	BorderNone = BorderStyle{}

	BorderSingle = BorderStyle{
		TopLeft: '┌', TopRight: '┐',
		BottomLeft: '└', BottomRight: '┘',
		Horizontal: '─', Vertical: '│',
	}

	BorderDouble = BorderStyle{
		TopLeft: '╔', TopRight: '╗',
		BottomLeft: '╚', BottomRight: '╝',
		Horizontal: '═', Vertical: '║',
	}

	BorderRounded = BorderStyle{
		TopLeft: '╭', TopRight: '╮',
		BottomLeft: '╰', BottomRight: '╯',
		Horizontal: '─', Vertical: '│',
	}

	BorderThick = BorderStyle{
		TopLeft: '┏', TopRight: '┓',
		BottomLeft: '┗', BottomRight: '┛',
		Horizontal: '━', Vertical: '┃',
	}

	BorderASCII = BorderStyle{
		TopLeft: '+', TopRight: '+',
		BottomLeft: '+', BottomRight: '+',
		Horizontal: '-', Vertical: '|',
	}
)

// Block is a container that draws a border and optional title around content.
type Block struct {
	Border BorderStyle
	Style  Style
	Title  string
}

// NewBlock creates a block with a single-line border.
func NewBlock() Block {
	return Block{Border: BorderSingle}
}

// Render draws the block's border onto the buffer in the given area.
// Returns the inner area available for content.
func (b Block) Render(buf *Buffer, area Rect) Rect {
	if area.IsEmpty() {
		return Rect{}
	}

	if b.Border == BorderNone {
		return area
	}

	style := b.Style

	// Corners
	buf.SetChar(area.X, area.Y, b.Border.TopLeft, style)
	buf.SetChar(area.Right()-1, area.Y, b.Border.TopRight, style)
	buf.SetChar(area.X, area.Bottom()-1, b.Border.BottomLeft, style)
	buf.SetChar(area.Right()-1, area.Bottom()-1, b.Border.BottomRight, style)

	// Top and bottom edges
	for x := area.X + 1; x < area.Right()-1; x++ {
		buf.SetChar(x, area.Y, b.Border.Horizontal, style)
		buf.SetChar(x, area.Bottom()-1, b.Border.Horizontal, style)
	}

	// Left and right edges
	for y := area.Y + 1; y < area.Bottom()-1; y++ {
		buf.SetChar(area.X, y, b.Border.Vertical, style)
		buf.SetChar(area.Right()-1, y, b.Border.Vertical, style)
	}

	// Title
	if b.Title != "" && area.Width > 4 {
		title := b.Title
		maxLen := area.Width - 4
		if len(title) > maxLen {
			title = title[:maxLen]
		}
		buf.SetChar(area.X+1, area.Y, ' ', style)
		buf.SetString(area.X+2, area.Y, title, style)
		buf.SetChar(area.X+2+len(title), area.Y, ' ', style)
	}

	// Return inner area (1 cell inset on all sides)
	return area.Inner(1)
}

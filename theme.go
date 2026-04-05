package tui

// Theme defines a coherent color scheme for an application.
type Theme struct {
	// Base colors
	Primary   Color
	Secondary Color
	Accent    Color

	// Backgrounds
	Background       Color
	Surface          Color
	SurfaceHighlight Color

	// Text
	TextPrimary   Color
	TextSecondary Color
	TextMuted     Color

	// Status
	Success Color
	Warning Color
	Error   Color
	Info    Color

	// Borders
	BorderColor      Color
	BorderFocusColor Color

	// Derived styles
	border      Style
	borderFocus Style
	title       Style
	statusBar   Style
	selected    Style
}

// Predefined themes.

// DefaultTheme is a dark theme with blue accents.
var DefaultTheme = Theme{
	Primary:          Blue,
	Secondary:        Cyan,
	Accent:           Magenta,
	Background:       NoColor,
	Surface:          NoColor,
	SurfaceHighlight: ANSI256(237),
	TextPrimary:      White,
	TextSecondary:    BrightWhite,
	TextMuted:        BrightBlack,
	Success:          Green,
	Warning:          Yellow,
	Error:            Red,
	Info:             Cyan,
	BorderColor:      BrightBlack,
	BorderFocusColor: Blue,
}

// NordTheme is inspired by the Nord color palette.
var NordTheme = Theme{
	Primary:          RGB(136, 192, 208),
	Secondary:        RGB(129, 161, 193),
	Accent:           RGB(180, 142, 173),
	Background:       RGB(46, 52, 64),
	Surface:          RGB(59, 66, 82),
	SurfaceHighlight: RGB(67, 76, 94),
	TextPrimary:      RGB(236, 239, 244),
	TextSecondary:    RGB(229, 233, 240),
	TextMuted:        RGB(76, 86, 106),
	Success:          RGB(163, 190, 140),
	Warning:          RGB(235, 203, 139),
	Error:            RGB(191, 97, 106),
	Info:             RGB(136, 192, 208),
	BorderColor:      RGB(76, 86, 106),
	BorderFocusColor: RGB(136, 192, 208),
}

// GruvboxTheme is inspired by the Gruvbox color scheme.
var GruvboxTheme = Theme{
	Primary:          RGB(215, 153, 33),
	Secondary:        RGB(152, 151, 26),
	Accent:           RGB(211, 134, 155),
	Background:       RGB(40, 40, 40),
	Surface:          RGB(60, 56, 54),
	SurfaceHighlight: RGB(80, 73, 69),
	TextPrimary:      RGB(235, 219, 178),
	TextSecondary:    RGB(213, 196, 161),
	TextMuted:        RGB(146, 131, 116),
	Success:          RGB(152, 151, 26),
	Warning:          RGB(215, 153, 33),
	Error:            RGB(204, 36, 29),
	Info:             RGB(69, 133, 136),
	BorderColor:      RGB(102, 92, 84),
	BorderFocusColor: RGB(215, 153, 33),
}

// MonochromeTheme uses only ANSI base colors for maximum compatibility.
var MonochromeTheme = Theme{
	Primary:          White,
	Secondary:        BrightWhite,
	Accent:           White,
	Background:       NoColor,
	Surface:          NoColor,
	SurfaceHighlight: NoColor,
	TextPrimary:      White,
	TextSecondary:    BrightWhite,
	TextMuted:        BrightBlack,
	Success:          White,
	Warning:          White,
	Error:            White,
	Info:             White,
	BorderColor:      White,
	BorderFocusColor: BrightWhite,
}

// --- Derived style accessors ---

// BorderStyled returns the style for unfocused borders.
func (t Theme) BorderStyled() Style {
	return NewStyle().Fg(t.BorderColor)
}

// BorderFocusStyle returns the style for focused borders.
func (t Theme) BorderFocusStyle() Style {
	return NewStyle().Fg(t.BorderFocusColor)
}

// TitleStyle returns a style for titles.
func (t Theme) TitleStyle() Style {
	return NewStyle().Fg(t.Primary).Bold(true)
}

// StatusBarStyle returns a style for status bars.
func (t Theme) StatusBarStyle() Style {
	return NewStyle().Fg(t.TextPrimary).Bg(t.Primary)
}

// SelectedStyle returns a style for selected items.
func (t Theme) SelectedStyle() Style {
	return NewStyle().Fg(t.TextPrimary).Bg(t.Primary)
}

// ErrorStyle returns a style for error text.
func (t Theme) ErrorStyle() Style {
	return NewStyle().Fg(t.Error)
}

// WarningStyle returns a style for warning text.
func (t Theme) WarningStyle() Style {
	return NewStyle().Fg(t.Warning)
}

// SuccessStyle returns a style for success text.
func (t Theme) SuccessStyle() Style {
	return NewStyle().Fg(t.Success)
}

// InfoStyle returns a style for info text.
func (t Theme) InfoStyle() Style {
	return NewStyle().Fg(t.Info)
}

// MutedStyle returns a style for muted/secondary text.
func (t Theme) MutedStyle() Style {
	return NewStyle().Fg(t.TextMuted)
}

// Block creates a Block styled with the theme.
func (t Theme) Block(title string, focused bool) Block {
	b := NewBlock()
	b.Title = title
	if focused {
		b.Style = t.BorderFocusStyle()
	} else {
		b.Style = t.BorderStyled()
	}
	return b
}

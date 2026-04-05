---
title: Style & Color
weight: 3
---

## Color

`Color` represents a terminal color. The zero value is "no color" (terminal default).

### Constructors

| Function | Description |
|----------|-------------|
| `RGB(r, g, b uint8)` | 24-bit true color |
| `Hex(string)` | From hex string (`"#FF6432"` or `"FF6432"`) |
| `ANSI256(uint8)` | 256-color palette index |

### ANSI 16 Constants

Standard: `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White`

Bright: `BrightBlack`, `BrightRed`, `BrightGreen`, `BrightYellow`, `BrightBlue`, `BrightMagenta`, `BrightCyan`, `BrightWhite`

### Methods

| Method | Description |
|--------|-------------|
| `IsZero() bool` | True if this is the default (no) color |

## Style

`Style` defines visual attributes for a cell. It is a value type with a chainable builder.

### Builder Methods

All methods return a new `Style` value:

| Method | Description |
|--------|-------------|
| `Fg(Color)` | Set foreground color |
| `Bg(Color)` | Set background color |
| `Bold(bool)` | Bold attribute |
| `Dim(bool)` | Dim/faint attribute |
| `Italic(bool)` | Italic attribute |
| `Underline(bool)` | Underline attribute |
| `Blink(bool)` | Blink attribute |
| `Reverse(bool)` | Reverse video attribute |
| `Strikethrough(bool)` | Strikethrough attribute |

### Accessors

| Method | Description |
|--------|-------------|
| `Foreground() Color` | Get foreground color |
| `Background() Color` | Get background color |

### Example

```go
style := tui.NewStyle().
    Fg(tui.RGB(200, 100, 50)).
    Bg(tui.Black).
    Bold(true).
    Italic(true)
```

## Theme

`Theme` defines a coherent color scheme. Use derived style methods for consistent visuals.

### Built-in Themes

| Theme | Description |
|-------|-------------|
| `DefaultTheme` | Dark theme with blue accents |
| `NordTheme` | Nord color palette |
| `GruvboxTheme` | Gruvbox color scheme |
| `MonochromeTheme` | ANSI-only for maximum compatibility |

### Theme Fields

| Field | Description |
|-------|-------------|
| `Primary`, `Secondary`, `Accent` | Base colors |
| `Background`, `Surface`, `SurfaceHighlight` | Background layers |
| `TextPrimary`, `TextSecondary`, `TextMuted` | Text colors |
| `Success`, `Warning`, `Error`, `Info` | Status colors |
| `BorderColor`, `BorderFocusColor` | Border colors |

### Derived Styles

| Method | Description |
|--------|-------------|
| `TitleStyle()` | Bold primary text |
| `StatusBarStyle()` | Primary bg with white text |
| `SelectedStyle()` | Primary bg for selection |
| `ErrorStyle()` | Error colored text |
| `WarningStyle()` | Warning colored text |
| `SuccessStyle()` | Success colored text |
| `InfoStyle()` | Info colored text |
| `MutedStyle()` | Muted/secondary text |
| `BorderStyled()` | Unfocused border style |
| `BorderFocusStyle()` | Focused border style |
| `Block(title, focused)` | Create a themed Block |

## Styled Text

Rich inline text with mixed formatting.

| Type | Description |
|------|-------------|
| `Span` | Run of text with uniform style |
| `StyledLine` | Line of multiple spans |
| `StyledText` | Block of styled lines |

### Span Constructors

| Function | Description |
|----------|-------------|
| `NewSpan(text, style)` | Custom styled span |
| `PlainSpan(text)` | Default style |
| `BoldSpan(text)` | Bold text |
| `ColorSpan(text, color)` | Colored text |
| `StyledSpan(text, style)` | Alias for NewSpan |

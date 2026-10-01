package tui

import "golang.org/x/text/width"

// RuneWidth returns the number of terminal columns a rune occupies: 2 for
// East Asian wide and fullwidth runes (CJK, most emoji), otherwise 1.
// Ambiguous-width runes are treated as narrow, matching most terminals.
func RuneWidth(r rune) int {
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	}
	return 1
}

// StringWidth returns the number of terminal columns a string occupies.
func StringWidth(s string) int {
	n := 0
	for _, r := range s {
		n += RuneWidth(r)
	}
	return n
}

// Truncate shortens s to at most maxWidth columns without splitting a rune.
// A wide rune that would straddle the limit is dropped.
func Truncate(s string, maxWidth int) string {
	n := 0
	for i, r := range s {
		n += RuneWidth(r)
		if n > maxWidth {
			return s[:i]
		}
	}
	return s
}

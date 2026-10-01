package tui

import "golang.org/x/text/width"

// runeWidth returns the number of terminal columns a rune occupies: 2 for
// East Asian wide and fullwidth runes (CJK, most emoji), otherwise 1.
// Ambiguous-width runes are treated as narrow, matching most terminals.
func runeWidth(r rune) int {
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	}
	return 1
}

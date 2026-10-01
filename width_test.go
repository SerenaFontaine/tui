package tui

import "testing"

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		r    rune
		want int
	}{
		{'a', 1},
		{' ', 1},
		{'é', 1},
		{'─', 1}, // box drawing (ambiguous) stays narrow
		{'日', 2},
		{'한', 2},
		{'Ａ', 2}, // fullwidth Latin
		{'ｱ', 1}, // halfwidth katakana
		{'😀', 2},
		{'🚀', 2},
	}

	for _, tt := range tests {
		if got := runeWidth(tt.r); got != tt.want {
			t.Errorf("runeWidth(%q) = %d, want %d", tt.r, got, tt.want)
		}
	}
}

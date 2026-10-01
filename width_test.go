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
		if got := RuneWidth(tt.r); got != tt.want {
			t.Errorf("RuneWidth(%q) = %d, want %d", tt.r, got, tt.want)
		}
	}
}

func TestStringWidth(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"hello", 5},
		{"héllo", 5},
		{"日本", 4},
		{"a日b", 4},
		{"🚀go", 4},
	}

	for _, tt := range tests {
		if got := StringWidth(tt.s); got != tt.want {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s        string
		maxWidth int
		want     string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hel"},
		{"hello", 0, ""},
		{"héllo", 2, "hé"}, // multi-byte rune is kept whole
		{"日本語", 4, "日本"},
		{"日本語", 3, "日"}, // 本 would straddle the limit
		{"a日", 1, "a"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		if got := Truncate(tt.s, tt.maxWidth); got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.maxWidth, got, tt.want)
		}
	}
}

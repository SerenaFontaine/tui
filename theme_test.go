package tui

import "testing"

func TestThemeDerivedStyles(t *testing.T) {
	themes := []struct {
		name  string
		theme Theme
	}{
		{"default", DefaultTheme},
		{"nord", NordTheme},
		{"gruvbox", GruvboxTheme},
		{"monochrome", MonochromeTheme},
	}

	for _, tt := range themes {
		s := tt.theme.TitleStyle()
		if s == (Style{}) {
			t.Errorf("%s: TitleStyle() should not be zero", tt.name)
		}

		s = tt.theme.StatusBarStyle()
		if s == (Style{}) {
			t.Errorf("%s: StatusBarStyle() should not be zero", tt.name)
		}

		s = tt.theme.SelectedStyle()
		if s == (Style{}) {
			t.Errorf("%s: SelectedStyle() should not be zero", tt.name)
		}

		s = tt.theme.ErrorStyle()
		if s == (Style{}) {
			t.Errorf("%s: ErrorStyle() should not be zero", tt.name)
		}
	}
}

func TestThemeBlock(t *testing.T) {
	theme := DefaultTheme

	focused := theme.Block("Title", true)
	if focused.Title != "Title" {
		t.Errorf("focused block title = %q, want 'Title'", focused.Title)
	}

	unfocused := theme.Block("Title", false)
	if unfocused.Style == focused.Style {
		t.Error("focused and unfocused blocks should have different styles")
	}
}

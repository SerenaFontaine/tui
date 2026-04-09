package tui

import "testing"

func TestCapabilitiesHasKGP(t *testing.T) {
	c := Capabilities{KGP: true}
	if !c.HasKGP() {
		t.Error("HasKGP() should return true when KGP is true")
	}

	c = Capabilities{KGP: false}
	if c.HasKGP() {
		t.Error("HasKGP() should return false when KGP is false")
	}
}

func TestCapabilitiesSupportsFormat(t *testing.T) {
	c := Capabilities{
		KGP:     true,
		Formats: []ImageFormat{ImagePNG, ImageRGBA},
	}
	if !c.SupportsFormat(ImagePNG) {
		t.Error("should support PNG")
	}
	if !c.SupportsFormat(ImageRGBA) {
		t.Error("should support RGBA")
	}
	if c.SupportsFormat(ImageRGB) {
		t.Error("should not support RGB")
	}
}

func TestCapabilitiesSupportsTransmission(t *testing.T) {
	c := Capabilities{
		KGP:           true,
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile},
	}
	if !c.SupportsTransmission(TransmitDirect) {
		t.Error("should support Direct")
	}
	if c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("should not support SharedMem")
	}
}

func TestCapabilitiesZeroValue(t *testing.T) {
	var c Capabilities
	if c.HasKGP() {
		t.Error("zero-value Capabilities should not have KGP")
	}
	if c.SupportsFormat(ImagePNG) {
		t.Error("zero-value should not support any format")
	}
	if c.SupportsTransmission(TransmitDirect) {
		t.Error("zero-value should not support any transmission")
	}
}

func TestKittyCapabilities(t *testing.T) {
	c := KittyCapabilities()
	if !c.HasKGP() {
		t.Error("Kitty should have KGP")
	}
	if !c.Animation {
		t.Error("Kitty should support animation")
	}
	if !c.Unicode {
		t.Error("Kitty should support unicode placeholders")
	}
	if !c.SupportsFormat(ImagePNG) {
		t.Error("Kitty should support PNG")
	}
	if !c.SupportsFormat(ImageRGBA) {
		t.Error("Kitty should support RGBA")
	}
	if !c.SupportsFormat(ImageRGB) {
		t.Error("Kitty should support RGB")
	}
	if !c.SupportsTransmission(TransmitDirect) {
		t.Error("Kitty should support direct transmission")
	}
	if !c.SupportsTransmission(TransmitFromFile) {
		t.Error("Kitty should support file transmission")
	}
	if !c.SupportsTransmission(TransmitFromTempFile) {
		t.Error("Kitty should support temp file transmission")
	}
	if !c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("Kitty should support shared memory")
	}
}

func TestWezTermCapabilities(t *testing.T) {
	c := WezTermCapabilities()
	if !c.HasKGP() {
		t.Error("WezTerm should have KGP")
	}
	if c.Animation {
		t.Error("WezTerm should not support animation")
	}
	if c.Unicode {
		t.Error("WezTerm should not support unicode placeholders")
	}
	if !c.SupportsFormat(ImagePNG) {
		t.Error("WezTerm should support PNG")
	}
	if !c.SupportsFormat(ImageRGBA) {
		t.Error("WezTerm should support RGBA")
	}
	if !c.SupportsFormat(ImageRGB) {
		t.Error("WezTerm should support RGB")
	}
	if c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("WezTerm should not support shared memory")
	}
}

func TestKonsoleCapabilities(t *testing.T) {
	c := KonsoleCapabilities()
	if !c.HasKGP() {
		t.Error("Konsole should have KGP")
	}
	if c.Animation {
		t.Error("Konsole should not support animation")
	}
	if !c.SupportsFormat(ImagePNG) {
		t.Error("Konsole should support PNG")
	}
	if c.SupportsFormat(ImageRGBA) {
		t.Error("Konsole should not support RGBA")
	}
	if c.SupportsFormat(ImageRGB) {
		t.Error("Konsole should not support RGB")
	}
	if len(c.Formats) != 1 {
		t.Errorf("Konsole should support exactly one format, got %d", len(c.Formats))
	}
	if !c.SupportsTransmission(TransmitDirect) {
		t.Error("Konsole should support direct transmission")
	}
	if c.SupportsTransmission(TransmitFromFile) {
		t.Error("Konsole should not support file transmission")
	}
	if c.SupportsTransmission(TransmitFromTempFile) {
		t.Error("Konsole should not support temp file transmission")
	}
	if c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("Konsole should not support shared memory transmission")
	}
	if len(c.Transmissions) != 1 {
		t.Errorf("Konsole should support exactly one transmission method, got %d", len(c.Transmissions))
	}
}

func TestGhosttyCapabilities(t *testing.T) {
	c := GhosttyCapabilities()
	if !c.HasKGP() {
		t.Error("Ghostty should have KGP")
	}
	if !c.Animation {
		t.Error("Ghostty should support animation")
	}
	if !c.Unicode {
		t.Error("Ghostty should support unicode placeholders")
	}
	if !c.SupportsFormat(ImagePNG) || !c.SupportsFormat(ImageRGBA) || !c.SupportsFormat(ImageRGB) {
		t.Error("Ghostty should support PNG, RGBA, and RGB")
	}
	if !c.SupportsTransmission(TransmitDirect) || !c.SupportsTransmission(TransmitFromFile) ||
		!c.SupportsTransmission(TransmitFromTempFile) || !c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("Ghostty should support all transmission methods")
	}
}

func TestFootCapabilities(t *testing.T) {
	c := FootCapabilities()
	if !c.HasKGP() {
		t.Error("foot should have KGP")
	}
	if c.Animation {
		t.Error("foot should not support animation")
	}
	if c.Unicode {
		t.Error("foot should not support unicode placeholders")
	}
	if !c.SupportsFormat(ImagePNG) || !c.SupportsFormat(ImageRGBA) || !c.SupportsFormat(ImageRGB) {
		t.Error("foot should support PNG, RGBA, and RGB")
	}
	if !c.SupportsTransmission(TransmitDirect) || !c.SupportsTransmission(TransmitFromFile) ||
		!c.SupportsTransmission(TransmitFromTempFile) {
		t.Error("foot should support direct, file, and temp file transmission")
	}
	if c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("foot should not support shared memory")
	}
}

func TestITermCapabilities(t *testing.T) {
	c := ITermCapabilities()
	if !c.HasKGP() {
		t.Error("iTerm2 should have KGP")
	}
	if c.Animation {
		t.Error("iTerm2 should not support animation")
	}
	if c.Unicode {
		t.Error("iTerm2 should not support unicode placeholders")
	}
	if !c.SupportsFormat(ImagePNG) {
		t.Error("iTerm2 should support PNG")
	}
	if c.SupportsFormat(ImageRGBA) || c.SupportsFormat(ImageRGB) {
		t.Error("iTerm2 should not support RGBA or RGB")
	}
	if !c.SupportsTransmission(TransmitDirect) {
		t.Error("iTerm2 should support direct transmission")
	}
	if c.SupportsTransmission(TransmitFromFile) || c.SupportsTransmission(TransmitFromTempFile) ||
		c.SupportsTransmission(TransmitFromSharedMem) {
		t.Error("iTerm2 should only support direct transmission")
	}
}

func TestDetectCapabilitiesTermProgram(t *testing.T) {
	tests := []struct {
		name        string
		termProgram string
		wantKGP     bool
		wantAnim    bool
	}{
		{"kitty", "kitty", true, true},
		{"WezTerm", "WezTerm", true, false},
		{"konsole", "konsole", true, false},
		{"ghostty", "ghostty", true, true},
		{"foot", "foot", true, false},
		{"iTerm.app", "iTerm.app", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TERM_PROGRAM", tt.termProgram)
			app := &App{}
			app.detectCapabilities()

			if app.capabilities.HasKGP() != tt.wantKGP {
				t.Errorf("TERM_PROGRAM=%s: HasKGP() = %v, want %v", tt.termProgram, app.capabilities.HasKGP(), tt.wantKGP)
			}
			if app.capabilities.Animation != tt.wantAnim {
				t.Errorf("TERM_PROGRAM=%s: Animation = %v, want %v", tt.termProgram, app.capabilities.Animation, tt.wantAnim)
			}
		})
	}
}

func TestDetectCapabilitiesSkipsWhenSet(t *testing.T) {
	app := &App{}
	app.capabilitiesSet = true
	app.capabilities = KittyCapabilities()
	app.detectCapabilities()

	if !app.capabilities.Animation {
		t.Error("should not overwrite when capabilitiesSet is true")
	}
}

func TestNoKGPCapabilities(t *testing.T) {
	c := NoKGPCapabilities()
	if c.HasKGP() {
		t.Error("NoKGP should not have KGP")
	}
	if len(c.Formats) != 0 {
		t.Error("NoKGP should have no formats")
	}
	if len(c.Transmissions) != 0 {
		t.Error("NoKGP should have no transmissions")
	}
}

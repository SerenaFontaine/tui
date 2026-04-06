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

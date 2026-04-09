package tui

import (
	"errors"
	"os"
	"time"
)

const detectTimeout = 200 * time.Millisecond

// ErrQueryTimeout is returned when a terminal capability query times out.
var ErrQueryTimeout = errors.New("terminal query timed out")

// Capabilities describes the KGP features supported by the connected terminal.
type Capabilities struct {
	// KGP is true if any Kitty Graphics Protocol support was detected.
	KGP bool
	// Formats lists the supported image formats.
	Formats []ImageFormat
	// Transmissions lists the supported transmission methods.
	Transmissions []TransmitMethod
	// Animation is true if frame-based animation is supported.
	Animation bool
	// Unicode is true if Unicode placeholder mode is supported.
	Unicode bool
}

// HasKGP returns true if the terminal has any KGP support.
func (c Capabilities) HasKGP() bool {
	return c.KGP
}

// SupportsFormat returns true if the given format is supported.
func (c Capabilities) SupportsFormat(f ImageFormat) bool {
	for _, format := range c.Formats {
		if format == f {
			return true
		}
	}
	return false
}

// SupportsTransmission returns true if the given transmission method is supported.
func (c Capabilities) SupportsTransmission(m TransmitMethod) bool {
	for _, tm := range c.Transmissions {
		if tm == m {
			return true
		}
	}
	return false
}

// --- Preset profiles for known terminals ---

// KittyCapabilities returns the full capability set for Kitty (v0.19.0+).
func KittyCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG, ImageRGBA, ImageRGB},
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile, TransmitFromTempFile, TransmitFromSharedMem},
		Animation:     true,
		Unicode:       true,
	}
}

// WezTermCapabilities returns the capability set for WezTerm.
func WezTermCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG, ImageRGBA, ImageRGB},
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile, TransmitFromTempFile},
		Animation:     false,
		Unicode:       false,
	}
}

// KonsoleCapabilities returns the capability set for Konsole.
func KonsoleCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
		Animation:     false,
		Unicode:       false,
	}
}

// GhosttyCapabilities returns the capability set for Ghostty.
func GhosttyCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG, ImageRGBA, ImageRGB},
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile, TransmitFromTempFile, TransmitFromSharedMem},
		Animation:     true,
		Unicode:       true,
	}
}

// FootCapabilities returns the capability set for foot.
func FootCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG, ImageRGBA, ImageRGB},
		Transmissions: []TransmitMethod{TransmitDirect, TransmitFromFile, TransmitFromTempFile},
		Animation:     false,
		Unicode:       false,
	}
}

// ITermCapabilities returns the capability set for iTerm2 (v3.5+).
func ITermCapabilities() Capabilities {
	return Capabilities{
		KGP:           true,
		Formats:       []ImageFormat{ImagePNG},
		Transmissions: []TransmitMethod{TransmitDirect},
		Animation:     false,
		Unicode:       false,
	}
}

// NoKGPCapabilities returns an empty capability set for terminals without KGP.
func NoKGPCapabilities() Capabilities {
	return Capabilities{}
}

// detectCapabilities probes the terminal for KGP support.
// No-op if capabilities were set via WithCapabilities.
func (a *App) detectCapabilities() {
	if a.capabilitiesSet {
		return
	}

	// Use preset capabilities for known terminals to avoid probe timeouts.
	switch os.Getenv("TERM_PROGRAM") {
	case "kitty":
		a.capabilities = KittyCapabilities()
		return
	case "WezTerm":
		a.capabilities = WezTermCapabilities()
		return
	case "konsole":
		a.capabilities = KonsoleCapabilities()
		return
	case "ghostty":
		a.capabilities = GhosttyCapabilities()
		return
	case "foot":
		a.capabilities = FootCapabilities()
		return
	case "iTerm.app":
		a.capabilities = ITermCapabilities()
		return
	}

	// Send the basic KGP support query
	resp, err := a.screen.QueryAndRead(QueryKGPSupport(), detectTimeout)
	if err != nil {
		a.capabilities = NoKGPCapabilities()
		return
	}

	parsed, err := ParseImageResponse(string(resp))
	if err != nil || !parsed.Success {
		a.capabilities = NoKGPCapabilities()
		return
	}

	// KGP is supported — probe individual features
	a.capabilities.KGP = true

	// Probe formats
	for _, fmt := range []ImageFormat{ImagePNG, ImageRGBA, ImageRGB} {
		resp, err := a.screen.QueryAndRead(QueryFormat(fmt), detectTimeout)
		if err == nil {
			if p, err := ParseImageResponse(string(resp)); err == nil && p.Success {
				a.capabilities.Formats = append(a.capabilities.Formats, fmt)
			}
		}
	}

	// Probe transmission methods
	for _, method := range []TransmitMethod{TransmitDirect, TransmitFromFile, TransmitFromTempFile, TransmitFromSharedMem} {
		resp, err := a.screen.QueryAndRead(QueryTransmitMethod(method), detectTimeout)
		if err == nil {
			if p, err := ParseImageResponse(string(resp)); err == nil && p.Success {
				a.capabilities.Transmissions = append(a.capabilities.Transmissions, method)
			}
		}
	}

	// Animation/Unicode: no direct KGP query exists for these features.
	// Infer from full format + transmission support, which indicates a fully
	// compliant terminal (e.g., Kitty). This heuristic may need refinement
	// as new terminals add partial KGP support.
	if len(a.capabilities.Formats) == 3 && len(a.capabilities.Transmissions) == 4 {
		a.capabilities.Animation = true
		a.capabilities.Unicode = true
	}
}

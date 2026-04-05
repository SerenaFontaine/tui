---
title: KGP Integration
weight: 20
---

TUI provides complete integration with the [Kitty Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/) via the [kgp](https://github.com/SerenaFontaine/kgp) package. All KGP features are exposed through idiomatic Go APIs.

## Architecture

TUI integrates KGP at the buffer level. During `Render`, components add `ImagePlacement` values to the buffer. After diff-based text rendering, the framework encodes and flushes all image placements as KGP escape sequences.

```
Component.Render()
    └─► Buffer
         ├── Cell grid (text)  ──► Diff ──► ANSI sequences
         └── ImagePlacement[]  ──► KGP  ──► Kitty escape sequences
```

## Protocol Overview

KGP commands use APC (Application Program Command) escape sequences:

```
ESC_G<control-data>;<payload>ESC\
```

- **Control data**: Comma-separated `key=value` pairs (action, format, dimensions, IDs)
- **Payload**: Base64-encoded binary data

## Actions

| Action | Code | TUI Integration |
|--------|------|-----------------|
| Transmit | `t` | `TransmitImageWithID`, `ImageManager.Transmit` |
| Transmit & Display | `T` | `ImagePlacement.Encode`, widget `Image` |
| Put | `p` | `ImagePlacement.WithImageID` |
| Delete | `d` | `DeleteImageByID`, `DeleteAllImages`, etc. |
| Frame | `f` | `Animation.Encode` |
| Animate | `a` | `PlayOnce`, `PlayLoop`, `StopAnimation`, etc. |
| Compose | `c` | `ComposeFrames` |
| Query | `q` | `QueryKGPSupport`, `QueryFormat` |

## Image Formats

| Format | KGP Code | TUI Constant | Notes |
|--------|----------|--------------|-------|
| PNG | 100 | `ImagePNG` | Embeds dimensions, already compressed |
| RGBA | 32 | `ImageRGBA` | 4 bytes/pixel, requires dimensions |
| RGB | 24 | `ImageRGB` | 3 bytes/pixel, requires dimensions |

## Transmission Methods

| Method | KGP Code | TUI Constant | Best For |
|--------|----------|--------------|----------|
| Direct | `d` | `TransmitDirect` | Small images, remote sessions |
| File | `f` | `TransmitFromFile` | Large local images |
| Temp File | `t` | `TransmitFromTempFile` | One-shot images (auto-deleted) |
| Shared Memory | `s` | `TransmitFromSharedMem` | Highest performance, local only |

## Chunked Transmission

Large images are automatically chunked into 4096-byte segments. Each chunk is a separate escape sequence with a continuation flag. TUI handles this transparently via `EncodeChunked`.

## Response Handling

Terminals respond to KGP commands with status messages. Use `ParseImageResponse` to parse them:

```go
resp, err := tui.ParseImageResponse(responseString)
if err != nil {
    // Parse error
}
if resp.Success {
    // Command accepted
} else {
    fmt.Printf("Error %s: %s\n", resp.ErrorCode, resp.Message)
}
```

## Coverage

TUI exposes 100% of the kgp package's API surface:

- All 8 builder types (Transmit, Put, Delete, Frame, Animate, Compose, Query, Command)
- All transmission methods and formats
- All 20 delete modes with free variants
- Complete animation control (frames, playback, composition)
- Response parsing (standard and strict)
- Pixel helpers (conversion, compression, solid colors)
- Low-level command construction for custom use cases

## References

- [Kitty Graphics Protocol Specification](https://sw.kovidgoyal.net/kitty/graphics-protocol/)
- [kgp Go Bindings](https://github.com/SerenaFontaine/kgp)
- [kgp Documentation](https://serenafontaine.github.io/kgp/)

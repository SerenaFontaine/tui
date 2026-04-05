---
title: Terminal Support
weight: 50
---

TUI works in any terminal that supports ANSI escape sequences. KGP image features require additional protocol support.

## Text Features

All text-based features (widgets, layout, styling, events) work in any modern terminal:

| Feature | Requirement |
|---------|-------------|
| Basic rendering | ANSI escape sequences |
| 16 colors | ANSI color support |
| 256 colors | 256-color mode |
| True color (24-bit) | True color support |
| Mouse tracking | SGR mouse mode |
| Alternate screen | xterm alternate screen |

## KGP Image Features

| Terminal | Support | Notes |
|----------|---------|-------|
| **Kitty** | Full | Version 0.19.0+ |
| **WezTerm** | Partial | Core features supported |
| **Konsole** | Experimental | KDE Konsole |

## Verify KGP Support

```go
screen.Flush(tui.QueryKGPSupport())
// Terminal responds with OK if supported
```

Query specific capabilities:

```go
screen.Flush(tui.QueryFormat(tui.ImagePNG))
screen.Flush(tui.QueryTransmitMethod(tui.TransmitFromSharedMem))
```

## Performance Tips

1. **Transmission method**: Use shared memory for local apps, file transmission for large images, direct for small images or remote sessions
2. **Format**: PNG for smallest payload, RGBA/RGB with compression for speed
3. **Image reuse**: Transmit once with `ImageManager`, then place multiple times with `WithImageID`
4. **Compression**: Enable ZLIB for raw RGBA/RGB data with `WithCompression()`
5. **Chunking**: Handled automatically; no action needed
6. **Diff rendering**: The framework only updates changed cells; avoid full buffer clears in `Render`

## Suspend/Resume

TUI handles Ctrl+Z suspension gracefully:

1. Restores terminal to normal mode
2. Sends `SuspendMsg` to your component
3. Sends SIGTSTP to the process
4. On SIGCONT, re-enters raw mode and sends `ResumeMsg`
5. Forces a full redraw

## References

- [Kitty Graphics Protocol Specification](https://sw.kovidgoyal.net/kitty/graphics-protocol/)
- [Kitty Terminal](https://sw.kovidgoyal.net/kitty/)
- [WezTerm](https://wezfurlong.org/wezterm/)

---
title: TUI — Terminal User Interface Framework for Go
---

A fully featured terminal user interface framework for Go with first-class [Kitty Graphics Protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/) support via [kgp](https://github.com/SerenaFontaine/kgp).

## Features

- **Elm Architecture** — Init/Update/Render loop with immutable message passing
- **Buffer-based rendering** with automatic diffing for efficient terminal output
- **Complete KGP integration** including:
  - Image display (PNG, RGBA, RGB formats)
  - All transmission methods (direct, file, temp file, shared memory)
  - Animation support (frames, playback, composition)
  - Image management with lifecycle tracking
  - Response parsing
- **14 built-in widgets** — Text, Input, List, Table, Tabs, Viewport, Progress, Gauge, Spinner, Dialog, Sparkline, Scrollbar, Tree, Form
- **Flexible layout system** — Constraint-based splitting with Fixed, Flex, Percent, Min, Max
- **True color support** — ANSI 16, 256-color, and 24-bit RGB
- **Theme system** — Built-in themes with derived styles
- **Rich styled text** — Inline spans with mixed formatting
- **Focus management** — Named focus tracking
- **Mouse support** — SGR mouse tracking with all events
- **Suspend/Resume** — Ctrl+Z suspension with state restore

## Quick Links

- [Getting Started](/docs/getting-started/) — Installation and first steps
- [Examples](/docs/examples/) — Practical usage examples
- [API Reference](/docs/api/) — Exhaustive technical reference
- [KGP Integration](/docs/protocol/) — Kitty Graphics Protocol details

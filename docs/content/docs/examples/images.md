---
title: Images
weight: 5
---

## Display from image.Image

The simplest way to show an image:

```go
imgWidget := widget.NewImage(myImage)
imgWidget.Render(buf, area)
```

The widget converts the Go image to PNG and transmits it via the Kitty Graphics Protocol.

## Display from PNG Data

```go
data, _ := os.ReadFile("photo.png")
imgWidget := widget.NewImageFromPNG(data)
imgWidget.Render(buf, area)
```

## Display from Raw RGBA

```go
pixels := tui.SolidColorImageData(100, 100, 255, 0, 0, 255) // red square
imgWidget := widget.NewImageFromRGBA(pixels, 100, 100)
imgWidget.SetCompression(true) // ZLIB compress for efficiency
imgWidget.Render(buf, area)
```

## Display from File Path

```go
imgWidget := widget.NewImageFromFile("/path/to/image.png")
imgWidget.Render(buf, area)
```

## Pre-Transmitted Images

Transmit once and place multiple times for efficiency:

```go
mgr := tui.NewImageManager()

// Transmit (do this once)
seq, id := mgr.Transmit(myImage)
screen.Flush(seq)

// Place at different positions (cheap, just a reference)
placement1 := tui.NewImagePlacement(0, 0, 20, 10).WithImageID(id)
placement2 := tui.NewImagePlacement(25, 0, 10, 5).WithImageID(id)
buf.AddImage(placement1)
buf.AddImage(placement2)
```

## Image Placement Options

```go
placement := tui.NewImagePlacement(x, y, cols, rows).
    WithImage(img).
    WithCrop(srcX, srcY, srcW, srcH).  // crop to a region
    WithZIndex(-1).                      // below text
    WithCompression().                   // ZLIB compress
    WithCellOffset(pixelX, pixelY).      // sub-cell positioning
    WithVirtual().                        // Unicode placeholder mode
    WithCursorMove(false).               // don't advance cursor
    WithRelativeTo(parentID, parentPlaceID, offsetH, offsetV)
```

## Transmission Methods

```go
// Direct (embedded in escape sequence)
placement := tui.NewImagePlacement(0, 0, 20, 10).WithPNG(data)

// From file
placement := tui.NewImagePlacement(0, 0, 20, 10).
    WithFile("/path/to/image.png", tui.ImagePNG)

// From temp file (terminal auto-deletes)
placement := tui.NewImagePlacement(0, 0, 20, 10).
    WithTempFile("/tmp/tty-graphics-protocol-img.png", tui.ImagePNG)

// From shared memory (most efficient for local apps)
placement := tui.NewImagePlacement(0, 0, 20, 10).
    WithSharedMemory("my-shm", size, tui.ImageRGBA)
```

## Deleting Images

```go
screen.Flush(tui.DeleteImageByID(imageID))   // remove specific image
screen.Flush(tui.DeleteAllImages())           // remove all
screen.Flush(tui.DeleteAtCursor())            // remove at cursor
screen.Flush(tui.DeleteByCell(10, 5))         // remove at cell
screen.Flush(tui.DeleteByZIndex(-1))          // remove by layer
```

## Querying Support

```go
screen.Flush(tui.QueryKGPSupport())          // check protocol support
screen.Flush(tui.QueryFormat(tui.ImagePNG))  // check format support
```

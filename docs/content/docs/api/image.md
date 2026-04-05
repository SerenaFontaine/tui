---
title: Image
weight: 6
---

Complete KGP image integration. All functions return escape sequence strings ready to flush to the terminal.

## ImagePlacement

`ImagePlacement` describes an image to render at a terminal position. Use the chainable `With*` methods to configure.

### Creating Placements

```go
func NewImagePlacement(x, y, cols, rows int) ImagePlacement
```

### Data Sources

| Method | Description |
|--------|-------------|
| `WithPNG([]byte)` | PNG-encoded data |
| `WithRGBA([]byte, w, h)` | Raw 32-bit RGBA pixels |
| `WithRGB([]byte, w, h)` | Raw 24-bit RGB pixels |
| `WithImage(image.Image)` | Convert Go image to PNG |
| `WithImageRGBA(image.Image)` | Convert Go image to RGBA |
| `WithImageRGB(image.Image)` | Convert Go image to RGB |
| `WithImageID(uint32)` | Reference a pre-transmitted image |

### Transmission Methods

| Method | Description |
|--------|-------------|
| `WithFile(path, format)` | Read from filesystem |
| `WithFileOffset(path, format, offset, size)` | Read byte range from file |
| `WithTempFile(path, format)` | Temp file (terminal auto-deletes) |
| `WithSharedMemory(name, size, format)` | POSIX shared memory |

### Display Options

| Method | Description |
|--------|-------------|
| `WithCrop(x, y, w, h)` | Source rectangle for cropping |
| `WithCellOffset(x, y)` | Pixel offset within starting cell |
| `WithZIndex(int)` | Layering (negative = below text) |
| `WithCursorMove(bool)` | Advance cursor after placement |
| `WithCompression()` | ZLIB compress RGBA/RGB data |
| `WithVirtual()` | Unicode placeholder mode |
| `WithRelativeTo(parentID, placeID, offsetH, offsetV)` | Relative positioning |

### ID Management

| Method | Description |
|--------|-------------|
| `WithID(uint32)` | Set image ID |
| `WithImageNumber(uint32)` | Set client-side image number |
| `WithPlacementID(uint32)` | Set placement ID |
| `WithResponseSuppression(mode)` | Control terminal responses |

## Image Manager

`ImageManager` tracks transmitted images and manages lifecycle.

```go
mgr := tui.NewImageManager()
```

| Method | Description |
|--------|-------------|
| `NextID() uint32` | Get unique image ID |
| `Transmit(image.Image) (string, uint32)` | Transmit PNG, returns seq + ID |
| `TransmitRGBA(image.Image, compress) (string, uint32)` | Transmit RGBA |
| `TransmitFromFile(path, format, w, h) (string, uint32)` | Transmit from file |
| `Place(id, x, y, cols, rows) string` | Create placement |
| `Delete(id) string` | Delete and untrack |
| `DeleteAll() string` | Delete all tracked images |
| `Get(id) *ManagedImage` | Look up tracked image |

## Transmission Helpers

| Function | Description |
|----------|-------------|
| `TransmitImageWithID(image.Image, id)` | Transmit PNG with explicit ID |
| `TransmitImageRGBAWithID(image.Image, id, compress)` | Transmit RGBA |
| `TransmitImageRGBWithID(image.Image, id, compress)` | Transmit RGB |
| `TransmitFileWithID(path, format, w, h, id)` | Transmit from file |
| `TransmitTempFileWithID(path, format, w, h, id)` | Transmit from temp file |
| `TransmitSharedMemWithID(name, size, format, w, h, id)` | Transmit from shared memory |
| `TransmitImageAuto(image.Image)` | Auto-ID PNG transmission |
| `TransmitImageRGBAAuto(image.Image, compress)` | Auto-ID RGBA transmission |

## Delete Helpers

| Function | Description |
|----------|-------------|
| `DeleteImageByID(id)` | Remove specific image |
| `DeleteImageFree(id)` | Remove and free memory |
| `DeleteAllImages()` | Remove all images |
| `DeleteAllImagesFree()` | Remove all and free memory |
| `DeleteAtCursor()` / `DeleteAtCursorFree()` | At cursor position |
| `DeleteByCell(x, y)` / `DeleteByCellFree(x, y)` | At cell position |
| `DeleteByColumn(x)` / `DeleteByColumnFree(x)` | By column |
| `DeleteByRow(y)` / `DeleteByRowFree(y)` | By row |
| `DeleteByZIndex(z)` / `DeleteByZIndexFree(z)` | By z-index |
| `DeleteByIDRange(start, end)` / `DeleteByIDRangeFree(start, end)` | By ID range |
| `DeleteByPlacementID(imgID, placeID)` | By placement |
| `DeleteByImageNumber(num)` | By image number |

## Query Helpers

| Function | Description |
|----------|-------------|
| `QueryKGPSupport()` | Check protocol support |
| `QueryFormat(ImageFormat)` | Check format support |
| `QueryTransmitMethod(TransmitMethod)` | Check transmission method |
| `QueryWithTestData([]byte)` | Verify pipeline |

## Response Parsing

```go
resp, err := tui.ParseImageResponse(responseString)
// resp.Success, resp.ImageID, resp.ErrorCode, resp.Message
```

## Pixel Helpers

| Function | Description |
|----------|-------------|
| `ConvertToRGBA(image.Image) []byte` | Image to RGBA bytes |
| `ConvertToRGB(image.Image) []byte` | Image to RGB bytes |
| `ConvertToPNG(image.Image) ([]byte, error)` | Image to PNG bytes |
| `CompressImageData([]byte) ([]byte, error)` | ZLIB compression |
| `SolidColorImageData(w, h, r, g, b, a) []byte` | Solid color RGBA |
| `MakeRGBAColor(r, g, b, a) uint32` | Pack RGBA to uint32 |
| `ValidateTempPath(string) error` | Check temp file path |

## Low-Level Access

```go
cmd := tui.NewKGPCommand(tui.KGPActionTransmit)
```

Action constants: `KGPActionTransmit`, `KGPActionTransmitDisplay`, `KGPActionPut`, `KGPActionDelete`, `KGPActionFrame`, `KGPActionAnimate`, `KGPActionCompose`, `KGPActionQuery`.

Response suppression: `ResponseAll`, `ResponseErrorsOnly`, `ResponseOKOnly`.

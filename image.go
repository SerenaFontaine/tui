package tui

import (
	"image"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/SerenaFontaine/kgp"
)

// ImageFormat specifies how image data is encoded.
type ImageFormat int

const (
	ImagePNG  ImageFormat = iota // PNG encoded (includes dimensions)
	ImageRGBA                    // Raw 32-bit RGBA pixels
	ImageRGB                     // Raw 24-bit RGB pixels
)

// TransmitMethod specifies how image data is sent to the terminal.
type TransmitMethod int

const (
	TransmitDirect        TransmitMethod = iota // Embed data in escape sequence
	TransmitFromFile                            // Read from filesystem path
	TransmitFromTempFile                        // Read from temp file (terminal deletes after)
	TransmitFromSharedMem                       // Read from POSIX shared memory
)

// ImagePlacement describes an image to render at a specific terminal position.
type ImagePlacement struct {
	// Position in terminal cells (0-based)
	X, Y int
	// Display size in terminal cells
	Columns, Rows int

	// Image data — set via With* methods
	format    ImageFormat
	imageID   uint32
	imageNum  uint32
	placeID   uint32
	data      []byte
	imgWidth  int
	imgHeight int

	// Transmission method
	transmit   TransmitMethod
	filePath   string
	fileOffset int
	fileSize   int
	shmName    string
	shmSize    int

	// Source rectangle for cropping (zero means full image)
	SrcX, SrcY, SrcW, SrcH int
	// Pixel offset within the starting cell
	CellOffsetX, CellOffsetY int
	// Z-index for layering (negative = below text, positive = above)
	ZIndex int
	// Whether the cursor should move after placing the image
	CursorMove bool
	// Enable ZLIB compression for transmission
	Compress bool
	// Virtual placement (Unicode placeholder mode)
	Virtual bool
	// Response suppression mode
	SuppressResponse kgp.ResponseSuppression
	// Alt text for placeholder display when KGP is unavailable
	alt string
	// isAnimation marks this placement as part of an animation
	isAnimation bool

	// Relative positioning
	RelativeParentID      uint32
	RelativeParentPlaceID uint32
	RelativeOffsetH       int
	RelativeOffsetV       int
}

// NewImagePlacement creates a placement at the given position and size.
func NewImagePlacement(x, y, cols, rows int) ImagePlacement {
	return ImagePlacement{X: x, Y: y, Columns: cols, Rows: rows}
}

// --- Data source methods ---

// WithPNG sets PNG-encoded image data.
func (p ImagePlacement) WithPNG(data []byte) ImagePlacement {
	p.format = ImagePNG
	p.data = data
	p.transmit = TransmitDirect
	return p
}

// WithRGBA sets raw RGBA pixel data with dimensions.
func (p ImagePlacement) WithRGBA(data []byte, width, height int) ImagePlacement {
	p.format = ImageRGBA
	p.data = data
	p.imgWidth = width
	p.imgHeight = height
	p.transmit = TransmitDirect
	return p
}

// WithRGB sets raw RGB pixel data with dimensions.
func (p ImagePlacement) WithRGB(data []byte, width, height int) ImagePlacement {
	p.format = ImageRGB
	p.data = data
	p.imgWidth = width
	p.imgHeight = height
	p.transmit = TransmitDirect
	return p
}

// WithImage converts a Go image.Image to PNG for transmission.
func (p ImagePlacement) WithImage(img image.Image) ImagePlacement {
	data, err := kgp.ImageToPNG(img)
	if err != nil {
		return p
	}
	p.format = ImagePNG
	p.data = data
	p.transmit = TransmitDirect
	return p
}

// WithImageRGBA converts a Go image.Image to raw RGBA for transmission.
func (p ImagePlacement) WithImageRGBA(img image.Image) ImagePlacement {
	bounds := img.Bounds()
	p.format = ImageRGBA
	p.data = kgp.ImageToRGBA(img)
	p.imgWidth = bounds.Dx()
	p.imgHeight = bounds.Dy()
	p.transmit = TransmitDirect
	return p
}

// WithImageRGB converts a Go image.Image to raw RGB (no alpha) for transmission.
func (p ImagePlacement) WithImageRGB(img image.Image) ImagePlacement {
	bounds := img.Bounds()
	p.format = ImageRGB
	p.data = kgp.ImageToRGB(img)
	p.imgWidth = bounds.Dx()
	p.imgHeight = bounds.Dy()
	p.transmit = TransmitDirect
	return p
}

// WithImageID places an already-transmitted image by its ID.
func (p ImagePlacement) WithImageID(id uint32) ImagePlacement {
	p.imageID = id
	return p
}

// --- Transmission method ---

// WithFile transmits from a file path on disk.
func (p ImagePlacement) WithFile(path string, format ImageFormat) ImagePlacement {
	p.transmit = TransmitFromFile
	p.filePath = path
	p.format = format
	return p
}

// WithFileOffset transmits from a file with a byte offset and size.
func (p ImagePlacement) WithFileOffset(path string, format ImageFormat, offset, size int) ImagePlacement {
	p.transmit = TransmitFromFile
	p.filePath = path
	p.format = format
	p.fileOffset = offset
	p.fileSize = size
	return p
}

// WithTempFile transmits from a temp file (terminal auto-deletes it).
// Path must contain "tty-graphics-protocol" per the Kitty spec.
func (p ImagePlacement) WithTempFile(path string, format ImageFormat) ImagePlacement {
	p.transmit = TransmitFromTempFile
	p.filePath = path
	p.format = format
	return p
}

// WithSharedMemory transmits via POSIX shared memory.
func (p ImagePlacement) WithSharedMemory(name string, size int, format ImageFormat) ImagePlacement {
	p.transmit = TransmitFromSharedMem
	p.shmName = name
	p.shmSize = size
	p.format = format
	return p
}

// --- Placement options ---

// WithID sets the image ID for this placement.
func (p ImagePlacement) WithID(id uint32) ImagePlacement {
	p.imageID = id
	return p
}

// WithImageNumber sets the image number (client-side identifier).
func (p ImagePlacement) WithImageNumber(num uint32) ImagePlacement {
	p.imageNum = num
	return p
}

// WithPlacementID sets a unique ID for this specific placement.
func (p ImagePlacement) WithPlacementID(id uint32) ImagePlacement {
	p.placeID = id
	return p
}

// WithCrop sets the source rectangle for cropping a portion of the image.
func (p ImagePlacement) WithCrop(x, y, w, h int) ImagePlacement {
	p.SrcX = x
	p.SrcY = y
	p.SrcW = w
	p.SrcH = h
	return p
}

// WithCellOffset sets the pixel offset within the starting cell.
func (p ImagePlacement) WithCellOffset(x, y int) ImagePlacement {
	p.CellOffsetX = x
	p.CellOffsetY = y
	return p
}

// WithZIndex sets the z-index for layering (negative = below text).
func (p ImagePlacement) WithZIndex(z int) ImagePlacement {
	p.ZIndex = z
	return p
}

// WithCursorMove controls whether the cursor advances after placement.
func (p ImagePlacement) WithCursorMove(move bool) ImagePlacement {
	p.CursorMove = move
	return p
}

// WithCompression enables ZLIB compression for the image data.
func (p ImagePlacement) WithCompression() ImagePlacement {
	p.Compress = true
	return p
}

// WithVirtual creates a virtual (Unicode placeholder) placement.
func (p ImagePlacement) WithVirtual() ImagePlacement {
	p.Virtual = true
	return p
}

// WithRelativeTo positions this image relative to another placement.
func (p ImagePlacement) WithRelativeTo(parentID, parentPlaceID uint32, offsetH, offsetV int) ImagePlacement {
	p.RelativeParentID = parentID
	p.RelativeParentPlaceID = parentPlaceID
	p.RelativeOffsetH = offsetH
	p.RelativeOffsetV = offsetV
	return p
}

// WithResponseSuppression sets the response suppression mode.
func (p ImagePlacement) WithResponseSuppression(mode kgp.ResponseSuppression) ImagePlacement {
	p.SuppressResponse = mode
	return p
}

// WithAlt sets alt text shown in the placeholder when KGP is unavailable.
func (p ImagePlacement) WithAlt(text string) ImagePlacement {
	p.alt = text
	return p
}

// Alt returns the alt text for this placement.
func (p ImagePlacement) Alt() string {
	return p.alt
}

// WithAnimation marks this placement as an animation frame.
func (p ImagePlacement) WithAnimation() ImagePlacement {
	p.isAnimation = true
	return p
}

// IsAnimation returns true if this placement is an animation.
func (p ImagePlacement) IsAnimation() bool {
	return p.isAnimation
}

// Encode generates the KGP escape sequences to display this image.
func (p ImagePlacement) Encode() string {
	// If we have a pre-transmitted image ID with no new data, just place it
	if p.imageID > 0 && p.data == nil && p.transmit == TransmitDirect {
		return p.encodePut()
	}
	return p.encodeTransmitDisplay()
}

func (p ImagePlacement) encodePut() string {
	builder := kgp.NewPut(p.imageID)
	if p.imageNum > 0 {
		builder = builder.ImageNumber(p.imageNum)
	}
	if p.placeID > 0 {
		builder = builder.PlacementID(p.placeID)
	}
	if p.Columns > 0 || p.Rows > 0 {
		builder = builder.DisplaySize(p.Columns, p.Rows)
	}
	if p.SrcW > 0 && p.SrcH > 0 {
		builder = builder.SourceRect(p.SrcX, p.SrcY, p.SrcW, p.SrcH)
	}
	if p.CellOffsetX > 0 || p.CellOffsetY > 0 {
		builder = builder.CellOffset(p.CellOffsetX, p.CellOffsetY)
	}
	if p.ZIndex != 0 {
		builder = builder.ZIndex(p.ZIndex)
	}
	if !p.CursorMove {
		builder = builder.CursorMovement(false)
	}
	if p.Virtual {
		builder = builder.VirtualPlacement()
	}
	if p.RelativeParentID > 0 {
		builder = builder.RelativeTo(p.RelativeParentID, p.RelativeParentPlaceID, p.RelativeOffsetH, p.RelativeOffsetV)
	}
	if p.SuppressResponse != 0 {
		builder = builder.ResponseSuppression(p.SuppressResponse)
	}
	cmd := builder.Build()
	return cursorPosition(p.X, p.Y) + cmd.Encode()
}

func (p ImagePlacement) encodeTransmitDisplay() string {
	builder := kgp.NewTransmitDisplay()

	// Format
	switch p.format {
	case ImagePNG:
		builder = builder.Format(kgp.FormatPNG)
	case ImageRGBA:
		builder = builder.Format(kgp.FormatRGBA).Dimensions(p.imgWidth, p.imgHeight)
	case ImageRGB:
		builder = builder.Format(kgp.FormatRGB).Dimensions(p.imgWidth, p.imgHeight)
	}

	// Transmission method
	switch p.transmit {
	case TransmitDirect:
		if p.Compress && (p.format == ImageRGBA || p.format == ImageRGB) {
			compressed, err := kgp.CompressZlib(p.data)
			if err == nil {
				builder = builder.Compress().TransmitDirect(compressed)
			} else {
				builder = builder.TransmitDirect(p.data)
			}
		} else {
			builder = builder.TransmitDirect(p.data)
		}
	case TransmitFromFile:
		if p.fileOffset > 0 || p.fileSize > 0 {
			builder = builder.TransmitFileWithOffset(p.filePath, p.fileOffset, p.fileSize)
		} else {
			builder = builder.TransmitFile(p.filePath)
		}
	case TransmitFromTempFile:
		var err error
		builder, err = builder.TryTransmitTemp(p.filePath)
		if err != nil {
			return ""
		}
	case TransmitFromSharedMem:
		builder = builder.TransmitSharedMemory(p.shmName, p.shmSize)
	}

	// IDs
	if p.imageID > 0 {
		builder = builder.ImageID(p.imageID)
	}
	if p.imageNum > 0 {
		builder = builder.ImageNumber(p.imageNum)
	}
	if p.placeID > 0 {
		builder = builder.PlacementID(p.placeID)
	}

	// Display
	if p.Columns > 0 || p.Rows > 0 {
		builder = builder.DisplaySize(p.Columns, p.Rows)
	}
	if p.SrcW > 0 && p.SrcH > 0 {
		builder = builder.SourceRect(p.SrcX, p.SrcY, p.SrcW, p.SrcH)
	}
	if p.CellOffsetX > 0 || p.CellOffsetY > 0 {
		builder = builder.CellOffset(p.CellOffsetX, p.CellOffsetY)
	}
	if p.ZIndex != 0 {
		builder = builder.ZIndex(p.ZIndex)
	}
	if !p.CursorMove {
		builder = builder.CursorMovement(false)
	}
	if p.Virtual {
		builder = builder.VirtualPlacement()
	}
	if p.RelativeParentID > 0 {
		builder = builder.RelativeTo(p.RelativeParentID, p.RelativeParentPlaceID, p.RelativeOffsetH, p.RelativeOffsetV)
	}
	if p.SuppressResponse != 0 {
		builder = builder.ResponseSuppression(p.SuppressResponse)
	}

	cmd := builder.Build()
	chunks := cmd.EncodeChunked(4096)
	var sb strings.Builder
	sb.WriteString(cursorPosition(p.X, p.Y))
	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
	return sb.String()
}

// ---------- Image Manager ----------

// ImageManager tracks transmitted images and manages their lifecycle.
// It assigns unique IDs and handles cleanup.
type ImageManager struct {
	mu     sync.Mutex
	nextID atomic.Uint32
	images map[uint32]*ManagedImage
	screen *Screen
}

// ManagedImage is a tracked image in the manager.
type ManagedImage struct {
	ID     uint32
	Width  int
	Height int
	Format ImageFormat
}

// NewImageManager creates an image manager.
func NewImageManager() *ImageManager {
	m := &ImageManager{
		images: make(map[uint32]*ManagedImage),
	}
	m.nextID.Store(1)
	return m
}

// NextID returns a unique image ID.
func (m *ImageManager) NextID() uint32 {
	return m.nextID.Add(1) - 1
}

// Transmit uploads an image to the terminal and tracks it.
// Returns the escape sequence to send and the assigned ID.
func (m *ImageManager) Transmit(img image.Image) (string, uint32) {
	id := m.NextID()
	bounds := img.Bounds()

	m.mu.Lock()
	m.images[id] = &ManagedImage{
		ID:     id,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Format: ImagePNG,
	}
	m.mu.Unlock()

	return TransmitImageWithID(img, id), id
}

// TransmitRGBA uploads an image as RGBA data and tracks it.
func (m *ImageManager) TransmitRGBA(img image.Image, compress bool) (string, uint32) {
	id := m.NextID()
	bounds := img.Bounds()

	m.mu.Lock()
	m.images[id] = &ManagedImage{
		ID:     id,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Format: ImageRGBA,
	}
	m.mu.Unlock()

	return TransmitImageRGBAWithID(img, id, compress), id
}

// TransmitFromFile uploads image data from a file.
func (m *ImageManager) TransmitFromFile(path string, format ImageFormat, width, height int) (string, uint32) {
	id := m.NextID()

	m.mu.Lock()
	m.images[id] = &ManagedImage{
		ID:     id,
		Width:  width,
		Height: height,
		Format: format,
	}
	m.mu.Unlock()

	return TransmitFileWithID(path, format, width, height, id), id
}

// Place creates a placement for a previously transmitted image.
func (m *ImageManager) Place(imageID uint32, x, y, cols, rows int) string {
	return NewImagePlacement(x, y, cols, rows).WithImageID(imageID).Encode()
}

// Delete removes an image from terminal memory and the manager.
func (m *ImageManager) Delete(id uint32) string {
	m.mu.Lock()
	delete(m.images, id)
	m.mu.Unlock()
	return DeleteImageByID(id)
}

// DeleteAll removes all tracked images.
func (m *ImageManager) DeleteAll() string {
	m.mu.Lock()
	m.images = make(map[uint32]*ManagedImage)
	m.mu.Unlock()
	return DeleteAllImages()
}

// Get returns a managed image by ID.
func (m *ImageManager) Get(id uint32) *ManagedImage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.images[id]
}

// ---------- Transmission helpers ----------

// TransmitImageWithID transmits a Go image as PNG to the terminal.
func TransmitImageWithID(img image.Image, imageID uint32) string {
	cmd, err := kgp.TransmitImageWithID(img, imageID)
	if err != nil {
		return ""
	}
	return encodeChunked(cmd)
}

// TransmitImageRGBAWithID transmits as raw RGBA with optional compression.
func TransmitImageRGBAWithID(img image.Image, imageID uint32, compress bool) string {
	bounds := img.Bounds()
	data := kgp.ImageToRGBA(img)

	builder := kgp.NewTransmit().
		ImageID(imageID).
		Format(kgp.FormatRGBA).
		Dimensions(bounds.Dx(), bounds.Dy())

	if compress {
		compressed, err := kgp.CompressZlib(data)
		if err == nil {
			builder = builder.Compress().TransmitDirect(compressed)
		} else {
			builder = builder.TransmitDirect(data)
		}
	} else {
		builder = builder.TransmitDirect(data)
	}

	return encodeChunked(builder.Build())
}

// TransmitImageRGBWithID transmits as raw RGB (no alpha) with optional compression.
func TransmitImageRGBWithID(img image.Image, imageID uint32, compress bool) string {
	bounds := img.Bounds()
	data := kgp.ImageToRGB(img)

	builder := kgp.NewTransmit().
		ImageID(imageID).
		Format(kgp.FormatRGB).
		Dimensions(bounds.Dx(), bounds.Dy())

	if compress {
		compressed, err := kgp.CompressZlib(data)
		if err == nil {
			builder = builder.Compress().TransmitDirect(compressed)
		} else {
			builder = builder.TransmitDirect(data)
		}
	} else {
		builder = builder.TransmitDirect(data)
	}

	return encodeChunked(builder.Build())
}

// TransmitFileWithID transmits image data from a file path.
func TransmitFileWithID(path string, format ImageFormat, width, height int, imageID uint32) string {
	builder := kgp.NewTransmit().
		ImageID(imageID).
		TransmitFile(path)

	switch format {
	case ImagePNG:
		builder = builder.Format(kgp.FormatPNG)
	case ImageRGBA:
		builder = builder.Format(kgp.FormatRGBA).Dimensions(width, height)
	case ImageRGB:
		builder = builder.Format(kgp.FormatRGB).Dimensions(width, height)
	}

	return encodeChunked(builder.Build())
}

// TransmitTempFileWithID transmits from a temp file (terminal auto-deletes).
func TransmitTempFileWithID(path string, format ImageFormat, width, height int, imageID uint32) (string, error) {
	builder := kgp.NewTransmit().
		ImageID(imageID)

	switch format {
	case ImagePNG:
		builder = builder.Format(kgp.FormatPNG)
	case ImageRGBA:
		builder = builder.Format(kgp.FormatRGBA).Dimensions(width, height)
	case ImageRGB:
		builder = builder.Format(kgp.FormatRGB).Dimensions(width, height)
	}

	builder, err := builder.TryTransmitTemp(path)
	if err != nil {
		return "", err
	}

	return encodeChunked(builder.Build()), nil
}

// TransmitSharedMemWithID transmits via POSIX shared memory.
func TransmitSharedMemWithID(name string, size int, format ImageFormat, width, height int, imageID uint32) string {
	builder := kgp.NewTransmit().
		ImageID(imageID).
		TransmitSharedMemory(name, size)

	switch format {
	case ImagePNG:
		builder = builder.Format(kgp.FormatPNG)
	case ImageRGBA:
		builder = builder.Format(kgp.FormatRGBA).Dimensions(width, height)
	case ImageRGB:
		builder = builder.Format(kgp.FormatRGB).Dimensions(width, height)
	}

	return encodeChunked(builder.Build())
}

// ---------- Delete helpers ----------

// DeleteImageByID removes a specific image by its ID.
func DeleteImageByID(imageID uint32) string {
	return kgp.DeleteImage(imageID).Encode()
}

// DeleteImageFree removes an image and frees its memory.
func DeleteImageFree(imageID uint32) string {
	return kgp.DeleteImageFree(imageID).Encode()
}

// DeleteAllImages removes all images from the terminal.
func DeleteAllImages() string {
	return kgp.DeleteAll().Encode()
}

// DeleteAllImagesFree removes all images and frees memory.
func DeleteAllImagesFree() string {
	return kgp.DeleteAllFree().Encode()
}

// DeleteAtCursor removes images at the current cursor position.
func DeleteAtCursor() string {
	return kgp.DeleteAtCursor().Encode()
}

// DeleteAtCursorFree removes images at cursor and frees memory.
func DeleteAtCursorFree() string {
	return kgp.DeleteAtCursorFree().Encode()
}

// DeleteByCell removes images at a specific cell position.
func DeleteByCell(x, y int) string {
	return kgp.NewDelete(kgp.DeleteByCell).Cell(x, y).Build().Encode()
}

// DeleteByColumn removes images in a specific column.
func DeleteByColumn(x int) string {
	return kgp.NewDelete(kgp.DeleteByColumn).Column(x).Build().Encode()
}

// DeleteByRow removes images in a specific row.
func DeleteByRow(y int) string {
	return kgp.NewDelete(kgp.DeleteByRow).Row(y).Build().Encode()
}

// DeleteByZIndex removes images at a specific z-index.
func DeleteByZIndex(z int) string {
	return kgp.NewDelete(kgp.DeleteByZIndex).ZIndex(z).Build().Encode()
}

// DeleteByIDRange removes images within an ID range.
func DeleteByIDRange(startID, endID int) string {
	return kgp.NewDelete(kgp.DeleteByIDRange).IDRange(startID, endID).Build().Encode()
}

// DeleteByPlacementID removes a specific placement.
func DeleteByPlacementID(imageID, placementID uint32) string {
	return kgp.NewDelete(kgp.DeleteByPlacementID).
		ImageID(imageID).PlacementID(placementID).Build().Encode()
}

// DeleteByImageNumber removes an image by its image number.
func DeleteByImageNumber(num uint32) string {
	return kgp.NewDelete(kgp.DeleteByImageNumber).ImageNumber(num).Build().Encode()
}

// "Free" variants that also release terminal-side memory:

// DeleteByCellFree removes images at a cell and frees memory.
func DeleteByCellFree(x, y int) string {
	return kgp.NewDelete(kgp.DeleteByCellFree).Cell(x, y).Build().Encode()
}

// DeleteByColumnFree removes images in a column and frees memory.
func DeleteByColumnFree(x int) string {
	return kgp.NewDelete(kgp.DeleteByColumnFree).Column(x).Build().Encode()
}

// DeleteByRowFree removes images in a row and frees memory.
func DeleteByRowFree(y int) string {
	return kgp.NewDelete(kgp.DeleteByRowFree).Row(y).Build().Encode()
}

// DeleteByZIndexFree removes images at a z-index and frees memory.
func DeleteByZIndexFree(z int) string {
	return kgp.NewDelete(kgp.DeleteByZIndexFree).ZIndex(z).Build().Encode()
}

// DeleteByIDRangeFree removes images in an ID range and frees memory.
func DeleteByIDRangeFree(startID, endID int) string {
	return kgp.NewDelete(kgp.DeleteByIDRangeFree).IDRange(startID, endID).Build().Encode()
}

// DeleteByPlacementIDFree removes a placement and frees memory.
func DeleteByPlacementIDFree(imageID, placementID uint32) string {
	return kgp.NewDelete(kgp.DeleteByPlacementIDFree).
		ImageID(imageID).PlacementID(placementID).Build().Encode()
}

// DeleteByImageNumberFree removes by image number and frees memory.
func DeleteByImageNumberFree(num uint32) string {
	return kgp.NewDelete(kgp.DeleteByImageNumberFree).ImageNumber(num).Build().Encode()
}

// ---------- Query helpers ----------

// QueryKGPSupport returns the escape sequence to check KGP support.
func QueryKGPSupport() string {
	return kgp.QuerySupport().Encode()
}

// QueryFormat queries support for a specific image format.
func QueryFormat(format ImageFormat) string {
	builder := kgp.NewQuery()
	switch format {
	case ImagePNG:
		builder = builder.Format(kgp.FormatPNG)
	case ImageRGBA:
		builder = builder.Format(kgp.FormatRGBA)
	case ImageRGB:
		builder = builder.Format(kgp.FormatRGB)
	}
	return builder.Build().Encode()
}

// QueryTransmitMethod queries support for a specific transmission method.
func QueryTransmitMethod(method TransmitMethod) string {
	builder := kgp.NewQuery()
	switch method {
	case TransmitDirect:
		builder = builder.TransmitMedium(kgp.TransmitDirect)
	case TransmitFromFile:
		builder = builder.TransmitMedium(kgp.TransmitFile)
	case TransmitFromTempFile:
		builder = builder.TransmitMedium(kgp.TransmitTemp)
	case TransmitFromSharedMem:
		builder = builder.TransmitMedium(kgp.TransmitSharedMem)
	}
	return builder.Build().Encode()
}

// QueryWithTestData sends a query with test data to verify the pipeline.
func QueryWithTestData(data []byte) string {
	return kgp.NewQuery().TestData(data).Build().Encode()
}

// ---------- Response parsing ----------

// ImageResponse represents a parsed response from the terminal.
type ImageResponse = kgp.Response

// ParseImageResponse parses a KGP response string from the terminal.
func ParseImageResponse(response string) (*ImageResponse, error) {
	return kgp.ParseResponse(response)
}

// ParseImageResponseStrict is like ParseImageResponse with strict validation.
func ParseImageResponseStrict(response string) (*ImageResponse, error) {
	return kgp.ParseResponseStrict(response)
}

// ---------- Color/pixel helpers ----------

// MakeRGBAColor packs r, g, b, a into a 32-bit RGBA value.
func MakeRGBAColor(r, g, b, a uint8) uint32 {
	return kgp.CreateRGBAColor(r, g, b, a)
}

// SolidColorImageData generates raw RGBA bytes for a solid-color image.
func SolidColorImageData(width, height int, r, g, b, a uint8) []byte {
	return kgp.SolidColorImage(width, height, r, g, b, a)
}

// CompressImageData compresses raw pixel data using ZLIB.
func CompressImageData(data []byte) ([]byte, error) {
	return kgp.CompressZlib(data)
}

// ConvertToRGBA converts any Go image to raw RGBA bytes.
func ConvertToRGBA(img image.Image) []byte {
	return kgp.ImageToRGBA(img)
}

// ConvertToRGB converts any Go image to raw RGB bytes (no alpha).
func ConvertToRGB(img image.Image) []byte {
	return kgp.ImageToRGB(img)
}

// ConvertToPNG encodes a Go image as PNG bytes.
func ConvertToPNG(img image.Image) ([]byte, error) {
	return kgp.ImageToPNG(img)
}

// ValidateTempPath checks if a path is valid for temp file transmission.
func ValidateTempPath(path string) error {
	return kgp.ValidateTempPath(path)
}

// ---------- Convenience transmission (wrapping kgp helpers) ----------

// TransmitImageAuto transmits an image as PNG with auto-assigned ID.
func TransmitImageAuto(img image.Image) (string, error) {
	cmd, err := kgp.TransmitImage(img)
	if err != nil {
		return "", err
	}
	return encodeChunked(cmd), nil
}

// TransmitImageRGBAAuto transmits as RGBA with optional compression and auto ID.
func TransmitImageRGBAAuto(img image.Image, compress bool) (string, error) {
	cmd, err := kgp.TransmitImageRGBA(img, compress)
	if err != nil {
		return "", err
	}
	return encodeChunked(cmd), nil
}

// ---------- Re-exported KGP constants for advanced users ----------

// KGP Action constants for building raw commands.
var (
	KGPActionTransmit        = kgp.ActionTransmit
	KGPActionTransmitDisplay = kgp.ActionTransmitDisplay
	KGPActionPut             = kgp.ActionPut
	KGPActionDelete          = kgp.ActionDelete
	KGPActionFrame           = kgp.ActionFrame
	KGPActionAnimate         = kgp.ActionAnimate
	KGPActionCompose         = kgp.ActionCompose
	KGPActionQuery           = kgp.ActionQuery
)

// KGP compression constant.
var KGPCompressionZlib = kgp.CompressionZlib

// Response suppression modes.
var (
	ResponseAll        = kgp.ResponseAll
	ResponseErrorsOnly = kgp.ResponseErrorsOnly
	ResponseOKOnly     = kgp.ResponseOKOnly
)

// ErrInvalidTempPath is returned when a temp file path doesn't meet KGP requirements.
var ErrInvalidTempPath = kgp.ErrInvalidTempPath

// ---------- Low-level command access ----------

// NewKGPCommand creates a raw KGP command for advanced use cases.
// Use KGPAction* constants for the action parameter.
func NewKGPCommand(action kgp.Action) *kgp.Command {
	return kgp.NewCommand(action)
}

func encodeChunked(cmd *kgp.Command) string {
	chunks := cmd.EncodeChunked(4096)
	var sb strings.Builder
	for _, chunk := range chunks {
		sb.WriteString(chunk)
	}
	return sb.String()
}

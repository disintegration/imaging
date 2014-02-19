package imaging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"code.google.com/p/go.image/bmp"
	"code.google.com/p/go.image/tiff"
)

// Open loads an image from file
func Open(filename string) (img image.Image, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return OpenWithReader(file)
}

// Loads an image from io.Reader
func OpenWithReader(reader io.Reader) (img image.Image, err error) {
	img, _, err = image.Decode(reader)
	if err != nil {
		return
	}

	img = toNRGBA(img)
	return
}

// Save saves the image to file with the specified filename.
// The format is determined from the filename extension: "jpg" (or "jpeg"), "png", "tif" (or "tiff") and "bmp" are supported.
func Save(img image.Image, filename string) (err error) {
	format := strings.ToLower(filepath.Ext(filename))
	m, err := regexp.MatchString(`^\.(jpg|jpeg|png|tif|tiff|bmp)$`, format)
	if err != nil || !m {
		err = fmt.Errorf(`imaging: unsupported image format: "%s"`, format)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return SaveToWrite(img, format, file)
}

// Save image to io.Writer
func SaveToWrite(img image.Image, format string, writer io.Writer) (err error) {
	switch format {
	case ".jpg", ".jpeg":
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(writer, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 95})
		}

	case ".png":
		err = png.Encode(writer, img)
	case ".tif", ".tiff":
		err = tiff.Encode(writer, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case ".bmp":
		err = bmp.Encode(writer, img)
	default:
		err = fmt.Errorf(`imaging: unsupported image format: "%s"`, format)
	}
	return
}

// New creates a new image with the specified width and height, and fills it with the specified color.
func New(width, height int, fillColor color.Color) *image.NRGBA {
	if width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)

	if c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0 {
		return dst
	}

	cs := []uint8{c.R, c.G, c.B, c.A}

	// fill the first row
	for x := 0; x < width; x++ {
		copy(dst.Pix[x*4:(x+1)*4], cs)
	}
	// copy the first row to other rows
	for y := 1; y < height; y++ {
		copy(dst.Pix[y*dst.Stride:y*dst.Stride+width*4], dst.Pix[0:width*4])
	}

	return dst
}

// This function used internally to convert any image type to NRGBA if needed.
func toNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return Clone(img)
}

// clamp & round float64 to uint8 (0..255)
func clamp(v float64) uint8 {
	return uint8(math.Min(math.Max(v, 0.0), 255.0) + 0.5)
}

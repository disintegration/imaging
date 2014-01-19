package imaging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Open loads an image from file
func Open(filename string) (img image.Image, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	if err != nil {
		return
	}

	img = toNRGBA(img)
	return
}

// Save saves the image to file with the specified filename.
// The format is determined from the filename extension, "jpg" (or "jpeg") and "png" are supported.
func Save(img image.Image, filename string) (err error) {
	format := strings.ToLower(filepath.Ext(filename))
	if format != ".jpg" && format != ".jpeg" && format != ".png" {
		err = fmt.Errorf("unknown image format: %s", format)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

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
			err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
		}

	case ".png":
		err = png.Encode(file, img)
	}
	return
}

// New creates a new image with the specified width and height, and fills it with the specified color.
func New(width, height int, fillColor color.Color) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
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

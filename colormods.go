package imaging

import (
	"image"
	"math"
)

// ColorMapping produces color adjusted version of the image with given map function
func ColorMapping(img image.Image, mapping func (r, g, b, a uint8) (uint8, uint8, uint8, uint8)) *image.NRGBA {
	src := toNRGBA(img)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	parallel(width * height, func(partStart, partEnd int) {
		for k := partStart; k < partEnd; k++ {
			i := k*4
			r, g, b, a := src.Pix[i+0], src.Pix[i+1], src.Pix[i+2], src.Pix[i+3]
			dst.Pix[i+0], dst.Pix[i+1], dst.Pix[i+2], dst.Pix[i+3] = mapping(r, g, b, a)
		}
	})

	return dst
}

// Grayscale produces gamma corrected version of the image.
func Gamma(img image.Image, gamma float64) *image.NRGBA {
	inv := 1 / gamma
	correct := func (r uint8) uint8 {
		return uint8(255 * (math.Pow(float64(r) / 255.0, inv)))
	}
	return ColorMapping(img, func (r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return correct(r), correct(g), correct(b), a
	})
}

// Grayscale produces grayscale version of the image.
func Grayscale(img image.Image) *image.NRGBA {
	return ColorMapping(img, func (r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		f := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
		c := uint8(f + 0.5)
		return c, c, c, a
	})
}

// Invert produces inverted (negated) version of the image.
func Invert(img image.Image) *image.NRGBA {
	return ColorMapping(img, func (r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return 255-r, 255-g, 255-b, a
	})
}

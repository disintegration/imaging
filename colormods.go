package imaging

import (
	"image"
)

// Grayscale produces grayscale version of the image.
func Grayscale(img image.Image) *image.NRGBA {
	src := toNRGBA(img)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				j := y*dst.Stride + x*4
				r := float64(src.Pix[i+0])
				g := float64(src.Pix[i+1])
				b := float64(src.Pix[i+2])
				f := 0.299*r + 0.587*g + 0.114*b
				c := uint8(f + 0.5)
				dst.Pix[j+0] = c
				dst.Pix[j+1] = c
				dst.Pix[j+2] = c
				dst.Pix[j+3] = src.Pix[i+3]
			}
		}
	})

	return dst
}

// Invert produces inverted (negated) version of the image.
func Invert(img image.Image) *image.NRGBA {
	src := toNRGBA(img)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				j := y*dst.Stride + x*4
				dst.Pix[j+0] = 255 - src.Pix[i+0]
				dst.Pix[j+1] = 255 - src.Pix[i+1]
				dst.Pix[j+2] = 255 - src.Pix[i+2]
				dst.Pix[j+3] = src.Pix[i+3]
			}
		}
	})

	return dst
}

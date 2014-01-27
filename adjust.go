package imaging

import (
	"image"
	"image/color"
	"math"
)

func applyColorMapping(img image.Image, fn func(c color.NRGBA) color.NRGBA) *image.NRGBA {
	src := toNRGBA(img)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				j := y*dst.Stride + x*4

				r := src.Pix[i+0]
				g := src.Pix[i+1]
				b := src.Pix[i+2]
				a := src.Pix[i+3]

				c := fn(color.NRGBA{r, g, b, a})

				dst.Pix[j+0] = c.R
				dst.Pix[j+1] = c.G
				dst.Pix[j+2] = c.B
				dst.Pix[j+3] = c.A
			}
		}
	})

	return dst
}

// clamp & round float64 to uint8 (0..255)
func clamp(v float64) uint8 {
	return uint8(math.Min(math.Max(v, 0.0), 255.0) + 0.5)
}

// AdjustGamma performs a gamma correction on the image and returns the adjusted image.
// Gamma parameter must be positive. Gamma = 1.0 gives the original image.
// Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.
//
// Example:
//
//	dstImage = imaging.AdjustGamma(srcImage, 0.7)
//
func AdjustGamma(img image.Image, gamma float64) *image.NRGBA {
	e := 1.0 / math.Max(gamma, 0.0001)
	lut := make([]uint8, 256)

	for i := 0; i < 256; i++ {
		lut[i] = clamp(math.Pow(float64(i)/255.0, e) * 255.0)
	}

	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}

	return applyColorMapping(img, fn)
}

// AdjustContrast changes the contrast of the image using the percentage parameter and returns the adjusted image.
// The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid grey image.
//
// Examples:
//
//	dstImage = imaging.AdjustContrast(srcImage, -10) // decrease image contrast by 10%
//	dstImage = imaging.AdjustContrast(srcImage, 20) // increase image contrast by 20%
//
func AdjustContrast(img image.Image, percentage float64) *image.NRGBA {
	percentage = math.Min(math.Max(percentage, -100.0), 100.0)
	lut := make([]uint8, 256)

	v := (100.0 + percentage) / 100.0
	for i := 0; i < 256; i++ {
		if 0 <= v && v <= 1 {
			lut[i] = clamp((0.5 + (float64(i)/255.0-0.5)*v) * 255.0)
		} else if 1 < v && v < 2 {
			lut[i] = clamp((0.5 + (float64(i)/255.0-0.5)*(1/(2.0-v))) * 255.0)
		} else {
			lut[i] = uint8(float64(i)/255.0+0.5) * 255
		}
	}

	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}

	return applyColorMapping(img, fn)
}

// AdjustBrightness changes the brightness of the image using the percentage parameter and returns the adjusted image.
// The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
//
// Examples:
//
//	dstImage = imaging.AdjustBrightness(srcImage, -15) // decrease image brightness by 15%
//	dstImage = imaging.AdjustBrightness(srcImage, 10) // increase image brightness by 10%
//
func AdjustBrightness(img image.Image, percentage float64) *image.NRGBA {
	percentage = math.Min(math.Max(percentage, -100.0), 100.0)
	lut := make([]uint8, 256)

	shift := 255.0 * percentage / 100.0
	for i := 0; i < 256; i++ {
		lut[i] = clamp(float64(i) + shift)
	}

	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}

	return applyColorMapping(img, fn)
}

// Grayscale produces grayscale version of the image.
func Grayscale(img image.Image) *image.NRGBA {
	fn := func(c color.NRGBA) color.NRGBA {
		f := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
		y := uint8(f + 0.5)
		return color.NRGBA{y, y, y, c.A}
	}
	return applyColorMapping(img, fn)
}

// Invert produces inverted (negated) version of the image.
func Invert(img image.Image) *image.NRGBA {
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{255 - c.R, 255 - c.G, 255 - c.B, c.A}
	}
	return applyColorMapping(img, fn)
}

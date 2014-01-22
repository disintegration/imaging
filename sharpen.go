package imaging

import (
	"image"
)

// Sharpen produces a sharpened version of the image.
// Sigma parameter must be positive and indicates how much the image will be sharpened.
//
// Usage example:
//
//		dstImage := imaging.Sharpen(srcImage, 3.5)
//
func Sharpen(img image.Image, sigma float64) *image.NRGBA {
	if sigma <= 0 {
		// sigma parameter must be positive!
		return Clone(img)
	}

	src := toNRGBA(img)
	blurred := Blur(img, sigma)

	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				for j := 0; j < 4; j++ {
					k := i + j
					val := int(src.Pix[k]) + (int(src.Pix[k]) - int(blurred.Pix[k]))
					if val < 0 {
						val = 0
					} else if val > 255 {
						val = 255
					}
					dst.Pix[k] = uint8(val)
				}
			}
		}
	})

	return dst
}

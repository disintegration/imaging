package imaging

import (
	"image"
)

// Crop cuts out a rectangular region with the specified bounds
// from the image and returns the cropped image.
func Crop(img image.Image, rect image.Rectangle) *image.NRGBA {
	src := toNRGBA(img)
	srcRect := rect.Sub(img.Bounds().Min)
	sub := src.SubImage(srcRect)
	return Clone(sub) // New image Bounds().Min point will be (0, 0)
}

// CropCenter cuts out a rectangular region with the specified size
// from the center of the image and returns the cropped image.
func CropCenter(img image.Image, width, height int) *image.NRGBA {
	cropW, cropH := width, height

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	centerX := srcMinX + srcW/2
	centerY := srcMinY + srcH/2

	x0 := centerX - cropW/2
	y0 := centerY - cropH/2
	x1 := x0 + cropW
	y1 := y0 + cropH

	return Crop(img, image.Rect(x0, y0, x1, y1))
}

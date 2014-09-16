package imaging

import (
	"image"
	"math"
)

// Overlay draws the img image over the background image at given position
// and returns the combined image. Opacity parameter is the opacity of the img
// image layer, used to compose the images, it must be from 0.0 to 1.0.
//
// Usage examples:
//
//		// draw the sprite over the background at position (50, 50)
//		dstImage := imaging.Overlay(backgroundImage, spriteImage, image.Pt(50, 50), 1.0)
//
//		// blend two opaque images of the same size
//		dstImage := imaging.Overlay(imageOne, imageTwo, image.Pt(0, 0), 0.5)
//
func Overlay(background, img image.Image, pos image.Point, opacity float64) *image.NRGBA {
	opacity = math.Min(math.Max(opacity, 0.0), 1.0) // check: 0.0 <= opacity <= 1.0

	src := toNRGBA(img)
	dst := Clone(background)                    // cloned image bounds start at (0, 0)
	startPt := pos.Sub(background.Bounds().Min) // so we should translate start point
	endPt := startPt.Add(src.Bounds().Size())
	pasteBounds := image.Rectangle{startPt, endPt}

	if dst.Bounds().Overlaps(pasteBounds) {
		intersectBounds := dst.Bounds().Intersect(pasteBounds)

		for y := intersectBounds.Min.Y; y < intersectBounds.Max.Y; y++ {
			for x := intersectBounds.Min.X; x < intersectBounds.Max.X; x++ {
				i := y*dst.Stride + x*4

				srcX := x - pasteBounds.Min.X
				srcY := y - pasteBounds.Min.Y
				j := srcY*src.Stride + srcX*4

				a1 := float64(dst.Pix[i+3])
				a2 := float64(src.Pix[j+3])

				coef2 := opacity * a2 / 255.0
				coef1 := (1 - coef2) * a1 / 255.0
				coefSum := coef1 + coef2
				coef1 /= coefSum
				coef2 /= coefSum

				dst.Pix[i+0] = uint8(float64(dst.Pix[i+0])*coef1 + float64(src.Pix[j+0])*coef2)
				dst.Pix[i+1] = uint8(float64(dst.Pix[i+1])*coef1 + float64(src.Pix[j+1])*coef2)
				dst.Pix[i+2] = uint8(float64(dst.Pix[i+2])*coef1 + float64(src.Pix[j+2])*coef2)
				dst.Pix[i+3] = uint8(math.Min(a1+a2*opacity*(255.0-a1)/255.0, 255.0))
			}
		}
	}

	return dst
}

package imaging

import (
	"image"
)

// Paste pastes the img image to the background image at the specified position and returns the combined image.
func Paste(background, img image.Image, pos image.Point) *image.NRGBA {
	src := toNRGBA(img)
	dst := Clone(background)                    // cloned image bounds start at (0, 0)
	startPt := pos.Sub(background.Bounds().Min) // so we should translate start point
	endPt := startPt.Add(src.Bounds().Size())
	pasteBounds := image.Rectangle{startPt, endPt}

	if dst.Bounds().Overlaps(pasteBounds) {
		intersectBounds := dst.Bounds().Intersect(pasteBounds)

		rowSize := intersectBounds.Dx() * 4
		numRows := intersectBounds.Dy()

		srcStartX := intersectBounds.Min.X - pasteBounds.Min.X
		srcStartY := intersectBounds.Min.Y - pasteBounds.Min.Y

		i0 := dst.PixOffset(intersectBounds.Min.X, intersectBounds.Min.Y)
		j0 := src.PixOffset(srcStartX, srcStartY)

		di := dst.Stride
		dj := src.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}
	}

	return dst
}

// PasteCenter pastes the img image to the center of the background image and returns the combined image.
func PasteCenter(background, img image.Image) *image.NRGBA {
	bgBounds := background.Bounds()
	bgW := bgBounds.Dx()
	bgH := bgBounds.Dy()
	bgMinX := bgBounds.Min.X
	bgMinY := bgBounds.Min.Y

	centerX := bgMinX + bgW/2
	centerY := bgMinY + bgH/2

	x0 := centerX - img.Bounds().Dx()/2
	y0 := centerY - img.Bounds().Dy()/2

	return Paste(background, img, image.Pt(x0, y0))
}

// Package imaging provides basic image manipulation functions 
// (resize, rotate, crop, etc.) as well as simplified image loading and saving.
// 
// This package is based on the standard Go image package. All the image 
// manipulation functions provided by the package return a new *image.NRGBA.
package imaging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
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
	return
}

// Save saves the image to file with the specified filename. Format parameter can be "jpeg" or "png".
func Save(img image.Image, filename string, format string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	formatLower := strings.ToLower(format)
	switch formatLower {
	case "jpeg", "jpg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(file, img)
	default:
		err = fmt.Errorf("unknown image format: %s", format)
	}
	return
}

// New creates a new image with the specified width and height, and fills it with the specified color. 
func New(width, height int, fillColor color.Color) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)

	i0 := dst.PixOffset(0, 0)
	for y := 0; y < height; y, i0 = y+1, i0+dst.Stride {
		for x, i := 0, i0; x < width; x, i = x+1, i+4 {
			dst.Pix[i+0] = c.R
			dst.Pix[i+1] = c.G
			dst.Pix[i+2] = c.B
			dst.Pix[i+3] = c.A
		}
	}

	return dst
}

// This function converts any image type to *image.NRGBA for faster pixel access
// Optimized for most standard image types: NRGBA64, RGBA, RGBA64, YCbCr, Gray, Gray16
// If clone is true, the new image bounds will start at (0,0), also, a new copy
// will be created even if the source image's type is already NRGBA
func toNRGBA(src image.Image, clone bool) *image.NRGBA {
	if !clone {
		if src0, ok := src.(*image.NRGBA); ok {
			return src0
		}
	}

	srcBounds := src.Bounds()
	dstBounds := srcBounds

	// if we need a copy - translate Min point to (0, 0)
	if clone {
		dstBounds = dstBounds.Sub(dstBounds.Min)
	}

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}

			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}

			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}

// This function is used internally to check if the image type is *image.NRGBA
// If not - converts any image type to *image.NRGBA for faster pixel access
func convertToNRGBA(img image.Image) *image.NRGBA {
	// 'false' indicates that we don't need a new copy of img if it is already NRGBA
	// and that the new image's bounds will be equal the bounds of the source image
	return toNRGBA(img, false)
}

// Clone returns a copy of the img. New image bounds will be (0, 0)-(width, height).
func Clone(img image.Image) *image.NRGBA {
	// 'true' indicates that we need a new copy of img even if it is already NRGBA
	// and that the new image's bounds will start at point (0, 0)
	return toNRGBA(img, true)
}

// Crop cuts out a rectangular region with the specified bounds
// from the image and returns the cropped image.
func Crop(img image.Image, rect image.Rectangle) *image.NRGBA {
	src := convertToNRGBA(img)
	sub := src.SubImage(rect)
	return Clone(sub) // New image Bounds().Min point will be (0, 0)
}

// Crop cuts out a rectangular region with the specified size
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

// Paste pastes the src image to the img image at the specified position and returns the combined image.
func Paste(img, src image.Image, pos image.Point) *image.NRGBA {
	srcBounds := src.Bounds()
	src0 := convertToNRGBA(src)

	dst := Clone(img)                    // cloned image bounds start at (0, 0)
	startPt := pos.Sub(img.Bounds().Min) // so we should translate start point
	endPt := startPt.Add(srcBounds.Size())
	pasteBounds := image.Rectangle{startPt, endPt}

	if dst.Bounds().Overlaps(pasteBounds) {
		intersectBounds := dst.Bounds().Intersect(pasteBounds)

		rowSize := intersectBounds.Dx() * 4
		numRows := intersectBounds.Dy()

		srcStartX := intersectBounds.Min.X - pasteBounds.Min.X + srcBounds.Min.X
		srcStartY := intersectBounds.Min.Y - pasteBounds.Min.Y + srcBounds.Min.Y

		i0 := dst.PixOffset(intersectBounds.Min.X, intersectBounds.Min.Y)
		j0 := src0.PixOffset(srcStartX, srcStartY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	}

	return dst
}

// Paste pastes the src image to the center of the img image and returns the combined image.
func PasteCenter(img, src image.Image) *image.NRGBA {
	imgBounds := img.Bounds()
	imgW := imgBounds.Dx()
	imgH := imgBounds.Dy()
	imgMinX := imgBounds.Min.X
	imgMinY := imgBounds.Min.Y

	centerX := imgMinX + imgW/2
	centerY := imgMinY + imgH/2

	x0 := centerX - src.Bounds().Dx()/2
	y0 := centerY - src.Bounds().Dy()/2

	return Paste(img, src, image.Pt(x0, y0))
}

// Rotate90 rotates the image 90 degrees clockwise and returns the transformed image.
func Rotate90(img image.Image) *image.NRGBA {
	src := convertToNRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMinY := srcBounds.Min.Y

	dstW := srcBounds.Dy()
	dstH := srcBounds.Dx()
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMaxX - dstY - 1
			srcY := srcMinY + dstX

			srcOff := src.PixOffset(srcX, srcY)
			dstOff := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOff+0] = src.Pix[srcOff+0]
			dst.Pix[dstOff+1] = src.Pix[srcOff+1]
			dst.Pix[dstOff+2] = src.Pix[srcOff+2]
			dst.Pix[dstOff+3] = src.Pix[srcOff+3]
		}
	}

	return dst
}

// Rotate180 rotates the image 180 degrees clockwise and returns the transformed image.
func Rotate180(img image.Image) *image.NRGBA {
	src := convertToNRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMaxX - dstX - 1
			srcY := srcMaxY - dstY - 1

			srcOff := src.PixOffset(srcX, srcY)
			dstOff := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOff+0] = src.Pix[srcOff+0]
			dst.Pix[dstOff+1] = src.Pix[srcOff+1]
			dst.Pix[dstOff+2] = src.Pix[srcOff+2]
			dst.Pix[dstOff+3] = src.Pix[srcOff+3]
		}
	}

	return dst
}

// Rotate270 rotates the image 270 degrees clockwise and returns the transformed image.
func Rotate270(img image.Image) *image.NRGBA {
	src := convertToNRGBA(img)
	srcBounds := src.Bounds()
	srcMaxY := srcBounds.Max.Y
	srcMinX := srcBounds.Min.X

	dstW := srcBounds.Dy()
	dstH := srcBounds.Dx()
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMinX + dstY
			srcY := srcMaxY - dstX - 1

			srcOff := src.PixOffset(srcX, srcY)
			dstOff := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOff+0] = src.Pix[srcOff+0]
			dst.Pix[dstOff+1] = src.Pix[srcOff+1]
			dst.Pix[dstOff+2] = src.Pix[srcOff+2]
			dst.Pix[dstOff+3] = src.Pix[srcOff+3]
		}
	}

	return dst
}

// FlipH flips the image horizontally (from left to right) and returns the transformed image.
func FlipH(img image.Image) *image.NRGBA {
	src := convertToNRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMinY := srcBounds.Min.Y

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMaxX - dstX - 1
			srcY := srcMinY + dstY

			srcOff := src.PixOffset(srcX, srcY)
			dstOff := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOff+0] = src.Pix[srcOff+0]
			dst.Pix[dstOff+1] = src.Pix[srcOff+1]
			dst.Pix[dstOff+2] = src.Pix[srcOff+2]
			dst.Pix[dstOff+3] = src.Pix[srcOff+3]
		}
	}

	return dst
}

// FlipV flips the image vertically (from top to bottom) and returns the transformed image.
func FlipV(img image.Image) *image.NRGBA {
	src := convertToNRGBA(img)
	srcBounds := src.Bounds()
	srcMaxY := srcBounds.Max.Y
	srcMinX := srcBounds.Min.X

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMinX + dstX
			srcY := srcMaxY - dstY - 1

			srcOff := src.PixOffset(srcX, srcY)
			dstOff := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOff+0] = src.Pix[srcOff+0]
			dst.Pix[dstOff+1] = src.Pix[srcOff+1]
			dst.Pix[dstOff+2] = src.Pix[srcOff+2]
			dst.Pix[dstOff+3] = src.Pix[srcOff+3]
		}
	}

	return dst
}

// Antialias filter for Resize is a basic quadratic function.
func antialiasFilter(x float64) float64 {
	x = math.Abs(x)
	if x <= 1.0 {
		return x*x*(1.4*x-2.4) + 1
	}
	return 0
}

// Resize resizes the image to the specified width and height and returns the transformed image.
// If one of width or height is 0, the image aspect ratio is preserved.
func Resize(img image.Image, width, height int) *image.NRGBA {
	// Antialiased resize algorithm. The quality is good, especially at downsizing, 
	// but the speed is not too good, some optimisations are needed.

	dstW, dstH := width, height

	if dstW < 0 || dstH < 0 {
		return &image.NRGBA{}
	}
	if dstW == 0 && dstH == 0 {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}
	}

	// if new width or height is 0 then preserve aspect ratio, minimum 1px  
	if dstW == 0 {
		tmpW := float64(dstH) * float64(srcW) / float64(srcH)
		dstW = int(math.Max(1.0, math.Floor(tmpW+0.5)))
	}
	if dstH == 0 {
		tmpH := float64(dstW) * float64(srcH) / float64(srcW)
		dstH = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	src := convertToNRGBA(img)
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dy := float64(srcH) / float64(dstH)
	dx := float64(srcW) / float64(dstW)

	radiusX := dx / 2.0
	radiusY := dy / 2.0

	// increase the radius of antialiasing a little to produce smoother output image
	radiusX = math.Ceil(radiusX * 1.25)
	radiusY = math.Ceil(radiusY * 1.25)

	coefs := make([]float64, int(radiusY+1)*2*4)
	xvals := make([]float64, int(radiusX+1)*2*4)

	for dstY := 0; dstY < dstH; dstY++ {
		fy := float64(srcMinY) + (float64(dstY)+0.5)*dy - 0.5

		for dstX := 0; dstX < dstW; dstX++ {
			fx := float64(srcMinX) + (float64(dstX)+0.5)*dx - 0.5

			startX := int(math.Ceil(fx - radiusX))
			if startX < srcMinX {
				startX = srcMinX
			}
			endX := int(math.Floor(fx + radiusX))
			if endX > srcMaxX-1 {
				endX = srcMaxX - 1
			}

			startY := int(math.Ceil(fy - radiusY))
			if startY < srcMinY {
				startY = srcMinY
			}
			endY := int(math.Floor(fy + radiusY))
			if endY > srcMaxY-1 {
				endY = srcMaxY - 1
			}

			// cache y weight coefficients
			for y := startY; y <= endY; y++ {
				coefs[y-startY] = antialiasFilter((fy - float64(y)) / radiusY)
			}

			var k, sumk, r, g, b, a float64
			var i int

			// Calculate combined RGBA values for each column according to weights:
			for x := startX; x <= endX; x++ {

				r, g, b, a, sumk = 0.0, 0.0, 0.0, 0.0, 0.0
				for y := startY; y <= endY; y++ {
					k = coefs[y-startY]
					sumk += k
					i = src.PixOffset(x, y)
					r += float64(src.Pix[i+0]) * k
					g += float64(src.Pix[i+1]) * k
					b += float64(src.Pix[i+2]) * k
					a += float64(src.Pix[i+3]) * k
				}

				i = (x - startX) * 4
				xvals[i+0] = r / sumk
				xvals[i+1] = g / sumk
				xvals[i+2] = b / sumk
				xvals[i+3] = a / sumk

			}

			// calculate final rgba values
			r, g, b, a, sumk = 0.0, 0.0, 0.0, 0.0, 0.0
			for x := startX; x <= endX; x++ {
				k = antialiasFilter((fx - float64(x)) / radiusX)
				sumk += k
				i = (x - startX) * 4

				r += xvals[i+0] * k
				g += xvals[i+1] * k
				b += xvals[i+2] * k
				a += xvals[i+3] * k
			}

			r = math.Min(r/sumk, 255.0)
			g = math.Min(g/sumk, 255.0)
			b = math.Min(b/sumk, 255.0)
			a = math.Min(a/sumk, 255.0)

			i = dst.PixOffset(dstX, dstY)
			dst.Pix[i+0] = uint8(r + 0.5)
			dst.Pix[i+1] = uint8(g + 0.5)
			dst.Pix[i+2] = uint8(b + 0.5)
			dst.Pix[i+3] = uint8(a + 0.5)
		}
	}
	return dst
}

// Fit scales down the image to fit the specified maximum width and height and returns the transformed image.
func Fit(img image.Image, width, height int) *image.NRGBA {
	maxW, maxH := width, height

	if maxW <= 0 || maxH <= 0 {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}
	}

	if srcW <= maxW && srcH <= maxH {
		return Clone(img)
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	maxAspectRatio := float64(maxW) / float64(maxH)

	var newW, newH int
	if srcAspectRatio > maxAspectRatio {
		newW = maxW
		newH = int(float64(newW) / srcAspectRatio)
	} else {
		newH = maxH
		newW = int(float64(newH) * srcAspectRatio)
	}

	return Resize(img, newW, newH)
}

// Thumbnail scales the image up or down, crops it to the specified size and returns the transformed image.
func Thumbnail(img image.Image, width, height int) *image.NRGBA {
	thumbW, thumbH := width, height

	if thumbW <= 0 || thumbH <= 0 {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.NRGBA{}
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	thumbAspectRatio := float64(thumbW) / float64(thumbH)

	var tmp image.Image
	if srcAspectRatio > thumbAspectRatio {
		tmp = Resize(img, 0, thumbH)
	} else {
		tmp = Resize(img, thumbW, 0)
	}

	return CropCenter(tmp, thumbW, thumbH)
}

// Package imaging provides basic image manipulation functions
// (resize, rotate, flip, crop, etc.) as well as simplified image loading and saving.
// 
// This package is based on the standard Go image package. All the image 
// manipulation functions provided by the package take any image type that 
// implements image.Image interface, and return a new image of 
// *image.NRGBA type (32 bit RGBA colors, not premultiplied by alpha).
//
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

// Overlay draws the source image over the background image at given position
// and returns the combined image. Opacity parameter is the opacity of the source
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
func Overlay(background, source image.Image, pos image.Point, opacity float64) *image.NRGBA {
	opacity = math.Min(math.Max(opacity, 0.0), 1.0) // check: 0.0 <= opacity <= 1.0

	src := convertToNRGBA(source)
	srcBounds := src.Bounds()

	dst := Clone(background)                    // cloned image bounds start at (0, 0)
	startPt := pos.Sub(background.Bounds().Min) // so we should translate start point
	endPt := startPt.Add(srcBounds.Size())
	pasteBounds := image.Rectangle{startPt, endPt}

	if dst.Bounds().Overlaps(pasteBounds) {
		intersectBounds := dst.Bounds().Intersect(pasteBounds)

		for y := intersectBounds.Min.Y; y < intersectBounds.Max.Y; y++ {
			for x := intersectBounds.Min.X; x < intersectBounds.Max.X; x++ {
				i := dst.PixOffset(x, y)
				srcX := x - pasteBounds.Min.X + srcBounds.Min.X
				srcY := y - pasteBounds.Min.Y + srcBounds.Min.Y
				j := src.PixOffset(srcX, srcY)

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

// Resize resizes the image to the specified width and height using the specified resampling
// filter and returns the transformed image. If one of width or height is 0, the image aspect
// ratio is preserved.
//
// Supported resample filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
// CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
//
// Usage example:
//
//		dstImage := imaging.Resize(srcImage, 800, 600, imaging.Lanczos)
//
func Resize(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
	if filter.Support <= 0.0 { // nearest-neighbor special case
		return resizeNearest(img, width, height)
	}

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
	var tmp, dst *image.NRGBA

	// two-pass resize 
	if srcW != dstW {
		tmp = resizeHorizontal(src, dstW, filter)
	} else {
		tmp = src
	}

	if srcH != dstH {
		dst = resizeVertical(tmp, dstH, filter)
	} else {
		dst = tmp
	}

	return dst
}

func resizeHorizontal(src *image.NRGBA, width int, filter ResampleFilter) *image.NRGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X

	dstW := width
	dstH := srcH

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dX := float64(srcW) / float64(dstW)
	scaleX := math.Max(dX, 1.0)
	rX := math.Ceil(scaleX * filter.Support)
	weights := make([]float64, int(rX+2)*2)

	for dstX := 0; dstX < dstW; dstX++ {
		fX := float64(srcMinX) + (float64(dstX)+0.5)*dX - 0.5

		startX := int(math.Ceil(fX - rX))
		if startX < srcMinX {
			startX = srcMinX
		}
		endX := int(math.Floor(fX + rX))
		if endX > srcMaxX-1 {
			endX = srcMaxX - 1
		}

		// cache weights
		weightSum := 0.0
		for x := startX; x <= endX; x++ {
			w := filter.Kernel((float64(x) - fX) / scaleX)
			weightSum += w
			weights[x-startX] = w
		}

		for dstY := 0; dstY < dstH; dstY++ {
			srcY := srcMinY + dstY

			r, g, b, a := 0.0, 0.0, 0.0, 0.0
			for x := startX; x <= endX; x++ {
				weight := weights[x-startX]
				i := src.PixOffset(x, srcY)
				r += float64(src.Pix[i+0]) * weight
				g += float64(src.Pix[i+1]) * weight
				b += float64(src.Pix[i+2]) * weight
				a += float64(src.Pix[i+3]) * weight
			}

			r = math.Min(math.Max(r/weightSum, 0.0), 255.0)
			g = math.Min(math.Max(g/weightSum, 0.0), 255.0)
			b = math.Min(math.Max(b/weightSum, 0.0), 255.0)
			a = math.Min(math.Max(a/weightSum, 0.0), 255.0)

			j := dst.PixOffset(dstX, dstY)
			dst.Pix[j+0] = uint8(r + 0.5)
			dst.Pix[j+1] = uint8(g + 0.5)
			dst.Pix[j+2] = uint8(b + 0.5)
			dst.Pix[j+3] = uint8(a + 0.5)
		}
	}

	return dst
}

func resizeVertical(src *image.NRGBA, height int, filter ResampleFilter) *image.NRGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxY := srcBounds.Max.Y

	dstW := srcW
	dstH := height

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dY := float64(srcH) / float64(dstH)
	scaleY := math.Max(dY, 1.0)
	rY := math.Ceil(scaleY * filter.Support)
	weights := make([]float64, int(rY+2)*2)

	for dstY := 0; dstY < dstH; dstY++ {
		fY := float64(srcMinY) + (float64(dstY)+0.5)*dY - 0.5

		startY := int(math.Ceil(fY - rY))
		if startY < srcMinY {
			startY = srcMinY
		}
		endY := int(math.Floor(fY + rY))
		if endY > srcMaxY-1 {
			endY = srcMaxY - 1
		}

		// cache weights
		weightSum := 0.0
		for y := startY; y <= endY; y++ {
			w := filter.Kernel((float64(y) - fY) / scaleY)
			weightSum += w
			weights[y-startY] = w
		}

		for dstX := 0; dstX < dstW; dstX++ {
			srcX := srcMinX + dstX

			r, g, b, a := 0.0, 0.0, 0.0, 0.0
			for y := startY; y <= endY; y++ {
				weight := weights[y-startY]
				i := src.PixOffset(srcX, y)
				r += float64(src.Pix[i+0]) * weight
				g += float64(src.Pix[i+1]) * weight
				b += float64(src.Pix[i+2]) * weight
				a += float64(src.Pix[i+3]) * weight
			}

			r = math.Min(math.Max(r/weightSum, 0.0), 255.0)
			g = math.Min(math.Max(g/weightSum, 0.0), 255.0)
			b = math.Min(math.Max(b/weightSum, 0.0), 255.0)
			a = math.Min(math.Max(a/weightSum, 0.0), 255.0)

			j := dst.PixOffset(dstX, dstY)
			dst.Pix[j+0] = uint8(r + 0.5)
			dst.Pix[j+1] = uint8(g + 0.5)
			dst.Pix[j+2] = uint8(b + 0.5)
			dst.Pix[j+3] = uint8(a + 0.5)
		}
	}

	return dst
}

// fast nearest-neighbor resize, no filtering
func resizeNearest(img image.Image, width, height int) *image.NRGBA {
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

	dx := float64(srcW) / float64(dstW)
	dy := float64(srcH) / float64(dstH)

	for dstY := 0; dstY < dstH; dstY++ {
		fy := float64(srcMinY) + (float64(dstY)+0.5)*dy - 0.5

		for dstX := 0; dstX < dstW; dstX++ {
			fx := float64(srcMinX) + (float64(dstX)+0.5)*dx - 0.5

			srcX := int(math.Min(math.Max(math.Floor(fx+0.5), float64(srcMinX)), float64(srcMaxX)))
			srcY := int(math.Min(math.Max(math.Floor(fy+0.5), float64(srcMinY)), float64(srcMaxY)))

			srcOffset := src.PixOffset(srcX, srcY)
			dstOffset := dst.PixOffset(dstX, dstY)

			dst.Pix[dstOffset+0] = src.Pix[srcOffset+0]
			dst.Pix[dstOffset+1] = src.Pix[srcOffset+1]
			dst.Pix[dstOffset+2] = src.Pix[srcOffset+2]
			dst.Pix[dstOffset+3] = src.Pix[srcOffset+3]
		}
	}

	return dst
}

// Fit scales down the image using the specified resample filter to fit the specified 
// maximum width and height and returns the transformed image.
//
// Supported resample filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
// CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
//
// Usage example:
//
//		dstImage := imaging.Fit(srcImage, 800, 600, imaging.Lanczos)
//
func Fit(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
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

	return Resize(img, newW, newH, filter)
}

// Thumbnail scales the image up or down using the specified resample filter, crops it 
// to the specified width and hight and returns the transformed image.
//
// Supported resample filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
// CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
//
// Usage example:
//
//		dstImage := imaging.Fit(srcImage, 100, 100, imaging.Lanczos)
//
func Thumbnail(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
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
		tmp = Resize(img, 0, thumbH, filter)
	} else {
		tmp = Resize(img, thumbW, 0, filter)
	}

	return CropCenter(tmp, thumbW, thumbH)
}

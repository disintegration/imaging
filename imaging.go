// Package imaging provides basic image manipulation functions as well as 
// simplified image loading and saving.
// 
// This package is based on the standard Go image package. All the image 
// manipulation functions provided by the package return a new draw.Image.
// Currently, *image.RGBA type is used to build new images.
package imaging

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
func New(width, height int, fillColor color.Color) draw.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	unf := image.NewUniform(fillColor)
	draw.Draw(dst, dst.Bounds(), unf, image.ZP, draw.Src)
	return dst
}

// Copy returns a copy of the img. New image bounds will be (0, 0)-(width, height).
func Copy(img image.Image) draw.Image {
	imgBounds := img.Bounds()
	newBounds := imgBounds.Sub(imgBounds.Min) // new image bounds start at (0, 0)
	dst := image.NewRGBA(newBounds)
	draw.Draw(dst, newBounds, img, imgBounds.Min, draw.Src)
	return dst
}

// This function is used internally to check if the image type is image.RGBA
// If not - converts any image type to image.RGBA for faster pixel access
func convertToRGBA(src image.Image) *image.RGBA {
	var dst *image.RGBA

	switch src.(type) {
	case *image.RGBA:
		dst = src.(*image.RGBA)
	default:
		b := src.Bounds()
		dst = image.NewRGBA(b) // converted image have the same bounds as a source
		draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	}

	return dst
}

// Crop cuts out a rectangular region with the specified bounds
// from the image and returns the cropped image.
func Crop(img image.Image, rect image.Rectangle) draw.Image {
	src := convertToRGBA(img)
	sub := src.SubImage(rect)
	return Copy(sub) // New image Bounds().Min point will be (0, 0)
}

// Crop cuts out a rectangular region with the specified size
// from the center of the image and returns the cropped image.
func CropCenter(img image.Image, width, height int) draw.Image {
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
func Paste(img, src image.Image, pos image.Point) draw.Image {
	dst := Copy(img)                     // copied image bounds start at (0, 0)
	startPt := pos.Sub(img.Bounds().Min) // so we should translate start point
	endPt := startPt.Add(src.Bounds().Size())
	draw.Draw(dst, image.Rectangle{startPt, endPt}, src, src.Bounds().Min, draw.Src)
	return dst
}

// Paste pastes the src image to the center of the img image and returns the combined image.
func PasteCenter(img, src image.Image) draw.Image {
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
func Rotate90(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMinY := srcBounds.Min.Y

	dstW := srcBounds.Dy()
	dstH := srcBounds.Dx()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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
func Rotate180(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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
func Rotate270(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMaxY := srcBounds.Max.Y
	srcMinX := srcBounds.Min.X

	dstW := srcBounds.Dy()
	dstH := srcBounds.Dx()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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
func FlipH(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMaxX := srcBounds.Max.X
	srcMinY := srcBounds.Min.Y

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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
func FlipV(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMaxY := srcBounds.Max.Y
	srcMinX := srcBounds.Min.X

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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
func Resize(img image.Image, width, height int) draw.Image {
	// Antialiased resize algorithm. The quality is good, especially at downsizing, 
	// but the speed is not too good, some optimisations are needed.

	dstW, dstH := width, height

	if dstW < 0 || dstH < 0 {
		return &image.RGBA{}
	}
	if dstW == 0 && dstH == 0 {
		return &image.RGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
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

	src := convertToRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

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

			// calculate combined rgba values for each column according to weights
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
func Fit(img image.Image, width, height int) draw.Image {
	maxW, maxH := width, height

	if maxW <= 0 || maxH <= 0 {
		return &image.RGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}

	if srcW <= maxW && srcH <= maxH {
		return Copy(img)
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
func Thumbnail(img image.Image, width, height int) draw.Image {
	thumbW, thumbH := width, height

	if thumbW <= 0 || thumbH <= 0 {
		return &image.RGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
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

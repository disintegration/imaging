package imaging

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// Loads image from file. Returns image.Image
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

// Saves image img to file with given filename. Format parameter may be either "jpeg" or "png"
func Save(img image.Image, filename string, format string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	formatLower := strings.ToLower(format)
	switch formatLower {
	case "jpeg", "jpg":
		err = jpeg.Encode(file, img, nil)
	case "png":
		err = png.Encode(file, img)
	default:
		err = fmt.Errorf("unknown image format: %s", format)
	}
	return
}

// Creates a new image with given size and fills it with given color. 
func New(width, height int, fillColor color.Color) draw.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	unf := image.NewUniform(fillColor)
	draw.Draw(dst, dst.Bounds(), unf, image.ZP, draw.Src)
	return dst
}

// Returns a copy of img
func Copy(img image.Image) draw.Image {
	b := img.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst
}

// This function is used internally to check if image type is image.RGBA
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

// Returns a copy of rectangular area of img
func Crop(img image.Image, rect image.Rectangle) draw.Image {
	src := convertToRGBA(img)
	sub := src.SubImage(rect)
	return Copy(sub)
}

// Returns a copy of rectangular area of given size from the center of img
func CropCenter(img image.Image, cropW, cropH int) draw.Image {
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

// Pastes image src to image img at given position. Returns resulting image.
func Paste(img, src image.Image, pos image.Point) draw.Image {
	dst := Copy(img)
	startPt := pos
	endPt := startPt.Add(src.Bounds().Size())
	draw.Draw(dst, image.Rectangle{startPt, endPt}, src, src.Bounds().Min, draw.Src)
	return dst
}

// Pastes image src to the center of image img. Returns resulting image.
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

// Rotates image img by 90 degrees clockwise
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

// Rotates image img by 180 degrees clockwise
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

// Rotates image img by 270 degrees clockwise
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

// Flips image img horizontally (left-to-right)
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

// Flips image img vertically (top-to-bottom)
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

// Simple image resize function. Will be removed later if not needed.
func resizeNearest(img image.Image, dstW, dstH int) draw.Image {
	if dstW <= 0 || dstH <= 0 {
		return &image.RGBA{}
	}

	src := convertToRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}

	dy := float64(srcH) / float64(dstH)
	dx := float64(srcW) / float64(dstW)

	for dstY := 0; dstY < dstH; dstY++ {
		srcY := int(float64(srcMinY) + float64(dstY)*dy)
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := int(float64(srcMinX) + float64(dstX)*dx)

			i := dst.PixOffset(dstX, dstY)
			j := src.PixOffset(srcX, srcY)

			dst.Pix[i+0], dst.Pix[i+1], dst.Pix[i+2], dst.Pix[i+3] = src.Pix[j+0], src.Pix[j+1], src.Pix[j+2], src.Pix[j+3]
		}
	}
	return dst
}

// Resizes image img to width=dstW and height=dstH.
// Bilinear interpolation algorithm used at the moment. This may change later.
func Resize(img image.Image, dstW, dstH int) draw.Image {
	if dstW <= 0 || dstH <= 0 {
		return &image.RGBA{}
	}
	if dstW == 1 || dstH == 1 {
		return resizeNearest(img, dstW, dstH)
	}

	src := convertToRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}
	if srcW == 1 || srcH == 1 {
		return resizeNearest(img, dstW, dstH)
	}

	dy := float64(srcH-1) / float64(dstH-1)
	dx := float64(srcW-1) / float64(dstW-1)

	var i int

	for dstY := 0; dstY < dstH; dstY++ {

		fy := float64(srcMinY) + float64(dstY)*dy
		srcY := int(fy)
		if srcY < srcMinY {
			srcY = srcMinY
		} else {
			if srcY >= srcMaxY-1 {
				srcY = srcMaxY - 2
			}
		}
		u := fy - float64(srcY)

		for dstX := 0; dstX < dstW; dstX++ {

			fx := float64(srcMinX) + float64(dstX)*dx
			srcX := int(fx)
			if srcX < srcMinX {
				srcX = srcMinX
			} else {
				if srcX >= srcMaxX-1 {
					srcX = srcMaxX - 2
				}
			}
			v := fx - float64(srcX)

			z1 := (1 - v) * (1 - u)
			z2 := v * (1 - u)
			z3 := v * u
			z4 := (1 - v) * u

			i = src.PixOffset(srcX, srcY)
			r1, g1, b1, a1 := src.Pix[i+0], src.Pix[i+1], src.Pix[i+2], src.Pix[i+3]

			i = src.PixOffset(srcX+1, srcY)
			r2, g2, b2, a2 := src.Pix[i+0], src.Pix[i+1], src.Pix[i+2], src.Pix[i+3]

			i = src.PixOffset(srcX+1, srcY+1)
			r3, g3, b3, a3 := src.Pix[i+0], src.Pix[i+1], src.Pix[i+2], src.Pix[i+3]

			i = src.PixOffset(srcX, srcY+1)
			r4, g4, b4, a4 := src.Pix[i+0], src.Pix[i+1], src.Pix[i+2], src.Pix[i+3]

			r := uint8(z1*float64(r1) + z2*float64(r2) + z3*float64(r3) + z4*float64(r4) + 0.5)
			g := uint8(z1*float64(g1) + z2*float64(g2) + z3*float64(g3) + z4*float64(g4) + 0.5)
			b := uint8(z1*float64(b1) + z2*float64(b2) + z3*float64(b3) + z4*float64(b4) + 0.5)
			a := uint8(z1*float64(a1) + z2*float64(a2) + z3*float64(a3) + z4*float64(a4) + 0.5)

			i = dst.PixOffset(dstX, dstY)
			dst.Pix[i+0], dst.Pix[i+1], dst.Pix[i+2], dst.Pix[i+3] = r, g, b, a
		}
	}

	return dst
}

// Scales image with given scale factor, keeps aspect ratio.
func Scale(img image.Image, scaleFactor float64) draw.Image {
	if scaleFactor <= 0.0 {
		return &image.RGBA{}
	}

	if scaleFactor == 1.0 {
		return Copy(img)
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}

	dstW := int(float64(srcW) * scaleFactor)
	dstH := int(float64(srcH) * scaleFactor)

	return Resize(img, dstW, dstH)
}

// Scales image  to given width, keeps aspect ratio.
func ScaleToWidth(img image.Image, dstW int) draw.Image {
	if dstW <= 0 {
		return &image.RGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}

	if dstW == srcW {
		return Copy(img)
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	dstH := int(float64(dstW) / srcAspectRatio)

	return Resize(img, dstW, dstH)
}

// Scales image  to given height, keeps aspect ratio.
func ScaleToHeight(img image.Image, dstH int) draw.Image {
	if dstH <= 0 {
		return &image.RGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW <= 0 || srcH <= 0 {
		return &image.RGBA{}
	}

	if dstH == srcH {
		return Copy(img)
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	dstW := int(float64(dstH) * srcAspectRatio)

	return Resize(img, dstW, dstH)
}

// Scales down image to fit given maximum width and height, keeps aspect ratio.
func Fit(img image.Image, maxW, maxH int) draw.Image {
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
		newH = int(srcAspectRatio / float64(newW))
	} else {
		newH = maxH
		newW = int(srcAspectRatio * float64(newH))
	}

	return Resize(img, newW, newH)
}

// Scales image up or down and crops to exact given size.
func Thumbnail(img image.Image, thumbW, thumbH int) draw.Image {
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
		tmp = ScaleToHeight(img, thumbH)
	} else {
		tmp = ScaleToWidth(img, thumbW)
	}

	return CropCenter(tmp, thumbW, thumbH)
}

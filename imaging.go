package imaging

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// Loads image from disk. Returns image.Image
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

// Saves image img to disk with given filename. Format parameter may be either "jpeg" or "png"
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

// Used internally to convert any image type to image.RGBA for faster pixel access
func convertToRGBA(src image.Image) *image.RGBA {
	var dst *image.RGBA

	switch src.(type) {
	case *image.RGBA:
		dst = src.(*image.RGBA)
	default:
		b := src.Bounds()
		dst = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	}

	return dst
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

// Returns a copy of img
func Copy(img image.Image) draw.Image {
	src := convertToRGBA(img)
	srcBounds := src.Bounds()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	dstW := srcBounds.Dx()
	dstH := srcBounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {

			srcX := srcMinX + dstX
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

// Simple image resize function. Will be removed later if not needed.
func resizeNearest(img image.Image, dstW, dstH int) draw.Image {
	src := convertToRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	dy := float64(srcH) / float64(dstH)
	dx := float64(srcW) / float64(dstW)

	for dstY := 0; dstY < dstH; dstY++ {
		srcY := int(float64(srcMinY) + float64(dstY)*dy)
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := int(float64(srcMinX) + float64(dstX)*dx)

			dst.Set(dstX, dstY, src.At(srcX, srcY))
		}
	}
	return dst
}

// Resizes image img to width=dstW and height=dstH.
// Bilinear interpolation algorithm used at the moment. This may change later.
func Resize(img image.Image, dstW, dstH int) draw.Image {
	// TODO: check src and dst image bounds, fix small images resizing (w or h <= 2px)

	src := convertToRGBA(img)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

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

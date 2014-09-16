package imaging

import (
	"image"
	"math"
)

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
	dstW, dstH := width, height

	if dstW < 0 || dstH < 0 {
		return &image.NRGBA{}
	}
	if dstW == 0 && dstH == 0 {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcW := src.Bounds().Max.X
	srcH := src.Bounds().Max.Y

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

	var dst *image.NRGBA

	if filter.Support <= 0.0 {
		// nearest-neighbor special case
		dst = resizeNearest(src, dstW, dstH)

	} else {
		// two-pass resize
		if srcW != dstW {
			dst = resizeHorizontal(src, dstW, filter)
		} else {
			dst = src
		}

		if srcH != dstH {
			dst = resizeVertical(dst, dstH, filter)
		}
	}

	return dst
}

func resizeHorizontal(src *image.NRGBA, width int, filter ResampleFilter) *image.NRGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Max.X
	srcH := srcBounds.Max.Y

	dstW := width
	dstH := srcH

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dX := float64(srcW) / float64(dstW)
	scaleX := math.Max(dX, 1.0)
	rX := math.Ceil(scaleX * filter.Support)

	parallel(dstW, func(partStart, partEnd int) {

		weights := make([]float64, int(rX+2)*2)

		for dstX := partStart; dstX < partEnd; dstX++ {

			fX := (float64(dstX)+0.5)*dX - 0.5

			startX := int(math.Ceil(fX - rX))
			if startX < 0 {
				startX = 0
			}
			endX := int(math.Floor(fX + rX))
			if endX > srcW-1 {
				endX = srcW - 1
			}

			// cache weights
			weightSum := 0.0
			for x := startX; x <= endX; x++ {
				w := filter.Kernel((float64(x) - fX) / scaleX)
				weightSum += w
				weights[x-startX] = w
			}

			for dstY := 0; dstY < dstH; dstY++ {
				r, g, b, a := 0.0, 0.0, 0.0, 0.0
				for x := startX; x <= endX; x++ {
					weight := weights[x-startX]
					i := dstY*src.Stride + x*4
					r += float64(src.Pix[i+0]) * weight
					g += float64(src.Pix[i+1]) * weight
					b += float64(src.Pix[i+2]) * weight
					a += float64(src.Pix[i+3]) * weight
				}

				r = math.Min(math.Max(r/weightSum, 0.0), 255.0)
				g = math.Min(math.Max(g/weightSum, 0.0), 255.0)
				b = math.Min(math.Max(b/weightSum, 0.0), 255.0)
				a = math.Min(math.Max(a/weightSum, 0.0), 255.0)

				j := dstY*dst.Stride + dstX*4
				dst.Pix[j+0] = uint8(r + 0.5)
				dst.Pix[j+1] = uint8(g + 0.5)
				dst.Pix[j+2] = uint8(b + 0.5)
				dst.Pix[j+3] = uint8(a + 0.5)
			}
		}

	})

	return dst
}

func resizeVertical(src *image.NRGBA, height int, filter ResampleFilter) *image.NRGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Max.X
	srcH := srcBounds.Max.Y

	dstW := srcW
	dstH := height

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dY := float64(srcH) / float64(dstH)
	scaleY := math.Max(dY, 1.0)
	rY := math.Ceil(scaleY * filter.Support)

	parallel(dstH, func(partStart, partEnd int) {

		weights := make([]float64, int(rY+2)*2)

		for dstY := partStart; dstY < partEnd; dstY++ {

			fY := (float64(dstY)+0.5)*dY - 0.5

			startY := int(math.Ceil(fY - rY))
			if startY < 0 {
				startY = 0
			}
			endY := int(math.Floor(fY + rY))
			if endY > srcH-1 {
				endY = srcH - 1
			}

			// cache weights
			weightSum := 0.0
			for y := startY; y <= endY; y++ {
				w := filter.Kernel((float64(y) - fY) / scaleY)
				weightSum += w
				weights[y-startY] = w
			}

			for dstX := 0; dstX < dstW; dstX++ {
				r, g, b, a := 0.0, 0.0, 0.0, 0.0
				for y := startY; y <= endY; y++ {
					weight := weights[y-startY]
					i := y*src.Stride + dstX*4
					r += float64(src.Pix[i+0]) * weight
					g += float64(src.Pix[i+1]) * weight
					b += float64(src.Pix[i+2]) * weight
					a += float64(src.Pix[i+3]) * weight
				}

				r = math.Min(math.Max(r/weightSum, 0.0), 255.0)
				g = math.Min(math.Max(g/weightSum, 0.0), 255.0)
				b = math.Min(math.Max(b/weightSum, 0.0), 255.0)
				a = math.Min(math.Max(a/weightSum, 0.0), 255.0)

				j := dstY*dst.Stride + dstX*4
				dst.Pix[j+0] = uint8(r + 0.5)
				dst.Pix[j+1] = uint8(g + 0.5)
				dst.Pix[j+2] = uint8(b + 0.5)
				dst.Pix[j+3] = uint8(a + 0.5)
			}
		}

	})

	return dst
}

// fast nearest-neighbor resize, no filtering
func resizeNearest(src *image.NRGBA, width, height int) *image.NRGBA {
	dstW, dstH := width, height

	srcBounds := src.Bounds()
	srcW := srcBounds.Max.X
	srcH := srcBounds.Max.Y

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dx := float64(srcW) / float64(dstW)
	dy := float64(srcH) / float64(dstH)

	parallel(dstH, func(partStart, partEnd int) {

		for dstY := partStart; dstY < partEnd; dstY++ {
			fy := (float64(dstY)+0.5)*dy - 0.5

			for dstX := 0; dstX < dstW; dstX++ {
				fx := (float64(dstX)+0.5)*dx - 0.5

				srcX := int(math.Min(math.Max(math.Floor(fx+0.5), 0.0), float64(srcW)))
				srcY := int(math.Min(math.Max(math.Floor(fy+0.5), 0.0), float64(srcH)))

				srcOff := srcY*src.Stride + srcX*4
				dstOff := dstY*dst.Stride + dstX*4

				copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
			}
		}

	})

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

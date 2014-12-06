package imaging

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func compareNRGBA(img1, img2 *image.NRGBA, delta int) bool {
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}

	if len(img1.Pix) != len(img2.Pix) {
		return false
	}

	for i := 0; i < len(img1.Pix); i++ {
		if absint(int(img1.Pix[i])-int(img2.Pix[i])) > delta {
			return false
		}
	}

	return true
}

func TestEncodeDecode(t *testing.T) {
	imgWithAlpha := image.NewNRGBA(image.Rect(0, 0, 3, 3))
	imgWithAlpha.Pix = []uint8{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138,
		244, 245, 246, 247, 248, 249, 250, 252, 252, 253, 254, 255,
	}

	imgWithoutAlpha := image.NewNRGBA(image.Rect(0, 0, 3, 3))
	imgWithoutAlpha.Pix = []uint8{
		0, 1, 2, 255, 4, 5, 6, 255, 8, 9, 10, 255,
		127, 128, 129, 255, 131, 132, 133, 255, 135, 136, 137, 255,
		244, 245, 246, 255, 248, 249, 250, 255, 252, 253, 254, 255,
	}

	for _, format := range []Format{JPEG, PNG, GIF, BMP, TIFF} {
		img := imgWithoutAlpha
		if format == PNG {
			img = imgWithAlpha
		}

		buf := &bytes.Buffer{}
		err := Encode(buf, img, format)
		if err != nil {
			t.Errorf("fail encoding format %s", format)
			continue
		}

		img2, err := Decode(buf)
		if err != nil {
			t.Errorf("fail decoding format %s", format)
			continue
		}
		img2cloned := Clone(img2)

		delta := 0
		if format == JPEG {
			delta = 3
		} else if format == GIF {
			delta = 16
		}

		if !compareNRGBA(img, img2cloned, delta) {
			t.Errorf("test [DecodeEncode %s] failed: %#v %#v", format, img, img2cloned)
			continue
		}
	}

	buf := &bytes.Buffer{}
	err := Encode(buf, imgWithAlpha, Format(100))
	if err != ErrUnsupportedFormat {
		t.Errorf("expected ErrUnsupportedFormat")
	}
}

func TestNew(t *testing.T) {
	td := []struct {
		desc      string
		w, h      int
		c         color.Color
		dstBounds image.Rectangle
		dstPix    []uint8
	}{
		{
			"New 1x1 black",
			1, 1,
			color.NRGBA{0, 0, 0, 0},
			image.Rect(0, 0, 1, 1),
			[]uint8{0x00, 0x00, 0x00, 0x00},
		},
		{
			"New 1x2 red",
			1, 2,
			color.NRGBA{255, 0, 0, 255},
			image.Rect(0, 0, 1, 2),
			[]uint8{0xff, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0xff},
		},
		{
			"New 2x1 white",
			2, 1,
			color.NRGBA{255, 255, 255, 255},
			image.Rect(0, 0, 2, 1),
			[]uint8{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
	}

	for _, d := range td {
		got := New(d.w, d.h, d.c)
		want := image.NewNRGBA(d.dstBounds)
		want.Pix = d.dstPix
		if !compareNRGBA(got, want, 0) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

func TestClone(t *testing.T) {
	td := []struct {
		desc string
		src  image.Image
		want *image.NRGBA
	}{
		{
			"Clone NRGBA",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone NRGBA64",
			&image.NRGBA64{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 8,
				Pix: []uint8{
					0x00, 0x00, 0x11, 0x11, 0x22, 0x22, 0x33, 0x33,
					0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee, 0xff, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone RGBA",
			&image.RGBA{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone RGBA64",
			&image.RGBA64{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 8,
				Pix: []uint8{
					0x00, 0x00, 0x11, 0x11, 0x22, 0x22, 0x33, 0x33,
					0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee, 0xff, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone Gray",
			&image.Gray{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 1,
				Pix:    []uint8{0x11, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x11, 0x11, 0x11, 0xff, 0xee, 0xee, 0xee, 0xff},
			},
		},
		{
			"Clone Gray16",
			&image.Gray16{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 2,
				Pix:    []uint8{0x11, 0x11, 0xee, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x11, 0x11, 0x11, 0xff, 0xee, 0xee, 0xee, 0xff},
			},
		},
		{
			"Clone Alpha",
			&image.Alpha{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 1,
				Pix:    []uint8{0x11, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0xff, 0xff, 0xff, 0x11, 0xff, 0xff, 0xff, 0xee},
			},
		},
		{
			"Clone YCbCr",
			&image.YCbCr{
				Rect:           image.Rect(-1, -1, 5, 0),
				SubsampleRatio: image.YCbCrSubsampleRatio444,
				YStride:        6,
				CStride:        6,
				Y:              []uint8{0x00, 0xff, 0x7f, 0x26, 0x4b, 0x0e},
				Cb:             []uint8{0x80, 0x80, 0x80, 0x6b, 0x56, 0xc0},
				Cr:             []uint8{0x80, 0x80, 0x80, 0xc0, 0x4b, 0x76},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 6, 1),
				Stride: 6 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff,
					0xff, 0xff, 0xff, 0xff,
					0x7f, 0x7f, 0x7f, 0xff,
					0x7f, 0x00, 0x00, 0xff,
					0x00, 0x7f, 0x00, 0xff,
					0x00, 0x00, 0x7f, 0xff,
				},
			},
		},
		{
			"Clone Paletted",
			&image.Paletted{
				Rect:   image.Rect(-1, -1, 5, 0),
				Stride: 6 * 1,
				Palette: color.Palette{
					color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
					color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
					color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0xff},
					color.NRGBA{R: 0x7f, G: 0x00, B: 0x00, A: 0xff},
					color.NRGBA{R: 0x00, G: 0x7f, B: 0x00, A: 0xff},
					color.NRGBA{R: 0x00, G: 0x00, B: 0x7f, A: 0xff},
				},
				Pix: []uint8{0x0, 0x1, 0x2, 0x3, 0x4, 0x5},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 6, 1),
				Stride: 6 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff,
					0xff, 0xff, 0xff, 0xff,
					0x7f, 0x7f, 0x7f, 0xff,
					0x7f, 0x00, 0x00, 0xff,
					0x00, 0x7f, 0x00, 0xff,
					0x00, 0x00, 0x7f, 0xff,
				},
			},
		},
	}

	for _, d := range td {
		got := Clone(d.src)
		want := d.want

		delta := 0
		if _, ok := d.src.(*image.YCbCr); ok {
			delta = 1
		}

		if !compareNRGBA(got, want, delta) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

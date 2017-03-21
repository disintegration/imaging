package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
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
	err := Encode(buf, imgWithAlpha, JPEG)
	if err != nil {
		t.Errorf("failed encoding alpha to JPEG format %s", err)
	}

	buf = &bytes.Buffer{}
	err = Encode(buf, imgWithAlpha, Format(100))
	if err != ErrUnsupportedFormat {
		t.Errorf("expected ErrUnsupportedFormat")
	}

	buf = bytes.NewBuffer([]byte("bad data"))
	_, err = Decode(buf)
	if err == nil {
		t.Errorf("decoding bad data, expected error")
	}
}

func TestEncodeWithOptionsDecode(t *testing.T) {
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

	options := map[Format]interface{}{
		JPEG: &jpeg.Options{Quality: 50},
		GIF:  &gif.Options{NumColors: 200},
	}

	for _, format := range []Format{JPEG, GIF} {
		img := imgWithoutAlpha

		buf := &bytes.Buffer{}
		err := EncodeWithOptions(buf, img, format, options[format])
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

		if compareNRGBA(img, img2cloned, delta) {
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
		{
			"New 0x0 white",
			0, 0,
			color.NRGBA{255, 255, 255, 255},
			image.Rect(0, 0, 0, 0),
			nil,
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
				Rect:   image.Rect(-1, -1, 0, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 3),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone RGBA64",
			&image.RGBA64{
				Rect:   image.Rect(-1, -1, 0, 2),
				Stride: 1 * 8,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x11, 0x11, 0x22, 0x22, 0x33, 0x33,
					0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee, 0xff, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 3),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
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
			"Clone YCbCr 444",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0x55, 0xd4, 0xff, 0x8e, 0x2c, 0x01, 0x6b, 0xaa, 0xc0, 0x95, 0x56, 0x40, 0x80, 0x80, 0x80, 0x80},
				Cr:             []uint8{0xff, 0xeb, 0x6b, 0x36, 0x15, 0x95, 0xc0, 0xb5, 0x76, 0x41, 0x4b, 0x8c, 0x80, 0x80, 0x80, 0x80},
				YStride:        4,
				CStride:        4,
				SubsampleRatio: image.YCbCrSubsampleRatio444,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x49, 0xe1, 0xca, 0xff, 0x0, 0xff, 0x0, 0xff, 0xff, 0xff, 0x0, 0xff, 0x7f, 0x0, 0x0, 0xff, 0x7f, 0x0, 0x7f, 0xff, 0x0, 0x0, 0x7f, 0xff, 0x0, 0x7f, 0x7f, 0xff, 0x0, 0x7f, 0x0, 0xff, 0x82, 0x7f, 0x0, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 440",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0x2c, 0x01, 0x6b, 0xaa, 0x80, 0x80, 0x80, 0x80},
				Cr:             []uint8{0x15, 0x95, 0xc0, 0xb5, 0x80, 0x80, 0x80, 0x80},
				YStride:        4,
				CStride:        4,
				SubsampleRatio: image.YCbCrSubsampleRatio440,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0x0, 0xb5, 0x0, 0xff, 0x86, 0x86, 0x0, 0xff, 0x77, 0x0, 0x0, 0xff, 0xfb, 0x7d, 0xfb, 0xff, 0x0, 0xff, 0x1, 0xff, 0xff, 0xff, 0x1, 0xff, 0x80, 0x0, 0x1, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0xe, 0xe, 0xe, 0xff, 0x59, 0x59, 0x59, 0xff, 0x4b, 0x4b, 0x4b, 0xff, 0x71, 0x71, 0x71, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 422",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0xd4, 0x8e, 0x01, 0xaa, 0x95, 0x40, 0x80, 0x80},
				Cr:             []uint8{0xeb, 0x36, 0x95, 0xb5, 0x41, 0x8c, 0x80, 0x80},
				YStride:        4,
				CStride:        2,
				SubsampleRatio: image.YCbCrSubsampleRatio422,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0xe2, 0x0, 0xe1, 0xff, 0xff, 0x0, 0xfe, 0xff, 0x0, 0x4d, 0x36, 0xff, 0x49, 0xe1, 0xca, 0xff, 0xb3, 0xb3, 0x0, 0xff, 0xff, 0xff, 0x1, 0xff, 0x70, 0x0, 0x70, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0x0, 0x34, 0x33, 0xff, 0x1, 0x7f, 0x7e, 0xff, 0x5c, 0x58, 0x0, 0xff, 0x82, 0x7e, 0x0, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 420",
			&image.YCbCr{
				Y:       []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:      []uint8{0x01, 0xaa, 0x80, 0x80},
				Cr:      []uint8{0x95, 0xb5, 0x80, 0x80},
				YStride: 4, CStride: 2,
				SubsampleRatio: image.YCbCrSubsampleRatio420,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0x69, 0x69, 0x0, 0xff, 0x86, 0x86, 0x0, 0xff, 0x67, 0x0, 0x67, 0xff, 0xfb, 0x7d, 0xfb, 0xff, 0xb3, 0xb3, 0x0, 0xff, 0xff, 0xff, 0x1, 0xff, 0x70, 0x0, 0x70, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0xe, 0xe, 0xe, 0xff, 0x59, 0x59, 0x59, 0xff, 0x4b, 0x4b, 0x4b, 0xff, 0x71, 0x71, 0x71, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
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

func TestFormats(t *testing.T) {
	formatNames := map[Format]string{
		JPEG:       "JPEG",
		PNG:        "PNG",
		GIF:        "GIF",
		BMP:        "BMP",
		TIFF:       "TIFF",
		Format(-1): "Unsupported",
	}
	for format, name := range formatNames {
		got := format.String()
		if got != name {
			t.Errorf("test [Format names] failed: got %#v want %#v", got, name)
			continue
		}
	}
}

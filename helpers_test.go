package imaging

import (
	"bytes"
	"image"
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
			t.Errorf("fail comparing images format=%s %#v %#v", format, img, img2cloned)
			continue
		}
	}

	buf := &bytes.Buffer{}
	err := Encode(buf, imgWithAlpha, Format(100))
	if err != ErrUnsupportedFormat {
		t.Errorf("expected ErrUnsupportedFormat")
	}
}

package imaging

import (
	"image"
	"testing"
)

func TestHistogram(t *testing.T) {
	testCases := []struct {
		name string
		img  image.Image
		want [256]float64
	}{
		{
			name: "grayscale",
			img: &image.RGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0x80, 0x80, 0x80, 0xff,
				},
			},
			want: [256]float64{0x00: 0.25, 0x80: 0.25, 0xff: 0.5},
		},
		{
			name: "colorful",
			img: &image.RGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff, 0x33, 0x44, 0x55, 0xff,
					0x55, 0x44, 0x33, 0xff, 0x77, 0x66, 0x55, 0xff,
				},
			},
			want: [256]float64{0x00: 0.25, 0x41: 0.25, 0x47: 0.25, 0x69: 0.25},
		},
		{
			name: "zero",
			img:  &image.RGBA{},
			want: [256]float64{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Histogram(tc.img)
			if got != tc.want {
				t.Fatalf("got histogram %#v want %#v", got, tc.want)
			}
		})
	}
}

func BenchmarkHistogram(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Histogram(testdataBranchesJPG)
	}
}

package imaging

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"testing"
)

var (
	testdataBranchesJPG     = mustOpen("testdata/branches.jpg")
	testdataBranchesPNG     = mustOpen("testdata/branches.png")
	testdataFlowersSmallPNG = mustOpen("testdata/flowers_small.png")
)

func mustOpen(filename string) image.Image {
	img, err := Open(filename)
	if err != nil {
		panic(err)
	}
	return img
}

func TestParallel(t *testing.T) {
	for _, n := range []int{0, 1, 10, 100, 1000} {
		for _, p := range []int{1, 2, 4, 8, 16, 100} {
			if !testParallelN(n, p) {
				t.Fatalf("test [parallel %d %d] failed", n, p)
			}
		}
	}
}

func testParallelN(n, procs int) bool {
	data := make([]bool, n)
	before := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(procs)
	parallel(0, n, func(is <-chan int) {
		for i := range is {
			data[i] = true
		}
	})
	runtime.GOMAXPROCS(before)
	for i := 0; i < n; i++ {
		if !data[i] {
			return false
		}
	}
	return true
}

func TestClamp(t *testing.T) {
	testCases := []struct {
		f float64
		u uint8
	}{
		{0, 0},
		{255, 255},
		{128, 128},
		{0.49, 0},
		{0.50, 1},
		{254.9, 255},
		{254.0, 254},
		{256, 255},
		{2500, 255},
		{-10, 0},
		{127.6, 128},
	}

	for _, tc := range testCases {
		if clamp(tc.f) != tc.u {
			t.Fatalf("test [clamp %v %v] failed: %v", tc.f, tc.u, clamp(tc.f))
		}
	}
}

func TestReverse(t *testing.T) {
	testCases := []struct {
		pix  []uint8
		want []uint8
	}{
		{
			pix:  []uint8{},
			want: []uint8{},
		},
		{
			pix:  []uint8{1, 2, 3, 4},
			want: []uint8{1, 2, 3, 4},
		},
		{
			pix:  []uint8{1, 2, 3, 4, 5, 6, 7, 8},
			want: []uint8{5, 6, 7, 8, 1, 2, 3, 4},
		},
		{
			pix:  []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			want: []uint8{9, 10, 11, 12, 5, 6, 7, 8, 1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			reverse(tc.pix)
			if !compareBytes(tc.pix, tc.want, 0) {
				t.Fatalf("got pix %v want %v", tc.pix, tc.want)
			}
		})
	}
}

func compareNRGBA(img1, img2 *image.NRGBA, delta int) bool {
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}
	return compareBytes(img1.Pix, img2.Pix, delta)
}

func compareBytes(a, b []uint8, delta int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if absint(int(a[i])-int(b[i])) > delta {
			return false
		}
	}
	return true
}

func compareFloat64(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

func compareUint8(a, b, delta uint8) bool {
	if a > b {
		return a - b <= delta
	} else {
		return b - a <= delta
	}
}

var nrgbaHslTestCases = []struct {
	rgb     color.NRGBA
	h, s, l float64
}{
	{
		rgb: color.NRGBA{R: 255, G: 0, B: 0},
		h:   0,
		s:   1,
		l:   0.5,
	},
	{
		rgb: color.NRGBA{R: 191, G: 191, B: 0},
		h:   60,
		s:   1,
		l:   0.375,
	},
	{
		rgb: color.NRGBA{R: 0, G: 128, B: 0},
		h:   120,
		s:   1,
		l:   0.25,
	},
	{
		rgb: color.NRGBA{R: 128, G: 255, B: 255},
		h:   180,
		s:   1,
		l:   0.75,
	},
	{
		rgb: color.NRGBA{R: 128, G: 128, B: 255},
		h:   240,
		s:   1,
		l:   0.75,
	},
	{
		rgb: color.NRGBA{R: 191, G: 64, B: 191},
		h:   300,
		s:   0.5,
		l:   0.5,
	},
	{
		rgb: color.NRGBA{R: 160, G: 164, B: 36},
		h:   61.8,
		s:   0.638,
		l:   0.393,
	},
	{
		rgb: color.NRGBA{R: 65, G: 27, B: 234},
		h:   251.1,
		s:   0.832,
		l:   0.511,
	},
	{
		rgb: color.NRGBA{R: 30, G: 172, B: 65},
		h:   134.9,
		s:   0.707,
		l:   0.396,
	},
	{
		rgb: color.NRGBA{R: 240, G: 200, B: 14},
		h:   49.5,
		s:   0.893,
		l:   0.497,
	},
	{
		rgb: color.NRGBA{R: 180, G: 48, B: 229},
		h:   283.7,
		s:   0.775,
		l:   0.542,
	},
	{
		rgb: color.NRGBA{R: 237, G: 118, B: 81},
		h:   14.3,
		s:   0.817,
		l:   0.624,
	},
	{
		rgb: color.NRGBA{R: 254, G: 248, B: 136},
		h:   56.9,
		s:   0.991,
		l:   0.765,
	},
	{
		rgb: color.NRGBA{R: 25, G: 203, B: 151},
		h:   162.4,
		s:   0.779,
		l:   0.447,
	},
	{
		rgb: color.NRGBA{R: 54, G: 38, B: 152},
		h:   248.3,
		s:   0.601,
		l:   0.373,
	},
	{
		rgb: color.NRGBA{R: 126, G: 126, B: 184},
		h:   240.5,
		s:   0.29,
		l:   0.607,
	},
}

func TestNrgbaToHSL(t *testing.T) {
	for _, tc := range nrgbaHslTestCases {
		t.Run("", func(t *testing.T) {
			h, s, l := nrgbaToHSL(tc.rgb)
			if !compareFloat64(h, tc.h, 1) || !compareFloat64(s, tc.s, 1) || !compareFloat64(l, tc.l, 1) {
				t.Fatalf("with %v, expected (%.2f, %.2f, %.2f) but got (%.2f, %.2f, %.2f)", tc.rgb, h, s, l, tc.h, tc.s, tc.l)
			}
		})
	}
}

func TestHslToNRGBA(t *testing.T) {
	for _, tc := range nrgbaHslTestCases {
		t.Run("", func(t *testing.T) {
			rgb := hslToNRGBA(tc.h, tc.s, tc.l)
			if !compareUint8(rgb.R, tc.rgb.R, 1) || !compareUint8(rgb.G, tc.rgb.G, 1) || !compareUint8(rgb.B, tc.rgb.B, 1) {
				t.Fatalf("expected %+v but got %+v", tc.rgb, rgb)
			}
		})
	}
}

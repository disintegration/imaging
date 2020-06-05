package imaging

import (
	"image"
	"math"
	"runtime"
	"sync/atomic"
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

func TestParallelMaxProcs(t *testing.T) {
	for _, n := range []int{0, 1, 10, 100, 1000} {
		for _, p := range []int{1, 2, 4, 8, 16, 100} {
			if !testParallelMaxProcsN(n, p) {
				t.Fatalf("test [parallel max procs %d %d] failed", n, p)
			}
		}
	}
}

func testParallelMaxProcsN(n, procs int) bool {
	data := make([]bool, n)
	SetMaxProcs(procs)
	parallel(0, n, func(is <-chan int) {
		for i := range is {
			data[i] = true
		}
	})
	SetMaxProcs(0)
	for i := 0; i < n; i++ {
		if !data[i] {
			return false
		}
	}
	return true
}

func TestSetMaxProcs(t *testing.T) {
	for _, p := range []int{-1, 0, 10} {
		SetMaxProcs(p)
		if int(atomic.LoadInt64(&maxProcs)) != p {
			t.Fatalf("test [set max procs %d] failed", p)
		}
	}

	SetMaxProcs(0)
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

// compareNRGBAGolden is a special version of compareNRGBA used in golden tests.
// All the golden images are generated on amd64 architecture. Due to differences
// in floating-point rounding on different architectures, we need to add some
// level of tolerance when comparing images on architectures other than amd64.
// See https://golang.org/ref/spec#Floating_point_operators for information on
// fused multiply and add (FMA) instruction.
func compareNRGBAGolden(img1, img2 *image.NRGBA) bool {
	delta := 0
	if runtime.GOARCH != "amd64" {
		delta = 1
	}
	return compareNRGBA(img1, img2, delta)
}

func compareFloat64(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

var rgbHSLTestCases = []struct {
	r, g, b uint8
	h, s, l float64
}{
	{
		r: 255,
		g: 0,
		b: 0,
		h: 0.000,
		s: 1.000,
		l: 0.500,
	},
	{
		r: 191,
		g: 191,
		b: 0,
		h: 0.167,
		s: 1.000,
		l: 0.375,
	},
	{
		r: 0,
		g: 128,
		b: 0,
		h: 0.333,
		s: 1.000,
		l: 0.251,
	},
	{
		r: 128,
		g: 255,
		b: 255,
		h: 0.500,
		s: 1.000,
		l: 0.751,
	},
	{
		r: 128,
		g: 128,
		b: 255,
		h: 0.667,
		s: 1.000,
		l: 0.751,
	},
	{
		r: 191,
		g: 64,
		b: 191,
		h: 0.833,
		s: 0.498,
		l: 0.500,
	},
	{
		r: 160,
		g: 164,
		b: 36,
		h: 0.172,
		s: 0.640,
		l: 0.392,
	},
	{
		r: 65,
		g: 27,
		b: 234,
		h: 0.697,
		s: 0.831,
		l: 0.512,
	},
	{
		r: 30,
		g: 172,
		b: 65,
		h: 0.374,
		s: 0.703,
		l: 0.396,
	},
	{
		r: 240,
		g: 200,
		b: 14,
		h: 0.137,
		s: 0.890,
		l: 0.498,
	},
	{
		r: 180,
		g: 48,
		b: 229,
		h: 0.788,
		s: 0.777,
		l: 0.543,
	},
	{
		r: 237,
		g: 119,
		b: 81,
		h: 0.040,
		s: 0.813,
		l: 0.624,
	},
	{
		r: 254,
		g: 248,
		b: 136,
		h: 0.158,
		s: 0.983,
		l: 0.765,
	},
	{
		r: 25,
		g: 203,
		b: 151,
		h: 0.451,
		s: 0.781,
		l: 0.447,
	},
	{
		r: 54,
		g: 38,
		b: 152,
		h: 0.690,
		s: 0.600,
		l: 0.373,
	},
	{
		r: 126,
		g: 126,
		b: 184,
		h: 0.667,
		s: 0.290,
		l: 0.608,
	},
}

func TestRGBToHSL(t *testing.T) {
	for _, tc := range rgbHSLTestCases {
		t.Run("", func(t *testing.T) {
			h, s, l := rgbToHSL(tc.r, tc.g, tc.b)
			if !compareFloat64(h, tc.h, 0.001) || !compareFloat64(s, tc.s, 0.001) || !compareFloat64(l, tc.l, 0.001) {
				t.Fatalf("(%d, %d, %d): got (%.3f, %.3f, %.3f) want (%.3f, %.3f, %.3f)", tc.r, tc.g, tc.b, h, s, l, tc.h, tc.s, tc.l)
			}
		})
	}
}

func TestHSLToRGB(t *testing.T) {
	for _, tc := range rgbHSLTestCases {
		t.Run("", func(t *testing.T) {
			r, g, b := hslToRGB(tc.h, tc.s, tc.l)
			if r != tc.r || g != tc.g || b != tc.b {
				t.Fatalf("(%.3f, %.3f, %.3f): got (%d, %d, %d) want (%d, %d, %d)", tc.h, tc.s, tc.l, r, g, b, tc.r, tc.g, tc.b)
			}
		})
	}
}

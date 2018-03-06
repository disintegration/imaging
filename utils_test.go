package imaging

import (
	"image"
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

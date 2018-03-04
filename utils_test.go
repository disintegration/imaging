package imaging

import (
	"image"
	"runtime"
	"testing"
)

var (
	testdataBranchesJPG     = mustOpen("testdata/branches.jpg")
	testdataBranchesPNG     = mustOpen("testdata/branches.png")
	testdataFlowersPNG      = mustOpen("testdata/flowers.png")
	testdataFlowersSmallPNG = mustOpen("testdata/flowers_small.png")
)

func mustOpen(filename string) image.Image {
	img, err := Open(filename)
	if err != nil {
		panic(err)
	}
	return img
}

func testParallelN(n, procs int) bool {
	data := make([]bool, n)
	before := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(procs)
	parallel(n, func(start, end int) {
		for i := start; i < end; i++ {
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

func TestParallel(t *testing.T) {
	for _, n := range []int{0, 1, 10, 100, 1000} {
		for _, p := range []int{1, 2, 4, 8, 16, 100} {
			if !testParallelN(n, p) {
				t.Errorf("test [parallel %d %d] failed", n, p)
			}
		}
	}
}

func TestClamp(t *testing.T) {
	td := []struct {
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

	for _, d := range td {
		if clamp(d.f) != d.u {
			t.Errorf("test [clamp %v %v] failed: %v", d.f, d.u, clamp(d.f))
		}
	}
}

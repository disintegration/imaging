package imaging

import (
	"image"
	"runtime"
	"sync"
	"math"
	"image/color"
)

// parallel processes the data in separate goroutines.
func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := runtime.GOMAXPROCS(0)
	if procs > count {
		procs = count
	}

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}

// absint returns the absolute value of i.
func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// clamp rounds and clamps float64 value to fit into uint8.
func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

func reverse(pix []uint8) {
	if len(pix) <= 4 {
		return
	}
	i := 0
	j := len(pix) - 4
	for i < j {
		pix[i+0], pix[j+0] = pix[j+0], pix[i+0]
		pix[i+1], pix[j+1] = pix[j+1], pix[i+1]
		pix[i+2], pix[j+2] = pix[j+2], pix[i+2]
		pix[i+3], pix[j+3] = pix[j+3], pix[i+3]
		i += 4
		j -= 4
	}
}

func toNRGBA(img image.Image) *image.NRGBA {
	if img, ok := img.(*image.NRGBA); ok {
		return &image.NRGBA{
			Pix:    img.Pix,
			Stride: img.Stride,
			Rect:   img.Rect.Sub(img.Rect.Min),
		}
	}
	return Clone(img)
}

// nrgbaToHSL converts NRGBA to HSL.
func nrgbaToHSL(c color.NRGBA) (float64, float64, float64) {
	var h, s, l float64

	r := float64(c.R) / float64(255)
	g := float64(c.G) / float64(255)
	b := float64(c.B) / float64(255)

	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)

	l = (max + min) / 2

	if min == max {
		s = 0
		h = 0
	} else {
		if l < 0.5 {
			s = (max - min) / (max + min)
		} else {
			s = (max - min) / (2.0 - max - min)
		}

		if max == r {
			h = (g - b) / (max - min)
		} else if max == g {
			h = 2.0 + (b-r)/(max-min)
		} else {
			h = 4.0 + (r-g)/(max-min)
		}

		h *= 60

		if h < 0 {
			h += 360
		}
	}

	return h, s, l
}

// hslToNRGBA converts HSL to NRGBA with A=1.
func hslToNRGBA(h, s, l float64) color.NRGBA {
	if s == 0 {
		c := uint8(l * 255)
		return color.NRGBA{R: c, G: c, B: c, A: 1}
	}

	var r, g, b float64
	var t1, t2, tr, tg, tb float64

	if l < 0.5 {
		t1 = l * (1.0 + s)
	} else {
		t1 = l + s - l*s
	}

	t2 = 2*l - t1
	h = h / 360
	tr = h + 1.0/3.0
	tg = h
	tb = h - 1.0/3.0

	if tr < 0 {
		tr++
	} else if tr > 1 {
		tr--
	}

	if tg < 0 {
		tg++
	} else if tg > 1 {
		tg--
	}

	if tb < 0 {
		tb++
	} else if tb > 1 {
		tb--
	}

	// Red
	if 6*tr < 1 {
		r = t2 + (t1-t2)*6*tr
	} else if 2*tr < 1 {
		r = t1
	} else if 3*tr < 2 {
		r = t2 + (t1-t2)*(2.0/3.0-tr)*6
	} else {
		r = t2
	}

	// Green
	if 6*tg < 1 {
		g = t2 + (t1-t2)*6*tg
	} else if 2*tg < 1 {
		g = t1
	} else if 3*tg < 2 {
		g = t2 + (t1-t2)*(2.0/3.0-tg)*6
	} else {
		g = t2
	}

	// Blue
	if 6*tb < 1 {
		b = t2 + (t1-t2)*6*tb
	} else if 2*tb < 1 {
		b = t1
	} else if 3*tb < 2 {
		b = t2 + (t1-t2)*(2.0/3.0-tb)*6
	} else {
		b = t2
	}

	return color.NRGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255)}
}
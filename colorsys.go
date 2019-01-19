package utils

//Conversion functions between any two of the color systems.

//All inputs and outputs are three floats in the range [0.0...1.0]
//(with the exception of I and Q, which covers a slightly larger range
// also with the exception of R, G and B, which range from 0 to 255
// also with the exception of H, which range from 0 to 360).

//Inputs outside the valid range may cause exceptions or invalid outputs.

//Supported color systems:
//RGB: Red, Green, Blue components
//YIQ: Luminance, Chrominance (used by composite video signals)
//HLS: Hue, Luminance, Saturation
//HSV: Hue, Saturation, Value

var ONE_THIRD float64
var ONE_SIXTH float64
var TWO_THIRD float64

func init(){
	ONE_THIRD = 1.0/3
	ONE_SIXTH = 1.0/6
	TWO_THIRD = 2.0/3
}

// Util Functions

func max3(a,b,c float64)float64{
	max := a
	if b > max{
		max = b
	}
	if c > max{
		max = c
	}
	return max
}

func min3(a,b,c float64)float64{
	min := a
	if b < min{
		min = b
	}
	if c < min{
		min = c
	}
	return min
}

func _v(m1, m2, hue float64)float64{
	hue = hue - float64(int(hue))
	if hue < ONE_SIXTH{
		return m1 + (m2-m1)*hue*6.0
	}
	if hue < 0.5{
		return m2
	}
	if hue < TWO_THIRD{
		return m1 + (m2-m1)*(TWO_THIRD-hue)*6.0
	}
	return m1
}

// YIQ: used by composite video signals (linear combinations of RGB)
// Y: perceived grey level (0.0 == black, 1.0 == white)
// I, Q: color components

func Rgb2Yiq(r,g,b  uint32)(float64,float64,float64){
	r1 := float64(r) / 255
	g1 := float64(g) / 255
	b1 := float64(b) / 255
	y := 0.30*float64(r1) + 0.59*float64(g1) + 0.11*float64(b1)
	i := 0.60*float64(r1) - 0.28*float64(g1) - 0.32*float64(b1)
	q := 0.21*float64(r1) - 0.52*float64(g1) + 0.31*float64(b1)
	return y, i, q
}

func Yiq2Rgb(y,i,q float64)(uint32, uint32, uint32){
	r := y + 0.948262*i + 0.624013*q
	g := y - 0.276066*i - 0.639810*q
	b := y - 1.105450*i + 1.729860*q
	if r < 0.0{
		r = 0.0
	}
	if g < 0.0{
		g = 0.0
	}
	if b < 0.0{
		b = 0.0
	}
	if r > 1.0{
		r = 1.0
	}
	if g > 1.0{
		g = 1.0
	}
	if b > 1.0{
		b = 1.0
	}
	return uint32(r*255), uint32(g*255), uint32(b*255)
}

// HLS: Hue, Luminance, Saturation
// H: position in the spectrum(0~360)
// L: color lightness
// S: color saturation

func Rgb2Hls(r, g, b uint32)(float64, float64, float64){
	var h, l, s float64
	r1 := float64(r) / 255
	g1 := float64(g) / 255
	b1 := float64(b) / 255
	maxc := max3(r1, g1, b1)
	minc := min3(r1, g1, b1)
	l = float64(maxc + minc)/2.0
	if minc == maxc{
		return 0.0, l, 0.0
	}
	if l <= 0.5{
		s = float64(maxc-minc) / float64(maxc+minc)
	}else{
		s = float64(maxc-minc) / float64(2.0-maxc-minc)
	}
	rc := (maxc-r1) / (maxc-minc)
	gc := (maxc-g1) / (maxc-minc)
	bc := (maxc-b1) / (maxc-minc)
	if r1 == maxc{
		h = bc - gc
	}else if g1 == maxc{
		h = 2.0 + rc - bc
	}else{
		h = 4.0 + gc - rc
	}
	h0 := (h/6.0)
	h = h0 - float64(int(h0))
	return h*360, l, s
}

func Hls2Rgb(h,l,s float64)(uint32,uint32,uint32){
	var m1, m2 float64
	if s == 0.0{
		return uint32(l), uint32(l), uint32(l)
	}
	if l <= 0.5{
		m2 = l * (1.0+s)
	}else{
		m2 = l+s-(l*s)
	}
	m1 = 2.0*l - m2
	h /= 360
	return uint32(_v(m1,m2,h+ONE_THIRD)*255), uint32(_v(m1, m2, h)*255), uint32(_v(m1, m2, h-ONE_THIRD)*255)
}

// HSV: Hue, Saturation, Value
// H: position in the spectrum(0~360)
// S: color saturation ("purity")
// V: color brightness

func Rgb2Hsv(r, g ,b uint32)(float64, float64, float64){
	var h, s, v float64
	r1 := float64(r) / 255
	g1 := float64(g) / 255
	b1 := float64(b) / 255
	maxc := max3(r1, g1, b1)
	minc := min3(r1, g1, b1)
	v = maxc
	if minc == maxc{
		return 0.0, 0.0, float64(v)
	}
	s = (maxc-minc) / maxc
	rc := (maxc-r1) / (maxc-minc)
	gc := (maxc-g1) / (maxc-minc)
	bc := (maxc-b1) / (maxc-minc)
	if r1 == maxc{
		h = bc-gc
	}else if(g1 == maxc){
		h = 2.0 + rc - bc
	}else{
		h = 4.0 + gc - rc
	}
	h0 := (h/6.0)
	h = h0 - float64(int(h0))
	return h*360, s, v
}

func Hsv2Rgb(h, s, v float64)(uint32, uint32, uint32){
	h /= 360
	if s == 0.0{
		return uint32(v*255), uint32(v*255), uint32(v*255)
	}
	i := int(h*6.0)
	f := (h*6.0) - float64(i)
	p := v*(1.0 - s)
	q := v*(1.0 - s*f)
	t := v*(1.0 - s*(1.0-f))
	i %= 6
	switch i{
	case 0:
		return uint32(v*255), uint32(t*255), uint32(p*255)
	case 1:
		return uint32(q*255), uint32(v*255), uint32(p*255)
	case 2:
		return uint32(p*255), uint32(v*255), uint32(t*255)
	case 3:
		return uint32(p*255), uint32(q*255), uint32(v*255)
	case 4:
		return uint32(t*255), uint32(p*255), uint32(v*255)
	case 5:
		return uint32(v*255), uint32(p*255), uint32(q*255)
	default:
		return 0, 0, 0
	}
}

func Yiq2Hls(y, i, q float64)(float64, float64, float64){
	r, g, b := Yiq2Rgb(y, i, q)
	return Rgb2Hls(r, g, b)
}

func Yiq2Hsv(y, i, q float64)(float64, float64, float64){
	r, g, b := Yiq2Rgb(y, i, q)
	return Rgb2Hsv(r, g, b)
}

func Hls2Yiq(h, l, s float64)(float64, float64, float64){
	r, g, b := Hls2Rgb(h, l, s)
	return Rgb2Yiq(r, g, b)
}

func Hls2Hsv(h, l, s float64)(float64, float64, float64){
	r, g, b := Hls2Rgb(h, l, s)
	return Rgb2Hsv(r, g, b)
}

func Hsv2Yiq(h, s, v float64)(float64, float64, float64){
	r, g, b := Hsv2Rgb(h, s, v)
	return Rgb2Yiq(r, g, b)
}

func Hsv2Hls(h, s, v float64)(float64, float64, float64){
	r, g, b := Hsv2Rgb(h, s, v)
	return Rgb2Hls(r, g, b)
}
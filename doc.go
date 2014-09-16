/*
Package imaging provides basic image manipulation functions (resize, rotate, flip, crop, etc.).
This package is based on the standard Go image package and works best along with it.

Image manipulation functions provided by the package take any image type
that implements `image.Image` interface as an input, and return a new image of
`*image.NRGBA` type (32bit RGBA colors, not premultiplied by alpha).

Imaging package uses parallel goroutines for faster image processing.
To achieve maximum performance, make sure to allow Go to utilize all CPU cores:

	runtime.GOMAXPROCS(runtime.NumCPU())

Here is the complete example that loades several images, makes thumbnails of them
and joins them together.

	package main

	import (
		"image"
		"image/color"
		"runtime"

		"github.com/disintegration/imaging"
	)

	func main() {
		// use all CPU cores for maximum performance
		runtime.GOMAXPROCS(runtime.NumCPU())

		// input files
		files := []string{"01.jpg", "02.jpg", "03.jpg"}

		// load images and make 100x100 thumbnails of them
		var thumbnails []image.Image
		for _, file := range files {
			img, err := imaging.Open(file)
			if err != nil {
				panic(err)
			}
			thumb := imaging.Thumbnail(img, 100, 100, imaging.CatmullRom)
			thumbnails = append(thumbnails, thumb)
		}

		// create a new blank image
		dst := imaging.New(100*len(thumbnails), 100, color.NRGBA{0, 0, 0, 0})

		// paste thumbnails into the new image side by side
		for i, thumb := range thumbnails {
			dst = imaging.Paste(dst, thumb, image.Pt(i*100, 0))
		}

		// save the combined image to file
		err := imaging.Save(dst, "dst.jpg")
		if err != nil {
			panic(err)
		}
	}
*/
package imaging

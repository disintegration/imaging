# Imaging

Package imaging provides basic image manipulation functions (resize, rotate, flip, crop, etc.). 
This package is based on the standard Go image package and works best along with it. 

## Installation

Imaging requires Go version 1.2 or greater.

    go get -u github.com/disintegration/imaging

*[Git](http://git-scm.com/) and [Mercurial](http://mercurial.selenic.com/) are needed.*
    
## Documentation

http://godoc.org/github.com/disintegration/imaging

## Overview

Image manipulation functions provided by the package take any image type 
that implements `image.Image` interface as an input, and return a new image of 
`*image.NRGBA` type (32bit RGBA colors, not premultiplied by alpha).

Some of examples below require importing standard `runtime`, `image` or `image/color` packages.

```go 
// allow Go to utilize all CPU cores:
runtime.GOMAXPROCS(runtime.NumCPU())

// resize srcImage to size = 800x600px using the high quality Lanczos filter
dstImage := imaging.Resize(srcImage, 800, 600, imaging.Lanczos)

// resize srcImage to width = 800px preserving the aspect ratio
dstImage := imaging.Resize(srcImage, 800, 0, imaging.Lanczos)

// scale down srcImage to fit the 800x600px bounding box
dstImage := imaging.Fit(srcImage, 800, 600, imaging.Lanczos)

// resize and crop the srcImage to make a 100x100px thumbnail
dstImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

// cut out a rectangular region from the image
dstImage := imaging.Crop(srcImage, image.Rect(50, 50, 100, 100)) 

// cut out a 100x100 px region from the center of the image
dstImage := imaging.CropCenter(srcImage, 100, 100)   

// paste the srcImage to the backgroundImage at the (50, 50) position
dstImage := imaging.Paste(backgroundImage, srcImage, image.Pt(50, 50))     

// paste the srcImage to the center of the backgroundImage
dstImage := imaging.PasteCenter(backgroundImage, srcImage)                   

// draw the srcImage over the backgroundImage at the (50, 50) position with opacity=0.5
dstImage := imaging.Overlay(backgroundImage, srcImage, image.Pt(50, 50), 0.5)

// blur the srcImage
dstImage := imaging.Blur(srcImage, 4.5)

// sharpen the srcImage
dstImage := imaging.Sharpen(srcImage, 3.0)

// gamma correction
dstImage := imaging.AdjustGamma(srcImage, 0.7)

// increase the brightness by 10%
dstImage := imaging.AdjustBrightness(srcImage, 10) 

// decrease the brightness by 15%
dstImage := imaging.AdjustBrightness(srcImage, -15) 

// increase the contrast by 20%
dstImage := imaging.AdjustContrast(srcImage, 20)

// decrease the contrast by 10%
dstImage := imaging.AdjustContrast(srcImage, -10) 

// increase the contrast using sigmoidal function
dstImage := imaging.AdjustSigmoid(srcImage, 0.5, 3.0) 

// decrease the contrast using sigmoidal function
dstImage := imaging.AdjustSigmoid(srcImage, 0.5, -3.0)

// produce the grayscaled version of an image
dstImage := imaging.Grayscale(srcImage)

// produce the inverted version of an image
dstImage := imaging.Invert(srcImage)

// read an image from io.Reader
img, err := imaging.Decode(r) 
if err != nil {
    panic(err)
}

// load an image from file
img, err := imaging.Open("src.png") 
if err != nil {
    panic(err)
}

// write the image to io.Writer
err := imaging.Encode(w, img, imaging.PNG) 
if err != nil {
    panic(err)
}

// save the image to file
err := imaging.Save(img, "dst.jpg") 
if err != nil {
    panic(err)
}

// create a new 800x600px image filled with red color
newImg := imaging.New(800, 600, color.NRGBA{255, 0, 0, 255})

// make a copy of the image
copiedImg := imaging.Clone(img)
```


## Code example
Here is the complete example that loads several images, makes thumbnails of them
and joins them together.

```go
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
```
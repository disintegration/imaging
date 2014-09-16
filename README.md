# Imaging

Package imaging provides basic image manipulation functions (resize, rotate, flip, crop, etc.). 
This package is based on the standard Go image package and works best along with it. 

## Installation

    go get -u github.com/disintegration/imaging

*[Git](http://git-scm.com/) and [Mercurial](http://mercurial.selenic.com/) are needed.*
    
## Documentation

http://godoc.org/github.com/disintegration/imaging

## Overview

Image manipulation functions provided by the package take any image type 
that implements `image.Image` interface as an input, and return a new image of 
`*image.NRGBA` type (32bit RGBA colors, not premultiplied by alpha).

Note: some of examples below require importing standard `image` or `image/color` packages.

### Parallelization

Imaging package uses parallel goroutines for faster image processing.
To achieve maximum performance, make sure to allow Go to utilize all CPU cores. Use standard `runtime` package:
```go 
runtime.GOMAXPROCS(runtime.NumCPU())
```

### Resize

There are three image resizing functions in the package: `Resize`, `Fit` and `Thumbnail`.

**Resize** resizes the image to the specified width and height using the specified resample
filter and returns the transformed image. If one of width or height is 0, the image aspect
ratio is preserved.

**Fit** scales down the image using the specified resample filter to fit the specified
maximum width and height and returns the transformed image. The image aspect
ratio is preserved.

**Thumbnail** scales the image up or down using the specified resample filter, crops it
to the specified width and hight and returns the transformed image.

All three resizing function take `ResampleFilter` as the last argument.
A complete list of supported filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
CatmullRom (cubic filter) and Lanczos are recommended for high quality general purpose image resizing. 
NearestNeighbor is the fastest one but applies no anti-aliasing.

```go
// resize srcImage to width = 800px preserving the aspect ratio
dstImage := imaging.Resize(srcImage, 800, 0, imaging.Lanczos)

// scale down srcImage to fit the 800x600px bounding box
dstImage = imaging.Fit(srcImage, 800, 600, imaging.Lanczos)

// make a 100x100px thumbnail from srcImage
dstImage = imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)
```

### Rotate & flip
    
Imaging package implements functions to rotate an image 90, 180 or 270 degrees (counter-clockwise)
and to flip an image horizontally or vertically.

```go
dstImage = imaging.Rotate90(srcImage)  // rotate 90 degrees counter-clockwise
dstImage = imaging.Rotate180(srcImage) // rotate 180 degrees counter-clockwise
dstImage = imaging.Rotate270(srcImage) // rotate 270 degrees counter-clockwise
dstImage = imaging.FlipH(srcImage)     // flip horizontally (from left to right)
dstImage = imaging.FlipV(srcImage)     // flip vertically (from top to bottom)
```

### Crop, Paste, Overlay

**Crop** cuts out a rectangular region with the specified bounds from the image and returns 
the cropped image. 

**CropCenter** cuts out a rectangular region with the specified size from the center of the image
and returns the cropped image.

**Paste** pastes one image into another at the specified position and returns the combined image.

**PasteCenter** pastes one image to the center of another image and returns the combined image.

**Overlay** draws one image over another image at the specified position with the specified opacity 
and returns the combined image. Opacity parameter must be from 0.0 (fully transparent) to 1.0 (opaque).

```go
// cut out a rectangular region from the image
dstImage = imaging.Crop(srcImage, image.Rect(50, 50, 100, 100)) 

// cut out a 100x100 px region from the center of the image
dstImage = imaging.CropCenter(srcImage, 100, 100)   

// paste the srcImage to the backgroundImage at the (50, 50) position
dstImage = imaging.Paste(backgroundImage, srcImage, image.Pt(50, 50))     

// paste the srcImage to the center of the backgroundImage
dstImage = imaging.PasteCenter(backgroundImage, srcImage)                   

// draw the srcImage over the backgroundImage at the (50, 50) position with opacity=0.5
dstImage = imaging.Overlay(backgroundImage, srcImage, image.Pt(50, 50), 0.5)
```
### Blur & Sharpen

**Blur** produces a blurred version of the image.

**Sharpen** produces a sharpened version of the image.

Both functions take the `sigma` argument that is used in a Gaussian function. 
Sigma must be a positive floating point value indicating how much the image 
will be blurred or sharpened and how many neighbours of each pixel will be affected.

```go
dstImage = imaging.Blur(srcImage, 4.5)
dstImage = imaging.Sharpen(srcImage, 3.0)
```

### Adjustments

**AdjustGamma** performs a gamma correction on the image and returns the adjusted image. 
Gamma parameter must be positive. Gamma = 1.0 gives the original image.
Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.

```go
dstImage = imaging.AdjustGamma(srcImage, 0.7)
```

**AdjustBrightness** changes the brightness of the image using the percentage parameter and returns the adjusted image. The percentage must be in range (-100, 100).
The percentage = 0 gives the original image. The percentage = -100 gives solid black image. The percentage = 100 gives solid white image. 

```go
dstImage = imaging.AdjustBrightness(srcImage, 10) // increase the brightness by 10%
dstImage = imaging.AdjustBrightness(srcImage, -15) // decrease the brightness by 15%
```

**AdjustContrast** changes the contrast of the image using the percentage parameter and returns the adjusted image. The percentage must be in range (-100, 100).
The percentage = 0 gives the original image. The percentage = -100 gives solid grey image.

```go
dstImage = imaging.AdjustContrast(srcImage, 20) // increase the contrast by 20%
dstImage = imaging.AdjustContrast(srcImage, -10) // decrease the contrast by 10%
```

**AdjustSigmoid** changes the contrast of the image using a sigmoidal function and returns the adjusted image.
It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.
The midpoint parameter is the midpoint of contrast that must be between 0 and 1, typically 0.5.
The factor parameter indicates how much to increase or decrease the contrast, typically in range (-10, 10).
If the factor parameter is positive the image contrast is increased otherwise the contrast is decreased.

```go
dstImage = imaging.AdjustSigmoid(srcImage, 0.5, 3.0) // increase the contrast
dstImage = imaging.AdjustSigmoid(srcImage, 0.5, -3.0) // decrease the contrast
```

**Grayscale** produces grayscale version of the image.

```go
dstImage = imaging.Grayscale(srcImage)
```

**Invert** produces inverted (negated) version of the image.

```go
dstImage = imaging.Invert(srcImage)
```

#### Open, Save, New, Clone

Imaging package provides useful shortcuts for image loading, saving, creation and copying.
Open and Save functions support JPEG, PNG, TIFF and BMP images. 
External libraries can be used to load other image formats.

```go
// load an image from file
img, err := imaging.Open("src.png") 
if err != nil {
    panic(err)
}

// save the image to file
err = imaging.Save(img, "dst.jpg") 
if err != nil {
    panic(err)
}

// create a new 800x600px image filled with red color
newImg := imaging.New(800, 600, color.NRGBA{255, 0, 0, 255})

// make a copy of the image
copiedImg := imaging.Clone(img)
```


## Code example
Here is the complete example that loades several images, makes thumbnails of them
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
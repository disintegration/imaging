# Imaging

Package imaging provides basic image manipulation functions 
(resize, rotate, flip, crop, etc.) as well as simplified image loading and saving.
This package is based on the standard Go image package. All the image 
manipulation functions provided by the package take any image type that 
implements `image.Image` interface, and return a new image of 
`*image.NRGBA` type (32 bit RGBA colors, not premultiplied by alpha). 

###Recent changes

- Resize, Fit and Thumbnail now take 4th argument - resample filter. 
Supported filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
- New function: Overlay. This function can be used to draw one (partially 
transparent) image over another, to blend two images, etc.


###Installation

    go get github.com/disintegration/imaging
    
### Documentation

http://godoc.org/github.com/disintegration/imaging
    
### Usage

```go
package main

import (
    "github.com/disintegration/imaging"
    "image"
    "image/color"
)

func main() {
    src, _ := imaging.Open("src.png") // load an image from file (returns image.Image interface)
    var dst *image.NRGBA
    
    dst = imaging.New(800, 600, color.NRGBA{255, 0, 0, 255}) // create a new 800x600px image filled with red color
    dst = imaging.Clone(src) // make a copy of the image
    
    dst = imaging.Rotate90(src) // rotate 90 degrees clockwise 
    dst = imaging.Rotate180(src) // rotate 180 degrees clockwise
    dst = imaging.Rotate270(src) // rotate 270 degrees clockwise

    dst = imaging.FlipH(src) // flip horizontally (from left to right)
    dst = imaging.FlipV(src) // flip vertically (from top to bottom)

    // Resize, Fit and Thumbnail functions take resampling filter as 4th argument.
    // Supported filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
    // CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.

    dst = imaging.Resize(src, 600, 400, imaging.CatmullRom) // resize to 600x400 px using CatmullRom cubic filter
    dst = imaging.Resize(src, 600, 0, imaging.CatmullRom) // resize to width = 600, preserve the image aspect ratio
    dst = imaging.Resize(src, 0, 400, imaging.CatmullRom) // resize to height = 400, preserve the image aspect ratio
    
    dst = imaging.Fit(src, 800, 600, imaging.CatmullRom) // scale down the image to fit the given maximum width and height
    dst = imaging.Thumbnail(src, 100, 100, imaging.CatmullRom) // resize and crop the image to make a 100x100 thumbnail
    
    dst = imaging.Crop(src, image.Rect(50, 50, 100, 100)) // cut out a rectangular region from the image
    dst = imaging.CropCenter(src, 200, 100) // cut out a 200x100 px region from the center of the image
    dst = imaging.Paste(dst, src, image.Pt(50, 50)) // paste the src image to the dst image at the given position
    dst = imaging.PasteCenter(dst, src) // paste the src image to the center of the dst image
    
    // draw one image over another at the given position and with the given opacity (from 0.0 to 1.0)
    dst = imaging.Overlay(dst, src, image.Pt(50, 30), 1.0)
    
    imaging.Save(dst, "dst.jpg") // save the image to file
}
```

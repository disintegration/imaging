# Imaging

Package imaging provides basic image manipulation functions 
(resize, rotate, crop, etc.) as well as simplified image loading and saving.
This package is based on the standard Go image package. All the image 
manipulation functions provided by the package return a new image of 
`*image.NRGBA` type (32 bit RGBA colors, not premultiplied by alpha). 

###Recent changes

- Format parameter removed from `Save` function. Now the format is determined
from the filename extension, `jpg` (or `jpeg`) and `png` are supported.
- All the image manipulation functions now return `*image.NRGBA` instead of
`draw.Image`.
- `Copy()` function renamed to `Clone()`. 
This function also can be used to convert any image type to *image.NRGBA for
fast pixel access. (`.Pix` slice, `.PixOffset` method)


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

    dst = imaging.Resize(src, 600, 400) // resize to 600x400 px    
    dst = imaging.Resize(src, 600, 0) // resize to width = 600, preserve the image aspect ratio
    dst = imaging.Resize(src, 0, 400) // resize to height = 400, preserve the image aspect ratio

    dst = imaging.Fit(src, 800, 600) // scale down the image to fit the given maximum width and height
    dst = imaging.Thumbnail(src, 100, 100) // resize and crop the image to make a 100x100 thumbnail
    
    dst = imaging.Crop(src, image.Rect(50, 50, 100, 100)) // cut out a rectangular region from the image
    dst = imaging.CropCenter(src, 200, 100) // cut out a 200x100 px region from the center of the image
    dst = imaging.Paste(dst, src, image.Pt(50, 50)) // paste the src image to the dst image at the given position
    dst = imaging.PasteCenter(dst, src) // paste the src image to the center of the dst image

    imaging.Save(dst, "dst.jpg") // save the image to file using jpeg format
}
```

# Imaging

Simple Go image processing package. 

Implements a number of basic image manipulation functions as well as simplified image loading, creation and saving. 
All the image manipulation functions are not in-place, they return a new image with bounds (0, 0) - (Width, Height)
This package integrates with standard "image" package well: most functions take image.Image interface as an argument 
and return draw.Image interface.

### Installation

    go get github.com/disintegration/imaging
    
### Usage

```go
package main

import (
    "github.com/disintegration/imaging"
    "image"
    "image/color"
    "image/draw"   
)

func main() {
    src, _ := imaging.Open("1.png") // load image from file (returns image.Image interface)
    var dst draw.Image
    
    dst = imaging.New(800, 600, color.NRGBA(255, 0, 0, 255)) // create a new 800x600px image filled with red color
    dst = imaging.Copy(src) // copy src image
    
    dst = imaging.Rotate90(src) // rotate by 90 degrees clockwise 
    dst = imaging.Rotate180(src) // rotate by 180 degrees clockwise
    dst = imaging.Rotate270(src) // rotate by 270 degrees clockwise

    dst = imaging.FlipH(src) // flip horizontally (left-to-right)
    dst = imaging.FlipV(src) // flip vertically (top-to-bottom)

    dst = imaging.Resize(src, 600, 400) // resize to 600x400 px    
    dst = imaging.Resize(src, 600, 0) // resize to width = 600, preserve the image aspect ratio
    dst = imaging.Resize(src, 0, 400) // resize to height = 400, preserve the image aspect ratio

    dst = imaging.Fit(src, 800, 600) // scale down to fit given maximum width and height, keep the aspect ratio
    dst = imaging.Thumbnail(src, 100, 100) // resize and crop the image to make a 100x100 thumbnail
    
    dst = imaging.Crop(src, image.Rect(50, 50, 100, 100)) // copy a rectangular region from image
    dst = imaging.CropCenter(src, 200, 100) // copy a rectangular region of given size from the center of image
    dst = imaging.Paste(dst, src, image.Pt(50, 50)) // paste src image to dst image at given position
    dst = imaging.PasteCenter(dst, src) // paste src image to the center of dst image
 
    imaging.Save(dst, "2.jpg", "jpeg") // save image to file using jpeg format
}
```

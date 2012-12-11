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
    src, _ := imaging.Open("1.png") // load an image from file (returns image.Image interface)
    var dst draw.Image
    
    dst = imaging.New(800, 600, color.NRGBA(255, 0, 0, 255)) // create a new 800x600px image filled with red color
    dst = imaging.Copy(src) // copy the image
    
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

    imaging.Save(dst, "2.jpg", "jpeg") // save the image to file using jpeg format
}
```

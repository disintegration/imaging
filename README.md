# Imaging

Simple Go image processing package. 

Implements a number of basic image transformations (resizing, rotation, flipping) 
as well as simplified image loading and saving.

Planned: thumbnails, cropping, fitting in given size, color filling, pasting one image into another, optimizations.

### Installation

    go get github.com/disintegration/imaging
    
### Usage

```go
package main

import (
  "github.com/disintegration/imaging"
  "image"
)

func main() {
	src, _ := imaging.Open("1.png") // load image from file (returns image.Image interface)
	var dst image.Image

	dst = imaging.Copy(src) // copy entire image
	dst = imaging.Rotate90(src) // rotate by 90 degrees clockwise 
	dst = imaging.Rotate180(src) // rotate by 180 degrees clockwise
	dst = imaging.Rotate270(src) // rotate by 270 degrees clockwise
	dst = imaging.FlipH(src) // flip horizontally (left-to-right)
	dst = imaging.FlipV(src) // flip vertically (top-to-bottom)
	dst = imaging.Resize(src, 600, 400) // resize to 600x400 px
  
	imaging.Save(dst, "2.jpg", "jpeg") // save image to file using jpeg format
}

```
    
package main

import (
	"image/jpeg"
	"image/png"
	"os"
)

func main() {
	r, _ := os.Open("images/equi1.jpeg")
	img, _ := jpeg.Decode(r)
	rendered := RenderOctant(NewEquirect(img), NewOctantMapper(true, 0))
	w, _ := os.Create("images/octant.png")
	defer w.Close()
	png.Encode(w, rendered)
}

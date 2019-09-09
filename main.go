package main

import (
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

func main() {
	r, _ := os.Open("images/equi1.jpeg")
	img, _ := jpeg.Decode(r)
	rendered := RenderStrip(NewEquirect(img), NewStripMapper(true, 0, math.Pi/8))
	w, _ := os.Create("images/strip.png")
	defer w.Close()
	png.Encode(w, rendered)
}

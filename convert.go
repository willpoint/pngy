package main

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

const dstExt = ".png"

// Decoder is the interface that wraps any type of file extension
// supported for conversion
type Decoder interface {
	Decode() (image.Image, error)
}

// JpegImage is a supported image type
type JpegImage struct {
	img *os.File
}

// GifImage is a supported image type
type GifImage struct {
	img *os.File
}

// Decode method for a GIF image to implement the Decoder interface
func (s *GifImage) Decode() (image.Image, error) {
	return gif.Decode(s.img)
}

// Decode method to for a jpeg|jpg to implement the Decoder interface
func (s *JpegImage) Decode() (image.Image, error) {
	return jpeg.Decode(s.img)
}

// Convert converts the given image to a png format
// dst is the new png image
func Convert(img image.Image, dst *os.File) error {
	return png.Encode(dst, img)
}
package convert

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

// Decoder is the interface that wraps any type of file extension
// supported for conversion
type Decoder interface {
	Decode() (image.Image, error)
}

// JpegImage is a supported image type
type JpegImage struct {
	Img *os.File
}

// GifImage is a supported image type
type GifImage struct {
	Img *os.File
}

// Decode method for a GIF image to implement the Decoder interface
func (s *GifImage) Decode() (image.Image, error) {
	return gif.Decode(s.Img)
}

// Decode method to for a jpeg|jpg to implement the Decoder interface
func (s *JpegImage) Decode() (image.Image, error) {
	return jpeg.Decode(s.Img)
}

// Convert converts the given image to a png format
// dst is the new png image
func Convert(img image.Image, dst *os.File) error {
	return png.Encode(dst, img)
}

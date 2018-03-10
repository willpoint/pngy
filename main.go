package main

/**
library reads a set of src images from a dir, converts the image type
to a png format and saves in a destination directory
- supported formats are: jpeg & gif
**/
import (
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"strings"
)

// Decoder type wraps around any type of file extension supported for convertion
type Decoder interface {
	Decode() (image.Image, error)
}

type jpegImage struct {
	image *os.File
}
type gifImage struct {
	image *os.File
}

func (s *gifImage) Decode() (image.Image, error) {
	return gif.Decode(s.image)
}

func (s *jpegImage) Decode() (image.Image, error) {
	return jpeg.Decode(s.image)
}

var supported map[string]Decoder

func main() {

	log := log.New(os.Stdout, "converter: ", log.Ltime|log.Lshortfile)

	source := flag.String("src", "source", "source folder for images")
	destination := flag.String("dest", "converted", "folder to save converted images")

	flag.Parse()

	file, err := os.Open(*source)
	if err != nil {
		log.Printf("error opening file: %v", *source)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		log.Printf("error reading directory: %v", err)
	}

	for _, v := range files {
		if v.IsDir() {
			continue
		}
		imgFile, err := os.Open(*source + "/" + v.Name())
		if err != nil {
			log.Printf("could not open image file: %v", v.Name())
		}
		extension := path.Ext(v.Name())
		var imgType Decoder
		switch extension {
		case ".jpg", ".jpeg":
			imgType = &jpegImage{imgFile}
			break
		case ".gif":
			imgType = &gifImage{imgFile}
			break
		}
		decodedImage, err := imgType.Decode()

		if err != nil {
			log.Printf("error decoding the image type: %v", err)
		}

		newFile := strings.TrimSuffix(v.Name(), extension)
		err = os.MkdirAll(*destination, os.ModePerm)
		if err != nil {
			log.Printf("error creating directory: %v", err)
		}

		n, err := os.Create(*destination + "/" + newFile + ".png")
		if err != nil {
			log.Printf("error :%v", err)
		}

		err = png.Encode(n, decodedImage)
		if err != nil {
			log.Printf("err encoding converting file: %v", err)
		}
	}
}

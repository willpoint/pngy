package main

/**
library reads a set of src images from a dir, converts the image type
to a png format and saves in a destination directory
- supported formats are: jpeg & gif
**/
import (
	"flag"
	"log"
	"os"
)

// handleErr logs the error with the provided logger
// and exits the program
func handleErr(log *log.Logger, msgf string, msgArgs ...interface{}) {
	log.Printf(msgf, msgArgs...)
	os.Exit(1)
}

func main() {
	log := log.New(os.Stdout, "pngy: ", log.Ltime|log.Lshortfile)

	source := flag.String("src", "source", "source folder for images")
	destination := flag.String("dest", "converted", "folder to save converted images")
	flag.Parse()

	file, err := os.Open(*source)
	handleErr(log, "error opening file: %v", *source)

	files, err := file.Readdir(-1)
	handleErr(log, "error reading directory: %v", err)

	for _, v := range files {
		if v.IsDir() {
			continue
		}
		imgFile, err := os.Open(*source + "/" + v.Name())
		handleErr(log, "could not open image file: %v", v.Name())

		var imgType Decoder

		extension := FileExt(v.Name())

		switch extension {
		case ".jpg", ".jpeg":
			imgType = &JpegImage{imgFile}
			break
		case ".gif":
			imgType = &GifImage{imgFile}
			break
		}

		decodedImage, err := imgType.Decode()
		handleErr(log, "error decoding the image type: %v", err)

		dstName := TrimExt(v.Name(), extension)
		err = EnsureDir(*destination)
		handleErr(log, "error creating directory: %v", err)

		dstFile, err := CreateFile(*destination + "/" + dstName + ".png")
		handleErr(log, "error creating destination file: %v", err)

		err = Convert(decodedImage, dstFile)
		handleErr(log, "error converting file: %v", err)

	}
}

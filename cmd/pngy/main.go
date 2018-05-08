package main

/**
a small tool to reads a set of src images from a directory, convert the images
to a png format and save in a destination directory
- supported formats are: jpeg & gif
- does not support recursive directory structure
**/
import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/willpoint/pngy/pkg/convert"
	"github.com/willpoint/pngy/pkg/utils"
)

// handleErr logs the error with the provided logger
// and exits the program
func handleErr(err error, log *log.Logger, msgf string, msgArgs ...interface{}) {
	if err != nil {
		log.Fatalf(msgf, msgArgs...)
	}
}

func main() {
	var wg sync.WaitGroup
	log := log.New(os.Stdout, "pngy: ", log.Ltime|log.Lshortfile)

	source := flag.String("src", "source", "source folder for images")
	destination := flag.String("dest", "converted", "folder to save converted images")
	flag.Parse()

	file, err := os.Open(*source)
	handleErr(err, log, "error opening file: %v", err)

	files, err := file.Readdir(-1)
	handleErr(err, log, "error reading directory: %v", err)

	err = utils.EnsureDir(*destination)
	handleErr(err, log, "error creating directory: %v", err)

	wg.Add(len(files))

	for _, v := range files {
		if v.IsDir() {
			wg.Done() // if it's a directory - since len(files) accounts for it
			continue
		}
		srcLocation := *source + "/" + v.Name()

		dstName := utils.TrimExt(v.Name(), utils.FileExt(srcLocation))
		dstLocation := *destination + "/" + dstName + ".png"
		go work(&wg, log, srcLocation, dstLocation)
	}

	wg.Wait()
}

// work encapsulates each job that can be done concurrently
func work(wg *sync.WaitGroup, log *log.Logger, srcLocation, dstLocation string) {
	defer wg.Done()
	start := time.Now()

	imgFile, err := os.Open(srcLocation)
	handleErr(err, log, "could not open image file: %v", err)

	var imgType convert.Decoder

	extension := utils.FileExt(srcLocation)
	switch extension {
	case ".jpg", ".jpeg":
		imgType = &convert.JpegImage{imgFile}
		break
	case ".gif":
		imgType = &convert.GifImage{imgFile}
		break
	}

	decodedImage, err := imgType.Decode()
	handleErr(err, log, "error decoding the image type: %v", err)

	dstFile, err := CreateFile(dstLocation)
	handleErr(err, log, "error creating destination file: %v", err)

	err = Convert(decodedImage, dstFile)
	handleErr(err, log, "error converting file: %v", err)

	elapsed := time.Since(start)
	log.Printf("time taken: %v", elapsed)
}

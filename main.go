/**
* Package main is a program that simply converts
* source images from one directory into jpeg images
* and writes them to another directory.
* It does not search recursively into directory
 */
package main

import (
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func main() {

	srcFolder := flag.String("src", "", "The source image file")
	dstFolder := flag.String("dst", "", "folder to save converted images")
	flag.Parse()

	if *dstFolder == "" || *srcFolder == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := ensureDir(*dstFolder)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(*srcFolder)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	files, err := file.Readdir(-1)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	for _, fi := range files {
		if fi.IsDir() {
			continue // we are not going deeper than the current level
		}

		name := fi.Name()
		srcfile, err := os.Open(*srcFolder + "/" + name)
		if err != nil {
			continue // could not open file, move on to the next file
		}

		img, _, err := image.Decode(srcfile)
		srcfile.Close() // close file since content is used
		if err != nil {
			continue // decoding failed, move on to next file, but closed already opened image
		}

		newName := strings.TrimSuffix(name, path.Ext(name)) + ".png"
		dstName := *dstFolder + "/" + newName
		dstfile, err := os.Create(dstName)
		if err != nil {
			continue // creating destination file failed, move on to next file
		}

		err = png.Encode(dstfile, img)
		dstfile.Close() // close file after writing to it
		if err != nil {
			continue // move on to the next file, ensure the already created file is closed
		}
	}

}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

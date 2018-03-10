package main

import (
	"log"
	"os"
	"path"
	"strings"
)

// handleErr logs the error with the provided logger
// and exits the program
func handleErr(err error, log *log.Logger, msgf string, msgArgs ...interface{}) {
	if err != nil {
		log.Printf(msgf, msgArgs...)
		os.Exit(1)
	}
}

// EnsureDir creates a directory if the given directory does not exist
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// FileExt returns the file extension of the given name
// this string contains the dot eg. (.jpeg)
func FileExt(name string) string {
	return path.Ext(name)
}

// TrimExt removes the extension from given name
func TrimExt(name, ext string) string {
	return strings.TrimSuffix(name, ext)
}

// CreateFile creates a new file
func CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

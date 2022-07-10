package main

import (
	"archive/zip"
	"log"
	"os"
)

func writerZip()  {
	// Create archive
	zipPath := "out.zip"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new zip archive.
	w := zip.NewWriter(zipFile)
	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"asong.txt", "This archive contains some text files."},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
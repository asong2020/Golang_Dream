package main

import (
	"io/ioutil"
	"log"
	"os"
)

func readAll(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	log.Printf("read %s content is %s", filename, data)
	return nil
}

func ReadAll2(filename string) error {
	file, err := os.Open("asong.txt")
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	log.Printf("read %s content is %s\n", filename, content)

	file.Close()
	return nil
}
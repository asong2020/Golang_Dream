package main

import (
	"os"
)

func writeAll(filename string) error {
	err := os.WriteFile("asong.txt", []byte("Hi asong\n"), 0666)
	if err != nil {
		return err
	}
	return nil
}

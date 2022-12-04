package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func readLine(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bufferedReader := bufio.NewReader(file)
	for {
		// ReadLine is a low-level line-reading primitive. Most callers should use
		// ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		lineBytes, err := bufferedReader.ReadBytes('\n')
		bufferedReader.ReadLine()
		line := strings.TrimSpace(string(lineBytes))
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		log.Printf("readline %s every line data is %s\n", filename, line)
	}
	file.Close()
	return nil
}
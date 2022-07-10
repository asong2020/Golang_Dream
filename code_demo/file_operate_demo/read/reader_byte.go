package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// use bufio.NewReader
func readByte(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	// 创建 Reader
	r := bufio.NewReader(file)

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}

// use os
func readByte2(filename string) error{
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}


// use os and io.ReadAtLeast
func readByte3(filename string) error{
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	// 每次读取 2 个字节
	buf := make([]byte, 2)
	for {
		n, err := io.ReadAtLeast(file, buf, 0)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}
		log.Printf("writeByte %s every read 2 byte is %s\n", filename, string(buf[:n]))
	}
	file.Close()
	return nil
}
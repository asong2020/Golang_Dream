package main

import (
	"bufio"
	"log"
	"os"
)

func writeLine(filename string) error {
	data := []string{
		"asong",
		"test",
		"123",
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil{
		return err
	}

	for _, line := range data{
		_,err := f.WriteString(line + "\n")
		if err != nil{
			return err
		}
	}
	f.Close()
	return nil
}


func writeLine2(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// 为这个文件创建buffered writer
	bufferedWriter := bufio.NewWriter(file)

	for i:=0; i < 2; i++{
		// 写字符串到buffer
		bytesWritten, err := bufferedWriter.WriteString(
			"asong真帅\n",
		)
		if err != nil {
			return err
		}
		log.Printf("Bytes written: %d\n", bytesWritten)
	}
	// 写内存buffer到硬盘
	err = bufferedWriter.Flush()
	if err != nil{
		return err
	}

	file.Close()
	return nil
}
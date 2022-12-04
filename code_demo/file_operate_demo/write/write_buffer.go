package main

import (
	"bufio"
	"log"
	"os"
)

func writeBuffer(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// 为这个文件创建buffered writer
	bufferedWriter := bufio.NewWriter(file)

	// 写字符串到buffer
	bytesWritten, err := bufferedWriter.WriteString(
		"asong真帅\n",
	)
	if err != nil {
		return err
	}
	log.Printf("Bytes written: %d\n", bytesWritten)

	// 检查缓存中的字节数
	unflushedBufferSize := bufferedWriter.Buffered()
	log.Printf("Bytes buffered: %d\n", unflushedBufferSize)

	// 还有多少字节可用（未使用的缓存大小）
	bytesAvailable := bufferedWriter.Available()
	if err != nil {
		return err
	}
	log.Printf("Available buffer: %d\n", bytesAvailable)
	// 写内存buffer到硬盘
	err = bufferedWriter.Flush()
	if err != nil{
		return err
	}

	file.Close()
	return nil
}
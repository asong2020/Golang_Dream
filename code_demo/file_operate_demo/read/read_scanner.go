package main

import (
	"bufio"
	"log"
	"os"
)

func readScanner(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	// 可以定制Split函数做分隔函数
	// ScanWords 是scanner自带的分隔函数用来找空格分隔的文本字
	scanner.Split(bufio.ScanWords)
	for {
		success := scanner.Scan()
		if success == false {
			// 出现错误或者EOF是返回Error
			err = scanner.Err()
			if err == nil {
				log.Println("Scan completed and reached EOF")
				break
			} else {
				return err
			}
		}
		// 得到数据，Bytes() 或者 Text()
		log.Printf("readScanner get data is %s", scanner.Text())
	}
	file.Close()
	return nil
}
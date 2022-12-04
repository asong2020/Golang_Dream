package main

import (
	"log"
	"os"
)

const name = "asong.txt"

func main()  {
	// 构造数据
	err := writeLine(name)
	if err != nil{
		log.Fatalf("write data failed err=%s\n", err)
	}
	// 读取文件所有内容
	err = readAll(name)
	if err != nil{
		log.Fatalf("read all data failed err=%s\n", err)
	}

	// 按行读取
	err = readLine(name)
	if err != nil{
		log.Fatalf("read line data failed err=%s\n", err)
	}

	// 按字节读取
	err = readByte(name)
	if err != nil{
		log.Fatalf("read byte failed err=%s\n", err)
	}

	// scanner
	err = readScanner(name)
	if err != nil{
		log.Fatalf("read scanner failed err=%s\n", err)
	}
}


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
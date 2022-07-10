package main

import (
	"io/ioutil"
	"log"
	"os"
)

const name = "asong.txt"

func main()  {
	// write all
	err := writeAll(name)
	if err != nil{
		log.Fatalf("write all failed err=%s\n", err)
	}
	err = PrintFileAllContent(name)
	if err != nil{
		log.Fatalf("printe file all content failed err=%s\n", err)
	}

	err = EmptyFile(name)
	if err != nil{
		log.Fatalf("emptyFile failed err=%s\n", err)
	}
	// write at
	err = writeAt(name)
	if err != nil{
		log.Fatalf("write at failed err=%s\n", err)
	}

	err = PrintFileAllContent(name)
	if err != nil{
		log.Fatalf("printe file all content failed err=%s\n", err)
	}

	err = EmptyFile(name)
	if err != nil{
		log.Fatalf("emptyFile failed err=%s\n", err)
	}

	// write line
	err = writeLine(name)
	if err != nil{
		log.Fatalf("write line failed err=%s\n", err)
	}

	err = PrintFileAllContent(name)
	if err != nil{
		log.Fatalf("printe file all content failed err=%s\n", err)
	}

	err = EmptyFile(name)
	if err != nil{
		log.Fatalf("emptyFile failed err=%s\n", err)
	}

	// write buffer
	err = writeBuffer(name)
	if err != nil{
		log.Fatalf("write buffer failed err=%s\n", err)
	}

	err = PrintFileAllContent(name)
	if err != nil{
		log.Fatalf("printe file all content failed err=%s\n", err)
	}

	err = EmptyFile(name)
	if err != nil{
		log.Fatalf("emptyFile failed err=%s\n", err)
	}
}

func PrintFileAllContent(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	log.Printf("%s file content is %s", filename, data)
	return nil
}


func EmptyFile(filename string) error {
	file, _ := os.OpenFile(filename, os.O_RDWR, 0666)
	// 清空文件
	err  := file.Truncate(0)
	if err != nil{
		return err
	}
	// 重置当前位置
	_, err = file.Seek(0, 0)
	if err != nil{
		return err
	}
	return nil
}
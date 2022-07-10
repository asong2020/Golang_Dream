package main

import "os"

func writeAt(filename string) error {
	data := []byte{
		0x41, // A
		0x73, // s
		0x20, // space
		0x20, // space
		0x67, // g
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil{
		return err
	}
	_, err = f.Write(data)
	if err != nil{
		return err
	}

	replaceSplace := []byte{
		0x6F, // o
		0x6E, // n
	}
	_, err = f.WriteAt(replaceSplace, 2)
	if err != nil{
		return err
	}
	f.Close()
	return nil
}
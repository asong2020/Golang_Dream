package main

import (
	"bytes"
	"fmt"
	"strings"
)

const base = "123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASFGHJKLZXCVBNM"
var baseSlice []string

func init()  {
	for i := 0; i < 200; i++ {
		baseSlice = append(baseSlice, base)
	}
}


func main()  {
}

func SumString()  string{
	res := ""
	for _,val := range baseSlice{
		res += val
	}
	return res
}

func SprintfString() string {
	res := ""
	for _,val := range baseSlice{
		res = fmt.Sprintf("%s%s", res, val)
	}
	return res
}

func BuilderString() string {
	var builder strings.Builder
	builder.Grow(200 * len(baseSlice))
	for _,val := range baseSlice{
		builder.WriteString(val)
	}
	return builder.String()
}

func bytesString() string {
	buf := new(bytes.Buffer)
	for _,val := range baseSlice{
		buf.WriteString(val)
	}
	return buf.String()
}

func byteSliceString() string {
	buf := make([]byte, 0)
	for _,val := range baseSlice{
		buf = append(buf, val...)
	}
	return string(buf)
}

func Joinstring() string {
	return strings.Join(baseSlice, "")
}

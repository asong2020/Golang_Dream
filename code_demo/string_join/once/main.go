package main

import (
	"bytes"
	"fmt"
	"strings"
)

var base  = "123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASFGHJKLZXCVBNM"

func main()  {
	fmt.Println(SumString(base))
	fmt.Println(SprintfString(base))
	fmt.Println(BuilderString(base))
	fmt.Println(bytesString(base))
	fmt.Println(byteSliceString(base))
	fmt.Println(Joinstring())
}

func SumString(str string)  string{
	return base + str
}

func SprintfString(str string) string {
	return fmt.Sprintf("%s%s", base, str)
}

func BuilderString(str string) string {
	var builder strings.Builder
	builder.Grow(2 * len(str))
	builder.WriteString(base)
	builder.WriteString(str)
	return builder.String()
}

func bytesString(str string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(base)
	buf.WriteString(str)
	return buf.String()
}

func byteSliceString(str string) string {
	buf := make([]byte, 0)
	buf = append(buf, base...)
	buf = append(buf, str...)
	return string(buf)
}

func Joinstring() string {
	return strings.Join([]string{base, base}, "")
}
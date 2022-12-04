package main

import "fmt"

type Base struct {
	Value string
}

func (b *Base) GetMsg() string {
	return b.Value
}


type Person struct {
	Base
	Name string
	Age uint64
}

func (p *Person) GetName() string {
	return p.Name
}

func (p *Person) GetAge() uint64 {
	return p.Age
}

func check(b *Base)  {
	b.GetMsg()
}

func main()  {
	m := Base{Value: "I Love You"}
	p := &Person{
		Base: m,
		Name: "asong",
		Age: 18,
	}
	fmt.Print(p.GetName(), "  ", p.GetAge(), " and say ",p.GetMsg())
	//check(p)
}

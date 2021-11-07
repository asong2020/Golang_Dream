package main

type Hero struct {
	Name string
	Age uint64
}

func NewHero() *Hero {
	return &Hero{
		Name: "盖伦",
		Age: 18,
	}
}

func (h *Hero) GetName() string {
	return h.Name
}

func (h *Hero) GetAge() uint64 {
	return h.Age
}


func main()  {
	h := NewHero()
	print(h.GetName())
	print(h.GetAge())
}

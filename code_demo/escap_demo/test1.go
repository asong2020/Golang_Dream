package main

func Add(x,y int) *int {
	res := 0
	res = x + y
	return &res
}

func main()  {
	Add(1,2)
}

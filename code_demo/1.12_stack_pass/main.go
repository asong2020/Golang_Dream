package main

func Test(a, b int) (int, int) {
	return a + b, a - b
}

func main() {
	Test(10, 20)
}

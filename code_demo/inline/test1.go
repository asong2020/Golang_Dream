package main
//
//func main(){
//	s := []int{90,100,24,18}
//	Sum(s)
//}

func Sum(s []int) int {
	sum :=0
	for i:=0;i<len(s);i++{
		sum = add(sum,s[i])
	}
	return sum
}

func add(x,y int) int{
	panic("panic")
	return x+y
}
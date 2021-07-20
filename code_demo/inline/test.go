package main

//func main(){
//	s := []int{10,12,3,14}
//	fmt.Println(GetMaxValue(s))
//}

func GetMaxValue(s []int) int {
	max :=0
	for i:=0;i<len(s);i++{
		max = maxValue(s[i],max)
	}
	return max
}

func maxValue(a,b int) int {
	if a > b{
		return a
	}
	return b
}

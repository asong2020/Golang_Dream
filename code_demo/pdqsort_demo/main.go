package main


import (
	"fmt"

	"github.com/zhangyunhao116/pdqsort"
)

func main() {
	x := []int{3, 1, 2, 4, 5, 9, 8, 7}
	pdqsort.Slice(x)
	fmt.Printf("sort_result = %v\n", x)
	search_result := pdqsort.Search(x, 4)
	fmt.Printf("search_result = %v\n", search_result)
	is_sort := pdqsort.SliceIsSorted(x)
	fmt.Printf("is_sort = %v\n", is_sort)
}

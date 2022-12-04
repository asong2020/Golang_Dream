package main

import (
	"fmt"
)

// 查找一个数在数组中下标的位置 (递归实现)
func Search(nums []int64, target int64) int {
	return binarySearchRecursive(target, nums, 0, len(nums))
}

// 二分查找递归实现
func binarySearchRecursive(target int64, nums []int64, left, right int) int {
	if left > right {
		return -1
	}

	mid := left + (right - left) / 2

	if target == nums[mid] {
		return mid
	}
	if nums[mid] < target {
		return binarySearchRecursive(target, nums, mid+1, right)
	}
	if nums[mid] > target {

		return binarySearchRecursive(target, nums, left, mid-1)
	}
	return -1

}

// 二分查找非递归实现
func binarySearch(target int64, nums []int64) int {

	left := 0
	right := len(nums)
	for left <= right {
		mid := left + (right - left) / 2

		if target == nums[mid] {
			return mid
		}

		if target > nums[mid] {
			left = mid + 1
			continue
		}

		if target < nums[mid] {
			right = mid - 1
			continue
		}

	}

	return -1

}

func main() {
	//nums := []int64{1, 2, 3, 4, 5, 6, 7}
	//
	////fmt.Println(Search(nums, 1))
	////fmt.Println(binarySearch(1, nums))
	//fmt.Println(sort.Search(len(nums), func(i int) bool{
	//	return nums[i] >= 10
	//}))
	var a uint = 100
	var b uint = 10
	fmt.Println(b - a)
}
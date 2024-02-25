package main

import (
	"fmt"
)

func mergeArrays(arr1 []int, arr2 []int) []int {
	result := make([]int, len(arr1)+len(arr2))
	i, j, k := 0, 0, 0

	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			result[k] = arr1[i]
			i++
		} else {
			result[k] = arr2[j]
			j++
		}
		k++
	}

	for i < len(arr1) {
		result[k] = arr1[i]
		i++
		k++
	}

	for j < len(arr2) {
		result[k] = arr2[j]
		j++
		k++
	}

	return result
}

func main() {
	arr1 := []int{1, 3, 5, 7}
	arr2 := []int{2, 4, 6, 8, 9}

	mergedArray := mergeArrays(arr1, arr2)

	fmt.Println("Merged Array:", mergedArray)
}

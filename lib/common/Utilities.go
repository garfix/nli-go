package common

import (
	"runtime"
	"path"
)

func IntArrayContains(haystack []int, needle int) bool {
	for _, value := range haystack  {
		if needle == value {
			return true
		}
	}
	return false
}

func IntArrayDeduplicate(array []int) []int {
	newArray := []int{}
	previous := map[int]bool{}

	for _, value := range array {
		_, found := previous[value]
		if !found {
			previous[value] = true
			newArray = append(newArray, value)
		}
	}

	return newArray
}

// Returns the directory of the file calling this function
// https://coderwall.com/p/_fmbug/go-get-path-to-current-file
func GetCurrentDir() string {

	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

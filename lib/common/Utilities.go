package common

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func IntArrayCopy(original []int) []int {
	copiedArray := make([]int, len(original))
	copy(copiedArray, original)
	return copiedArray
}

func StringArrayCopy(original []string) []string {
	copiedArray := []string{}
	for _, element := range original {
		copiedArray = append(copiedArray, element)
	}
	return copiedArray
}

func StringMatrixCopy(original [][]string) [][]string {
	copiedArray := [][]string{}
	for _, array := range original {
		anArray := []string{}
		for _, element := range array {
			anArray = append(anArray, element)
		}
		copiedArray = append(copiedArray, anArray)
	}
	return copiedArray
}

func IntArrayContains(haystack []int, needle int) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}

func StringArrayContains(haystack []string, needle string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}

func IntArrayEquals(haystack []int, needle []int) bool {
	if len(haystack) != len(needle) {
		return false
	}

	for i := range haystack {
		if haystack[i] != needle[i] {
			return false
		}
	}
	return true
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

func StringArrayDeduplicate(array []string) []string {
	newArray := []string{}
	uniques := map[string]bool{}

	for _, value := range array {
		_, found := uniques[value]
		if !found {
			uniques[value] = true
			newArray = append(newArray, value)
		}
	}

	return newArray
}

// Returns the directory of the file calling this function
// https://coderwall.com/p/_fmbug/go-get-path-to-current-file
func Dir() string {

	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func ReadFile(path string) (string, error) {

	source := ""

	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		source = string(bytes)
	}

	return source, err
}

func WriteFile(path string, contents string) error {
	return ioutil.WriteFile(path, []byte(contents), 0666)
}

// If path is an absolute path, returns path
// Otherwise, adds path to baseDir to create an absolute path
func AbsolutePath(baseDir string, path string) string {

	absolutePath := path

	if len(path) > 0 && path[0] != os.PathSeparator {
		absolutePath, _ = filepath.Abs(baseDir + string(os.PathSeparator) + path)
	}

	return absolutePath
}

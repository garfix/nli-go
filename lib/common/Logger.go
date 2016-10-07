package common

import "fmt"

var active = false

func Log(text string) {
	if active {
		fmt.Print(text)
	}
}

func Logf(text string, vals ...interface{}) {
	Log(fmt.Sprintf(text, vals))
}

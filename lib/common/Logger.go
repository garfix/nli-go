package common

import "fmt"

var LoggerActive = false

func Log(text string) {
	if LoggerActive {
		fmt.Print(text)
	}
}

func Logf(text string, vals ...interface{}) {
	Log(fmt.Sprintf(text, vals...))
}

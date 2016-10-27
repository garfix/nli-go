package common

import (
	"fmt"
	"strings"
)

var LoggerActive = false

var logStack = []string{}

func Log(text string) {
	if LoggerActive {
		fmt.Print(text)
	}
}

func Logf(text string, vals ...interface{}) {
	Log(fmt.Sprintf(text, vals...))
}

func LogTree(text string, vals ...interface{}) {

	if !LoggerActive {
		return
	}

	sameText := len(logStack) > 0 && text == logStack[len(logStack) - 1]

	if !sameText {
		logStack = append(logStack, text)
	}

	stmt := strings.Repeat("  ", len(logStack)) + text + " "
	for _, val := range vals {
		stmt += fmt.Sprintf("%v", val) + " "
	}

	Log(stmt + "\n")

	if sameText {
		logStack = logStack[0:len(logStack) - 1]
	}

}

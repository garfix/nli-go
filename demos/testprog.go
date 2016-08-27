package main

import (
	"nli-go/lib/example3"
	"fmt"
)

func main() {

	transformations := []example3.SimpleRelationTransformation{}
	ok := true
	lastLine := 0

	parser := example3.NewSimpleRelationTransformationParser()
	transformations, lastLine, ok = parser.ParseString("father(A, B) :- parent(A, B), male(A)")

	if !ok {
		fmt.Print("Parse error")
	}
	if lastLine != 1 {
		fmt.Printf("Error in line: %d", lastLine)
	}
	if len(transformations) != 1 {
		fmt.Printf("Wrong number of transformations: %d", len(transformations))
	}
}

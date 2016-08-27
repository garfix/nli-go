package tests

import (
	"testing"
	"nli-go/lib/example3"
	"fmt"
)

func TestSimpleGrammarTransformationParser(test *testing.T) {

	parser := example3.NewSimpleRelationTransformationParser()
	transformations := []example3.SimpleRelationTransformation{}
	ok := true
	lastLine := 0

	transformations, lastLine, ok = parser.ParseString("father(A, B) :- parent(A, B), male(A)")
	if !ok {
		test.Error("Parse error")
	}
	if lastLine != 1 {
		test.Error(fmt.Printf("Error in line: %d", lastLine))
	}
	if len(transformations) != 1 {
		test.Error(fmt.Printf("Wrong number of transformations: %d", len(transformations)))
	}

	transformations, lastLine, ok = parser.ParseString("father(A, B) :- ")
	if ok {
		test.Error("Parse should have failed")
	}

	transformations, lastLine, ok = parser.ParseString(":- parent(A, B), male(A)")
	if ok {
		test.Error("Parse should have failed")
	}
	transformations, lastLine, ok = parser.ParseString("\n")
	if !ok {
		test.Error("Parse should have succeeded")
	}

	transformations, lastLine, ok = parser.ParseString("\n" +
		"\tfather(A, B) :- parent(A, B), male(A)\n" +
		"\tcommand(A1), owner(Owner), done(true), fixed() :- tell(A1, Owner), prime_number(7)\n")
	if !ok {
		test.Error("Parse error")
	}
	if lastLine != 3 {
		test.Error(fmt.Printf("Last line was: %d", lastLine))
	}
	if len(transformations) != 2 {
		test.Error(fmt.Printf("Wrong number of transformations: %d", len(transformations)))
	}
}
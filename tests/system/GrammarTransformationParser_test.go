package tests

import (
	"testing"
	"fmt"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
)

func TestGrammarTransformationParser(test *testing.T) {

	parser := importer.NewInternalGrammarParser()
	transformations := []mentalese.RelationTransformation{}

	transformations = parser.CreateTransformations("[ father(A, B) :- parent(A, B) male(A); ]")
	if !parser.GetLastParseResult().Ok {
		test.Error("Parse should have succeeded")
	}
	if len(transformations) != 1 {
		test.Error(fmt.Printf("Wrong number of transformations: %d", len(transformations)))
	}

	transformations = parser.CreateTransformations("[\n]")
	if !parser.GetLastParseResult().Ok {
		test.Error("Parse should have succeeded")
	}

	transformations = parser.CreateTransformations("father(A, B) :- ")
	if parser.GetLastParseResult().Ok {
		test.Error("Parse should have failed")
	}

	transformations = parser.CreateTransformations(":- parent(A, B), male(A)")
	if parser.GetLastParseResult().Ok {
		test.Error("Parse should have failed")
	}

	transformations = parser.CreateTransformations("[\n" +
		"father(A, B) :- parent(A, B) male(A);\n" +
		"command(A1) owner(Owner) done(true) fixed() :- tell(A1, Owner) prime_number(7);\n" +
		"]")
	if !parser.GetLastParseResult().Ok {
		test.Error("Parse error")
	}
	if parser.GetLastParseResult().LineNumber != 4 {
		test.Error(fmt.Printf("Last line was: %d", parser.GetLastParseResult().LineNumber))
	}
	if len(transformations) != 2 {
		test.Error(fmt.Printf("Wrong number of transformations: %d", len(transformations)))
	}
}
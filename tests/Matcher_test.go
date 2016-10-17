package tests

import (
	"testing"
	"nli-go/lib/mentalese"
	"nli-go/lib/importer"
)

func TestMatchTwoRelations(t *testing.T) {
	var tests = []struct {
		needle string
		haystack string
		binding string
		wantBinding string
		wantMatch bool
	} {
		{"parent(X, Y)", "parent('Luke', 'George')", "{}", "{X: 'Luke', Y: 'George'}", true},
		{"parent('Luke', 'George')", "parent(X, Y)", "{}", "{}", true},
		{"parent('Luke', Y)", "parent('Luke', 'George')", "{}", "{Y: 'George'}", true},
		{"parent('Luke', 'Richard')", "parent('Luke', 'George')", "{}", "{}", false},
		{"parent(X, Y)", "parent(A, B)", "{}", "{X:A, Y:B}", true},
		{"parent(X, Y)", "parent('Luke', 'George')", "{X: 'Luke'}", "{X: 'Luke', Y: 'George'}", true},
		{"parent(X, Y)", "parent('Luke', 'George')", "{X: 'Vincent'}", "{X: 'Vincent'}", false},
	}

	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher()

	for _, test := range tests {

		needle, _ := parser.CreateRelation(test.needle)
		haystack, _ := parser.CreateRelation(test.haystack)
		binding, _ := parser.CreateBinding(test.binding)
		wantBinding, _ := parser.CreateBinding(test.wantBinding)
		wantMatch := test.wantMatch

		resultBinding, resultMatch := matcher.MatchTwoRelations(needle, haystack, binding)

		if !resultBinding.Equals(wantBinding) || resultMatch != wantMatch {
			t.Errorf("MatchTwoRelations(%v %v %v): got %v %v, want %v %v", needle, haystack, binding, resultBinding, resultMatch, wantBinding, wantMatch)
		}
	}
}
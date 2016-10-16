package tests

import (
	"testing"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
)

func TestBinder(t *testing.T) {
	var tests = []struct {
		subject string
		object string
		binding string
		want string
	} {
		{"A", "13", "{}", "{A:13}"},
	}

	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher()

	for _, test := range tests {

		subject, _ := parser.CreateTerm(test.subject)
		object, _ := parser.CreateTerm(test.object)
		binding, _ := parser.CreateBinding(test.binding)
		want, _ := parser.CreateBinding(test.want)
		result, _ := matcher.BindTerm(subject, object, binding)

		if !result.Equals(want) {
			t.Errorf("bindTerm(%v, %v, %v): got %v, want %v", subject, object, binding, result, test.want)
		}
	}
}
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
		wantBinding string
		wantOk bool
	} {
		{"A", "13", "{}", "{A:13}", true},
		{"A", "B", "{}", "{A:B}", true},
		{"13", "B", "{}", "{}", true},
		{"A", "13", "{B: 14}", "{B:14, A:13}", true},
		{"A", "13", "{A:6}", "{A:6}", false},
	}

	parser := importer.NewInternalGrammarParser()
	matcher := mentalese.NewRelationMatcher()

	for _, test := range tests {

		subject, _ := parser.CreateTerm(test.subject)
		object, _ := parser.CreateTerm(test.object)
		binding, _ := parser.CreateBinding(test.binding)
		originalBinding := binding
		wantBinding, _ := parser.CreateBinding(test.wantBinding)
		resultBinding, resultOk := matcher.BindTerm(subject, object, binding)

		if !resultBinding.Equals(wantBinding) || resultOk != test.wantOk {
			t.Errorf("bindTerm(%v, %v, %v): got %v %v, want %v %v", subject, object, binding, resultBinding, resultOk, wantBinding, test.wantOk)
		}
		if !binding.Equals(originalBinding) {
			t.Errorf("bindTerm input changed: got %v, want %v", binding, originalBinding)
		}
	}
}
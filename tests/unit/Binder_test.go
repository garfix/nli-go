package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestBinder(t *testing.T) {
	var tests = []struct {
		subject     string
		object      string
		binding     string
		wantBinding string
		wantOk      bool
	}{
		{"A", "13", "{}", "{A:13}", true},
		{"A", "B", "{}", "{A:B}", true},
		{"'John'", "'John'", "{}", "{}", true},
		{"'John'", "'Jack'", "{}", "{}", false},
		{"21", "_", "{}", "{}", true},
		{"_", "21", "{}", "{}", true},
		{"13", "B", "{}", "{}", true},
		{"A", "13", "{B: 14}", "{B:14, A:13}", true},
		{"A", "13", "{A:6}", "{A:6}", false},
	}

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)

	for _, test := range tests {

		subject := parser.CreateTerm(test.subject)
		object := parser.CreateTerm(test.object)
		binding := parser.CreateBinding(test.binding)
		originalBinding := binding
		wantBinding := parser.CreateBinding(test.wantBinding)
		resultBinding, resultOk := matcher.BindTerm(subject, object, binding)

		if !resultBinding.Equals(wantBinding) || resultOk != test.wantOk {
			t.Errorf("bindTerm(%v, %v, %v): got %v %v, want %v %v", subject, object, binding, resultBinding, resultOk, wantBinding, test.wantOk)
		}
		if !binding.Equals(originalBinding) {
			t.Errorf("bindTerm input changed: got %v, want %v", binding, originalBinding)
		}
	}
}

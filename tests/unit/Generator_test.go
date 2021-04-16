package tests

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog()

	grammarRules := internalGrammarParser.CreateGenerationGrammar(`
        { rule: s(P) -> np(E) vp(P),              condition: grammatical_subject(E) subject(P, E) }
		{ rule: s(P) -> named_number(P),          condition: result(P) }
        { rule: np(E) -> proper_noun(E),          condition: name(E, Name) }
		{ rule: proper_noun(E) -> text(Name),     condition: name(E, Name) }
        { rule: np(E) -> det(E) noun(E) }
        { rule: vp(V) -> verb(V) np(E),           condition: object(V, E) }
		{ rule: noun(E) -> 'book',                condition: instance_of(E, book) }
		{ rule: verb(E) -> 'kissed',		      condition: predication(E, kiss) }
		{ rule: verb(E) -> 'married',		      condition: predication(E, marry) }
		{ rule: named_number(1) -> 'one' }
		{ rule: named_number(2) -> 'two' }
	`)
	matcher := central.NewRelationMatcher(log)
	meta := mentalese.NewMeta()
	matcher.AddFunctionBase(knowledge.NewSystemFunctionBase("system-function", meta, log))
	generator := generate.NewGenerator(log, matcher)

	tests := []struct {
		input string
		want  string
	}{
		{"predication(P1, marry) subject(P1, E1) object(P1, E2) name(E1, 'John') name(E2, 'Mary') grammatical_subject(E1)", "John married Mary"},
		{"result(2)", "two"},
	}

	for _, test := range tests {

		input := internalGrammarParser.CreateRelationSet(test.input)
		result := generator.Generate(grammarRules, input)
		if strings.Join(result, " ") != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.input, strings.Join(result, " "), test.want)
			fmt.Println(log.String())
		}
	}
}

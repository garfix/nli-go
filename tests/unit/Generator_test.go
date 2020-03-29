package tests

import (
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
	log := common.NewSystemLog(false)

	grammar := internalGrammarParser.CreateGenerationGrammar(`[
        { rule: s(P) -> np(E) vp(P),              condition: grammatical_subject(E) subject(P, E) }
        { rule: np(E) -> proper_noun(E),          condition: name(E, Name) }
		{ rule: proper_noun(E) -> text(Name),     condition: name(E, Name) }
        { rule: np(E) -> det(E) noun(E) }
        { rule: vp(V) -> verb(V) np(E),           condition: object(V, E) }
		{ rule: noun(E) -> 'book',                condition: instance_of(E, book) }
		{ rule: verb(E) -> 'kissed',		      condition: predication(E, kiss) }
		{ rule: verb(E) -> 'married',		      condition: predication(E, marry) }
	]`)
	lexicon := internalGrammarParser.CreateGenerationLexicon(`[
		
	]`, log)
	matcher := mentalese.NewRelationMatcher(log)
	matcher.AddFunctionBase(knowledge.NewSystemFunctionBase("system-function"))
	generator := generate.NewGenerator(grammar, lexicon, log, matcher)

	tests := []struct {
		input string
		want  string
	}{
		{"[predication(P1, marry) subject(P1, E1) object(P1, E2) name(E1, 'John') name(E2, 'Mary') grammatical_subject(E1)]", "John married Mary"},
	}

	for _, test := range tests {

		input := internalGrammarParser.CreateRelationSet(test.input)
		result := generator.Generate(input)
		if strings.Join(result, " ") != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.input, strings.Join(result, " "), test.want)
		}
	}
}

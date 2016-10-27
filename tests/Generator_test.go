package tests

import (
	"testing"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"strings"
	"nli-go/lib/common"
)

func TestGenerator(t *testing.T) {

	common.LoggerActive=false

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar, _, _ := internalGrammarParser.CreateGenerationGrammar(`[
	    {
	        rule: s(P) :- np(E), vp(P)
	        condition: grammatical_subject(E), subject(P, E)
	    } {
	        rule: np(E) :- proper_noun(E)
	    } {
	        rule: np(E) :- det(E), noun(E)
	    }{
	        rule: vp(V) :- verb(V), proper_noun(E)
	        condition: object(V, E)
	    }
	]`)
	lexicon, _, _ := internalGrammarParser.CreateGenerationLexicon(`[
		{
			form: 'book'
			pos: noun
			condition: instance_of(This, book)
		}
		{
			form: 'married'
			pos: verb
			condition: predication(This, marry)
		}
		{
			form: 'name'
			pos: proper_noun
			condition: name(This, Name)
		}
	]`)
	generator := generate.NewGenerator(grammar, lexicon)

	tests := []struct {
		input string
		want string
	} {
		{"[predication(P1, marry) subject(P1, E1) object(P1, E2) name(E1, 'John') name(E2, 'Mary') grammatical_subject(E1)]", "John married Mary"},
	}

	for _, test := range tests {

		input, _, _ := internalGrammarParser.CreateRelationSet(test.input)
		result := generator.Generate(input)
		if strings.Join(result, " ") != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.input, strings.Join(result, " "), test.want)
		}
	}
}
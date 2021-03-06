package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/function"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"testing"
)

func TestGeneralizedQuantifier(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog()

	grammarRules := internalGrammarParser.CreateGrammarRules(`
		{ rule: qp(_) -> quantifier(Result, Range),                     	sense: go:quantifier(Result, Range, $quantifier) }
		{ rule: quantifier(Result, Range) -> 'all', 						sense: go:equals(Result, Range) } 
		{ rule: quantifier(Result, Range) -> 'some', 						sense: go:greater_than(Result, 0) }
		{ rule: quantifier(Result, Range) -> 'no', 							sense: go:equals(Result, 0) }
		{ rule: quantifier(Result, Range) -> number(N1), 	    			sense: go:equals(Result, N1) }
		{ rule: quantifier(Result, Range) -> 'more' 'than' number(N1),		sense: go:greater_than(Result, N1) }
		{ rule: quantifier(Result, Range) -> quantifier(Result, Range) 'or' quantifier(Result, Range),	
																			sense: go:or($quantifier1, $quantifier2) }

		{ rule: number(N1) -> ~^[0-9]+~ }

		{ rule: nbar(E1) -> 'books', 										sense: book(E1) }
		{ rule: np(E1) -> qp(_) nbar(E1), 									sense: go:quant($qp, E1, $nbar) }
		{ rule: s(S1) -> 'did' 'abraham' 'read' np(E1),     				sense: go:quant_check($np, read('abraham', E1)) }
	`)

	facts := internalGrammarParser.CreateRelationSet(`
		book('Dracula')
		book('Frankenstein')
		book('Curse of the mummy')

		read('abraham', 'Dracula')
		read('abraham', 'Frankenstein')
		read('abraham', 'Curse of the mummy')
		read('sarah', 'Dracula')
	`)

	readMap := internalGrammarParser.CreateRules(`
		book(A) :- book(A);
		read(A, B) :- read(A, B);
	`)
	writeMap := []mentalese.Rule{}

	matcher := central.NewRelationMatcher(log)
	meta := mentalese.NewMeta()
	solver := central.NewProblemSolverAsync(matcher, log)
	factBase := knowledge.NewInMemoryFactBase("in-memory", facts, matcher, readMap, writeMap, nil, log)
	solver.AddFactBase(factBase)
	anaphoraQueue := central.NewAnaphoraQueue()
	deicticCenter := central.NewDeicticCenter()
	nestedStructureBase := function.NewSystemSolverFunctionBase(anaphoraQueue, deicticCenter, meta, log)
	solver.AddSolverFunctionBase(nestedStructureBase)
	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", meta, log)
	solver.AddFunctionBase(systemFunctionBase)
	solver.Reindex()
	runner := central.NewProcessRunner(solver, log)
	tokenizer := parse.NewTokenizer(parse.DefaultTokenizerExpression)
	parser := parse.NewParser(grammarRules, log)

	tests := []struct {
		input string
		want  string
	}{
		{"did Abraham read all books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read 3 books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read 2 books", "[]"},
		{"did Abraham read more than 2 books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read more than 3 books", "[]"},
		{"did Abraham read 2 or 3 books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read 3 or 4 books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read 4 or 5 books", "[]"},
		{"did Abraham read some books", "[{E5:'Dracula'} {E5:'Frankenstein'} {E5:'Curse of the mummy'}]"},
		{"did Abraham read no books", "[]"},
	}

	for _, test := range tests {

		words := tokenizer.Process(test.input)
		trees := parser.Parse(words, "s", []string{"S"})
		if len(trees) == 0 {
			t.Errorf("Cannot parse: %s", test.input)
			continue
		}
		relationizer := parse.NewRelationizer(log)
		input, _ := relationizer.Relationize(trees[0], []string{ "S"})
		result := runner.RunRelationSetWithBindings(input, mentalese.InitBindingSet( mentalese.NewBinding() ))
		if result.String() != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.input, result.String(), test.want)
			print(log.String())
		}
		log.Clear()
	}
}
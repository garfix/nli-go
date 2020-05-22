package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/nested"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestGeneralizedQuantifier(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)

	grammar := internalGrammarParser.CreateGrammar(`[
		{ rule: qp(_) -> quantifier(Result, Range),                     	sense: quantifier(Result, Range, sem(1)) }
		{ rule: quantifier(Result, Range) -> 'all', 						sense: equals(Result, Range) } 
		{ rule: quantifier(Result, Range) -> 'some', 						sense: greater_than(Result, 0) }
		{ rule: quantifier(Result, Range) -> 'no', 							sense: equals(Result, 0) }
		{ rule: quantifier(Result, Range) -> number(N1), 	    			sense: equals(Result, N1) }
		{ rule: quantifier(Result, Range) -> 'more' 'than' number(N1),		sense: greater_than(Result, N1) }
		{ rule: quantifier(Result, Range) -> quantifier(Result, Range) 'or' quantifier(Result, Range),	
																			sense: or(P1, sem(1), sem(3)) }

		{ rule: number(N1) -> /^[0-9]+/ }

		{ rule: nbar(E1) -> 'books', 										sense: book(E1) }
		{ rule: np(E1) -> qp(_) nbar(E1), 									sense: quant(sem(1), E1, sem(2)) }
		{ rule: s(S1) -> 'did' 'abraham' 'read' np(E1),     				sense: find(sem(4), read('abraham', E1)) }
	]`)

	facts := internalGrammarParser.CreateRelationSet(`[
		book('Dracula')
		book('Frankenstein')
		book('Curse of the mummy')

		read('abraham', 'Dracula')
		read('abraham', 'Frankenstein')
		read('abraham', 'Curse of the mummy')
		read('sarah', 'Dracula')
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		book(A) :- book(A);
		read(A, B) :- read(A, B);
	]`)
	ds2dbWrite := internalGrammarParser.CreateRules(`[]`)

	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	entities := mentalese.Entities{}
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	factBase := knowledge.NewInMemoryFactBase("in-memory", facts, matcher, ds2db, ds2dbWrite, entities, log)
	solver.AddFactBase(factBase)
	nestedStructureBase := nested.NewSystemNestedStructureBase(solver, dialogContext, predicates, log)
	solver.AddNestedStructureBase(nestedStructureBase)
	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", log)
	solver.AddFunctionBase(systemFunctionBase)
	tokenizer := parse.NewTokenizer(log)
	nameResolver := central.NewNameResolver(solver, matcher, predicates, log, dialogContext)
	parser := earley.NewParser(grammar, nameResolver, predicates, log)

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
		trees := parser.Parse(words)
		if len(trees) == 0 {
			t.Errorf("Cannot parse: %s", test.input)
			continue
		}
		relationizer := earley.NewRelationizer(log)
		input, _ := relationizer.Relationize(trees[0], nameResolver)
		result := solver.SolveRelationSet(input, mentalese.Bindings{})
		if result.String() != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.input, result.String(), test.want)
			print(log.String())
		}
		log.Clear()
	}
}
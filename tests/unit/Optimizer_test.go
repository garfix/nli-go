package tests


import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"testing"
)

func TestOptimizer(t *testing.T) {

	log := common.NewSystemLog(false)
	parser := importer.NewInternalGrammarParser()

	facts1 := mentalese.RelationSet{}
	facts2 := mentalese.RelationSet{}

	ds2db1 := parser.CreateTransformations(`[
		married_to(A, B) => spouse(A, B);
		name(A, B) => name(A, B);
		first_name(A, F) last_name(A, L) join(N, ' ', F, L) => name(A, N);
		sibling(A, B) => parent(P, A) parent(P, B) unequal(A, B);
	]`)

	ds2db2 := parser.CreateTransformations(`[
		child(A, B) => parent(B, A);
		name(A, B) => relation(A, B);
		king(A) => rules(C, A);
		queen(A) => title(A, 'king');
		lord(A) => title(A, 'lord');
	]`)

	stats1 := mentalese.DbStats{
		"spouse": {Size: 100, DistinctValues: []int{75, 75}},
		"name": {Size: 200, DistinctValues: []int{200, 180}},
	}

	stats2 := mentalese.DbStats{
		"parent": {Size: 300, DistinctValues: []int{200, 200}},
		"title": {Size: 10, DistinctValues: []int{10}},
		"relation": {Size: 200, DistinctValues: []int{200, 180}},
	}

	matcher := mentalese.NewRelationMatcher(log)

	factBase1 := knowledge.NewInMemoryFactBase(facts1, matcher, ds2db1, stats1, log)
	factBase2 := knowledge.NewInMemoryFactBase(facts2, matcher, ds2db2, stats2, log)

	factBases := []knowledge.KnowledgeBase{factBase1, factBase2}

	systemFunctionBase := knowledge.NewSystemFunctionBase()
	matcher.AddFunctionBase(systemFunctionBase)
	optimizer := central.NewOptimizer(matcher)

	tests := []struct {
		input           string
		output 			string
		remaining       string
	}{
		// one predicate ('name') has stats and occurs in both fact bases; it is put up front
		{"[married_to(A, B) name(B, 'Lord Byron')]", "[[[name(B, 'Lord Byron')]@0, [married_to(A, B)]@0]\n [[name(B, 'Lord Byron')]@1, [married_to(A, B)]@0]]", "[]"},

		// a predicate ('who') that does not occur in any of the fact bases will appear in the remaining set
		{"[married_to(A, B) who(A) name(B, 'Lord Byron')]", "[]", "[who(A)]"},

		// a predicate ('king') that occurs in one fact base, but has no stats will appear last
		{"[married_to(A, B) king(A) name(B, 'Lord Byron')]", "[[[name(B, 'Lord Byron')]@0, [married_to(A, B)]@0, [king(A)]@1]\n [[name(B, 'Lord Byron')]@1, [married_to(A, B)]@0, [king(A)]@1]]", "[]"},

		// a predicate ('queen') with stats that occurs in one fact base, but not in the other will have a high position
		// in this example, all variables are unbound; the order is purely determined by the sizes of the tables
		{"[married_to(A, B) queen(A) name(A, B)]", "[[[queen(A)]@1, [married_to(A, B)]@0, [name(A, B)]@0]\n [[queen(A)]@1, [married_to(A, B)]@0, [name(A, B)]@1]]", "[]"},

		// match 2 predicates
		{"[first_name(C, 'Elvis') first_name(A, 'Lord') last_name(A, 'Byron')]", "[]", "[first_name(C, 'Elvis')]"},

		// 2 predicates that is more bound should precede 2 predicates that is less
		{"[first_name(A, 'Lord') last_name(A, 'Byron') first_name(11, 'Lord') last_name(11, 'Byron')]", "[[[first_name(11, 'Lord') last_name(11, 'Byron')]@0, [first_name(A, 'Lord') last_name(A, 'Byron')]@0]]", "[]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		routes, remainingRelations, _ := optimizer.CreateSolutionRoutes(input, factBases)
		outputString := routes.String()

		if outputString != test.output {
			t.Errorf("Optimizer: output, got %v, want %s", outputString, test.output)
		}
		if remainingRelations.String() != test.remaining {
			t.Errorf("Optimizer: remaining, got %v, want %s", remainingRelations.String(), test.remaining)
		}
	}
}

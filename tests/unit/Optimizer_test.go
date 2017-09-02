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

	ds2db1 := parser.CreateDbMappings(`[
		married_to(A, B) ->> spouse(A, B);
		name(A, B) ->> name(A, B);
	]`)

	ds2db2 := parser.CreateDbMappings(`[
		child(A, B) ->> parent(B, A);
		name(A, B) ->> relation(A, B);
		king(A) ->> title(A, 'king');
		queen(A) ->> title(A, 'king');
	]`)

	stats1 := mentalese.DbStats{
		"married_to": {Size: 100, DistinctValues: []int{75, 75}},
		"name": {Size: 200, DistinctValues: []int{200, 180}},
	}

	stats2 := mentalese.DbStats{
		"child": {Size: 300, DistinctValues: []int{200, 200}},
		"queen": {Size: 10, DistinctValues: []int{10}},
		"name": {Size: 200, DistinctValues: []int{200, 180}},
	}

	factBase1 := knowledge.NewInMemoryFactBase(facts1, ds2db1, stats1, log)
	factBase2 := knowledge.NewInMemoryFactBase(facts2, ds2db2, stats2, log)

	factBases := []knowledge.FactBase{factBase1, factBase2}

	optimizer := central.Optimizer{}

	tests := []struct {
		input           string
		output 			string
	}{
		// one predicate ('name') has stats and occurs in both fact bases; it is put up front
		{"[married_to(A, B) name(B, 'Lord Byron')]", "[name(B, 'Lord Byron') married_to(A, B)]"},

		// a predicate ('who') that does not occur in any of the fact bases will appear last
		// that way, it will profit from the bindings that are collected from the fact bases
		{"[married_to(A, B) who(A) name(B, 'Lord Byron')]", "[name(B, 'Lord Byron') married_to(A, B) who(A)]"},

		// none of the predicates occurs in any of the fact bases, no position change
		{"[abc(A, B) def(B, 'Lord Byron')]", "[abc(A, B) def(B, 'Lord Byron')]"},

		// a predicate ('king') that occurs in one fact base, but has no stats will appear last
		{"[married_to(A, B) king(A) name(B, 'Lord Byron')]", "[name(B, 'Lord Byron') married_to(A, B) king(A)]"},

		// a predicate ('queen') with stats that occurs in one fact base, but not in the other will have a high position
		// in this example, all variables are unbound; the order is purely determined by the sizes of the tables
		{"[married_to(A, B) queen(A) name(A, B)]", "[queen(A) married_to(A, B) name(A, B)]"},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		resultRelationSet := optimizer.Optimize(input, factBases)

		if resultRelationSet.String() != test.output {
			t.Errorf("Optimizer: got %v, want %s", resultRelationSet.String(), test.output)
		}
	}
}

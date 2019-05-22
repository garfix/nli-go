package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"testing"
)

func TestCost(t *testing.T) {

	log := common.NewSystemLog(false)
	parser := importer.NewInternalGrammarParser()

	stats1 := mentalese.DbStats{
		"name": {Size: 120000, DistinctValues: []int{100000, 120000}},
		"father": {Size: 50000, DistinctValues: []int{40000, 80000}},
		"country": {Size: 100, DistinctValues: []int{100, 90}},
	}

	matcher := mentalese.NewRelationMatcher(log)

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function")
	matcher.AddFunctionBase(systemFunctionBase)

	tests := []struct {
		input           string
		output 			float64
	}{
		{"[name(A, 'Byron')]", 1.0},
		{"[country(A, B)]", 100},
		{"[father('Ada', 'Byron')]", 0.000015625},
		{"[sibling('Ada', 'Byron')]", 100000000.0},
		{"[name(A, 'Byron') country(A, B) father('Ada', 'Byron')]", 101.000015625},
	}

	for _, test := range tests {

		input := parser.CreateRelationSet(test.input)

		cost := knowledge.CalculateCost(input, stats1)

		if cost != test.output {
			t.Errorf("Optimizer: cost, got %f, want %f", cost, test.output)
		}
	}
}

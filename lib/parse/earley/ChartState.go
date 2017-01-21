package earley

import (
	"nli-go/lib/parse"
	"nli-go/lib/mentalese"
)

type chartState struct {
	rule parse.GrammarRule
	dotPosition int
	startWordIndex int
	endWordIndex int
	sense mentalese.RelationSet
	children []int
	id int
}

func newChartState(rule parse.GrammarRule, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule: rule,
		dotPosition: dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex: endWordIndex,
		sense: mentalese.RelationSet{},
		children: []int{},
		id: 0,
	}
}
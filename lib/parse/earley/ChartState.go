package earley

import (
	"nli-go/lib/parse"
)

type chartState struct {
	rule parse.GrammarRule
	dotPosition int
	startWordIndex int
	endWordIndex int
	children []int
	id int
}

func newChartState(rule parse.GrammarRule, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule: rule,
		dotPosition: dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex: endWordIndex,
		children: []int{},
		id: 0,
	}
}
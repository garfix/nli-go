package earley

import (
	"nli-go/lib/parse"
)

type chartState struct {
	rule           parse.GrammarRule
	dotPosition    int
	startWordIndex int
	endWordIndex   int

	childStateIds []int
	id            int
}

func newChartState(rule parse.GrammarRule, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule:           rule,
		dotPosition:    dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex:   endWordIndex,

		childStateIds: []int{},
		id:            0,
	}
}

func (state chartState) isLeafState() bool {
	return len(state.childStateIds) == 0
}

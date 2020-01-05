package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/parse"
)

type chartState struct {
	rule           parse.GrammarRule
	sSelection	   parse.SSelection
	dotPosition    int
	startWordIndex int
	endWordIndex   int

	nameInformations []central.NameInformation

	childStateIds []int
	id            int
}

func newChartState(rule parse.GrammarRule, sSelection parse.SSelection, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule:           rule,
		sSelection:     sSelection,
		dotPosition:    dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex:   endWordIndex,

		nameInformations: []central.NameInformation{},

		childStateIds: []int{},
		id:            0,
	}
}

func (state chartState) isLeafState() bool {
	return len(state.childStateIds) == 0
}

func (state chartState) isIncomplete() bool {

	return state.dotPosition < state.rule.GetConsequentCount()+1
}

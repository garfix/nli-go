package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/parse"
	"strconv"
)

type chartState struct {
	rule           parse.GrammarRule
	sSelection	   parse.SSelection
	dotPosition    int
	startWordIndex int
	endWordIndex   int

	nameInformations []central.NameInformation

	id            int
	parentIds	  []int
}

func newChartState(id int, rule parse.GrammarRule, sSelection parse.SSelection, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule:           rule,
		sSelection:     sSelection,
		dotPosition:    dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex:   endWordIndex,

		nameInformations: []central.NameInformation{},

		parentIds: 		[]int{},
		id:            	id,
	}
}

func (state chartState) isIncomplete() bool {

	return state.dotPosition < state.rule.GetConsequentCount()+1
}

func (state chartState) Equals(otherState chartState) bool {
	return state.rule.Equals(otherState.rule) &&
		state.dotPosition == otherState.dotPosition &&
		state.startWordIndex == otherState.startWordIndex &&
		state.endWordIndex == otherState.endWordIndex &&
		common.IntArrayEquals(state.parentIds, otherState.parentIds)
}

func (state chartState) ToString(chart *chart) string {
	s := strconv.Itoa(state.id) + " ["
	for i, category := range state.rule.SyntacticCategories {
		if i == 0 {
			s += " " + category + " ->"
		} else {
			if i == state.dotPosition {
				s += " *"
			}
			s += " " + category
		}
	}
	if len(state.rule.SyntacticCategories) == state.dotPosition {
		s += " *"
	}
	s += " ] "

	s += "<"
	for i, word := range chart.words {
		if i >= state.startWordIndex && i < state.endWordIndex {
			s += " " + word
		}
	}
	s += " >"
	s += " ("
	for _, parentId := range state.parentIds {
		s += " " + strconv.Itoa(parentId)
	}
	s += " )"
	return s
}
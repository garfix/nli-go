package parse

import (
	"nli-go/lib/mentalese"
	"strconv"
)

type chartState struct {
	rule           mentalese.GrammarRule
	dotPosition    int
	startWordIndex int
	endWordIndex   int
}

func newChartState(rule mentalese.GrammarRule, dotPosition int, startWordIndex int, endWordIndex int) chartState {
	return chartState{
		rule:           rule,
		dotPosition:    dotPosition,
		startWordIndex: startWordIndex,
		endWordIndex:   endWordIndex,
	}
}

func (state chartState) isTerminal() bool {
	if len(state.rule.GetConsequentVariables(0)) == 0 {
		return false
	}
	return state.rule.GetConsequentVariables(0)[0] == terminal
}

func (state chartState) isComplete() bool {

	return state.dotPosition >= state.rule.GetConsequentCount()+1
}

func (state chartState) Equals(otherState chartState) bool {
	return state.rule.Equals(otherState.rule) &&
		state.dotPosition == otherState.dotPosition &&
		state.startWordIndex == otherState.startWordIndex &&
		state.endWordIndex == otherState.endWordIndex
}

func (state chartState) BasicForm() string {
	s := state.rule.BasicForm()
	s += " [" + strconv.Itoa(state.startWordIndex) + "-" + strconv.Itoa(state.endWordIndex) + "]"
	return s
}

func (state chartState) StartForm() string {
	s := state.rule.BasicForm()
	s += " " + strconv.Itoa(state.startWordIndex)
	return s
}

func (state chartState) ToString(chart *chart) string {
	s := state.rule.GetAntecedent() + " ->"
	for i, category := range state.rule.GetConsequents() {
		if i + 1 == state.dotPosition {
			s += " *"
		}
		s += " " + category
	}
	if len(state.rule.GetConsequents()) + 1 == state.dotPosition {
		s += " *"
	}

	s += " <"
	for i, word := range chart.words {
		if i >= state.startWordIndex && i < state.endWordIndex {
			s += " " + word
		}
	}
	s += " >"
	return s
}
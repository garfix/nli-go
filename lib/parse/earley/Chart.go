package earley

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

// Contains more than the strict chart that the Earley algorithm prescribes; it is used to hold all state of a parse.

type chart struct {
	states [][]chartState
	words  []string
	stateIdGenerator int
	children map[string][][]chartState
}

func newChart(words []string) *chart {
	return &chart{
		states:           make([][]chartState, len(words)+1),
		words:            words,
		stateIdGenerator: 0,
		children: map[string][][]chartState{},
	}
}

func (chart *chart) buildIncompleteGammaState() chartState {
	return newChartState(
		chart.generateId(),
		parse.NewGrammarRule(
			[]string{ parse.PosTypeRelation, parse.PosTypeRelation },
			[]string{"gamma", "s"},
			[][]string{{"G"}, {"S"}},
			mentalese.RelationSet{},
		),
		[][]string{{""}, {""}},
		1, 0, 0)
}

func (chart *chart) buildCompleteGammaState() chartState {
	state := chart.buildIncompleteGammaState()
	state.dotPosition = 2
	state.endWordIndex = len(chart.words)
	return state
}

func (chart *chart) generateId() int {
	chart.stateIdGenerator++
	return chart.stateIdGenerator
}

func (chart *chart) indexChildren(state chartState) {

	canonical := state.BasicForm()

	_, found := chart.children[canonical]
	if !found {
		chart.children[canonical] = [][]chartState{}
	}

	chart.children[canonical] = append(chart.children[canonical], state.children)
}

func (chart *chart) enqueue(state chartState, position int) bool {

	found := chart.isStateInChart(state, position)
	if !found {
		chart.pushState(state, position)
	}

	return found
}

func (chart *chart) isStateInChart(state chartState, position int) bool {

	for _, presentState := range chart.states[position] {
		if presentState.Equals(state) {
			return true
		}
	}

	return false
}

func (chart *chart) pushState(state chartState, position int) {

	chart.states[position] = append(chart.states[position], state)
}

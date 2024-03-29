package parse

import (
	"nli-go/lib/mentalese"
)

type chart struct {
	rootCategory  string
	rootVariables []string
	words         []string
	states        [][]chartState
	advanced      map[string][][]chartState
	completed     map[string][][]chartState
}

func NewChart(words []string, rootCategory string, rootVariables []string) *chart {
	return &chart{
		rootCategory:  rootCategory,
		rootVariables: rootVariables,
		words:         words,
		states:        make([][]chartState, len(words)+1),
		advanced:      map[string][][]chartState{},
		completed:     map[string][][]chartState{},
	}
}

func (chart *chart) buildIncompleteGammaState() chartState {
	return newChartState(
		mentalese.NewGrammarRule(
			[]string{mentalese.PosTypeRelation, mentalese.PosTypeRelation},
			[]string{"gamma", chart.rootCategory},
			[][]string{{"G"}, chart.rootVariables},
			mentalese.RelationSet{},
		),
		1, 0, 0)
}

func (chart *chart) buildCompleteGammaState() chartState {
	state := chart.buildIncompleteGammaState()
	state.dotPosition = 2
	state.endWordIndex = len(chart.words)
	return state
}

func (chart *chart) updateAdvancedStatesIndex(completedState chartState, advancedState chartState) {

	canonical := advancedState.StartForm()
	completedConsequentsCount := advancedState.dotPosition - 2

	_, found := chart.advanced[canonical]
	if !found {
		chart.advanced[canonical] = [][]chartState{}
	}

	children := []chartState{}
	if completedConsequentsCount == 0 {
		children = []chartState{completedState}
		chart.addAdvancedStateIndex(advancedState, children)
	} else {
		for _, previousChildren := range chart.advanced[canonical] {
			if len(previousChildren) == completedConsequentsCount {
				if previousChildren[len(previousChildren)-1].endWordIndex == completedState.startWordIndex {
					children = chart.appendState(previousChildren, completedState)
					chart.addAdvancedStateIndex(advancedState, children)
				}
			}
		}
	}
}

func (chart *chart) addAdvancedStateIndex(advancedState chartState, children []chartState) {
	canonical := advancedState.StartForm()
	chart.advanced[canonical] = append(chart.advanced[canonical], children)

	if advancedState.isComplete() {
		chart.updateCompletedStatesIndex(advancedState, children)
	}
}

func (chart *chart) updateCompletedStatesIndex(advancedState chartState, children []chartState) {

	canonical := advancedState.BasicForm()

	_, found := chart.completed[canonical]
	if !found {
		chart.completed[canonical] = [][]chartState{}
	}

	chart.completed[canonical] = append(chart.completed[canonical], children)
}

func (chart *chart) appendState(oldStates []chartState, newState chartState) []chartState {
	newStates := []chartState{}
	for _, state := range oldStates {
		newStates = append(newStates, state)
	}
	newStates = append(newStates, newState)
	return newStates
}

func (chart *chart) enqueue(state chartState, position int) bool {

	found := chart.containsState(state, position)
	if !found {
		chart.pushState(state, position)
	}

	return found
}

func (chart *chart) containsState(state chartState, position int) bool {

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

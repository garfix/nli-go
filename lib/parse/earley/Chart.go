package earley

// Contains more than the strict chart that the Earley algorithm prescribes; it is used to hold all state of a parse.

type chart struct {
	states [][]chartState
	words  []string

	sentenceStates   []chartState
	indexedStates    map[int]chartState
	stateIdGenerator int
}

func newChart(words []string) *chart {
	return &chart{
		states:           make([][]chartState, len(words)+1),
		words:            words,
		sentenceStates:   []chartState{},
		indexedStates:    map[int]chartState{},
		stateIdGenerator: 0,
	}
}

func (chart *chart) enqueue(state chartState, position int) {

	if !chart.isStateInChart(state, position) {
		chart.pushState(state, position)
	}
}

func (chart *chart) isStateInChart(state chartState, position int) bool {

	for _, presentState := range chart.states[position] {

		if presentState.rule.Equals(state.rule) &&
			presentState.dotPosition == state.dotPosition &&
			presentState.startWordIndex == state.startWordIndex &&
			presentState.endWordIndex == state.endWordIndex {

			return true
		}
	}

	return false
}

func (chart *chart) pushState(state chartState, position int) {

	// index the state for later lookup
	chart.stateIdGenerator++
	state.id = chart.stateIdGenerator
	chart.indexedStates[state.id] = state

	chart.states[position] = append(chart.states[position], state)
}

func (chart *chart) storeStateInfo(completedState chartState, chartedState chartState, advancedState chartState) (bool, chartState) {

	treeComplete := false

	// store the state's "children" to ease building the parse trees from the packed forest
	advancedState.childStateIds = append(chartedState.childStateIds, completedState.id)

	// rule complete?
	if chartedState.dotPosition == chartedState.rule.GetConsequentCount() {

		// complete sentence?
		if chartedState.rule.GetAntecedent() == "gamma" {

			// that matches all words?
			if completedState.endWordIndex == len(chart.words) {

				chart.sentenceStates = append(chart.sentenceStates, advancedState)

				// set a flag to allow the Parser to stop at the first complete parse
				treeComplete = true
			}
		}
	}

	return treeComplete, advancedState
}
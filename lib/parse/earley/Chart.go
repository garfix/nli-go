package earley

// Contains more than the strict chart that the Earley algorithm prescribes; it is used to hold all state of a parse.

type chart struct {
	states [][]chartState
	words  []string
	stateIdGenerator int
	children map[string][][]int
}

func newChart(words []string) *chart {
	return &chart{
		states:           make([][]chartState, len(words)+1),
		words:            words,
		stateIdGenerator: 0,
		children: map[string][][]int{},
	}
}

func (chart *chart) generateId() int {
	chart.stateIdGenerator++
	return chart.stateIdGenerator
}

func (chart *chart) indexChildren(state chartState) {

	canonical := state.Canonical()

	_, found := chart.children[canonical]
	if !found {
		chart.children[canonical] = [][]int{}
	}

	chart.children[canonical] = append(chart.children[canonical], state.parentIds)
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

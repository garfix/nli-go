package earley

type chart struct {
	states           [][]chartState
	words            []string

	sentenceStates   []chartState
	indexedStates    map[int]chartState
	stateIdGenerator int
}

func newChart(words []string) *chart {
	return &chart{
		states: make([][]chartState, len(words) + 1),
		words: words,
		sentenceStates: []chartState{},
		indexedStates: map[int]chartState{},
		stateIdGenerator: 0,
	}
}
package earley

type chart struct {
	states [][]chartState
	words []string

	treeInfoSentences []chartState
	treeInfoStates map[int]chartState
	stateIdGenerator int
}

func newChart(words []string) *chart {
	return &chart{
		states: make([][]chartState, len(words) + 1),
		words: words,
		treeInfoSentences: []chartState{},
		treeInfoStates: map[int]chartState{},
		stateIdGenerator: 0,
	}
}
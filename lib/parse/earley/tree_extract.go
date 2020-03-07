package earley


func extractTreeRoots(chart *chart) []ParseTreeNode {

	rootStates := []chartState{}
	stateIndex := map[int]chartState{}
	rootNodes := []ParseTreeNode{}

	for _, states := range chart.states {
		for _, state := range states {

			stateIndex[state.id] = state

			if state.rule.GetAntecedent() == "gamma" && !state.isIncomplete() && state.endWordIndex == len(chart.words) {
				sentenceState := stateIndex[state.parentIds[0]]
				rootStates = append(rootStates, sentenceState)
			}
		}
	}

	for _, rootState := range rootStates {
		rootNodes = append(rootNodes, extractTreesForState(chart, rootState, &stateIndex))
	}

	return rootNodes
}

func extractTreesForState(chart *chart, state chartState, stateIndex *map[int]chartState) ParseTreeNode {

	rule := state.rule
	branch := ParseTreeNode{category: rule.GetAntecedent(), constituents: []ParseTreeNode{}, form: "", rule: state.rule, nameInformations: state.nameInformations}

	if len(state.parentIds) == 0 {

		branch.form = rule.GetConsequent(0)

	} else {

		for _, childStateId := range state.parentIds {

			childState := (*stateIndex)[childStateId]
			branch.constituents = append(branch.constituents, extractTreesForState(chart, childState, stateIndex))
		}
	}

	return branch
}

// Returns the word that could not be parsed (or ""), and the index of the last completed word
func findLastCompletedWordIndex(chart *chart) (int, string) {

	nextWord := ""
	lastIndex := -1

	// find the last completed nextWord

	for i := len(chart.states) - 1; i >= 0; i-- {
		states := chart.states[i]
		for _, state := range states {
			if !state.isIncomplete() {

				lastIndex = state.endWordIndex - 1
				goto done
			}
		}
	}

done:

	if lastIndex <= len(chart.words)-2 {
		nextWord = chart.words[lastIndex+1]
	}

	return lastIndex, nextWord
}


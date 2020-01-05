package earley


func extractFirstTree(chart *chart) ParseTreeNode {

	tree := ParseTreeNode{}

	if len(chart.sentenceStates) > 0 {

		rootStateId := chart.sentenceStates[0].childStateIds[0]
		root := chart.indexedStates[rootStateId]
		tree = extractParseTreeBranch(chart, root)
	}

	return tree
}

func extractParseTreeBranch(chart *chart, state chartState) ParseTreeNode {

	rule := state.rule
	branch := ParseTreeNode{category: rule.GetAntecedent(), constituents: []ParseTreeNode{}, form: "", rule: state.rule, nameInformations: state.nameInformations}

	if state.isLeafState() {

		branch.form = rule.GetConsequent(0)

	} else {

		for _, childStateId := range state.childStateIds {

			childState := chart.indexedStates[childStateId]
			branch.constituents = append(branch.constituents, extractParseTreeBranch(chart, childState))
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


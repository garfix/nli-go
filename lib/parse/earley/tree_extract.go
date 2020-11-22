package earley

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

// no backtracking! uses custom stacks

type treeExtracter struct {
	trees []*ParseTreeNode
	stateIndex map[int]chartState
	chart *chart
}

func extractTreeRoots(chart *chart) []ParseTreeNode {

	extracter := &treeExtracter{
		trees: []*ParseTreeNode{},
		stateIndex: map[int]chartState{},
		chart: chart,
	}

	for _, states := range chart.states {
		for _, state := range states {
			extracter.stateIndex[state.id] = state
		}
	}

	extracter.extract()

	roots := []ParseTreeNode{}
	for _, root := range extracter.trees {
		if len(root.constituents) == 0 {
			continue
		}
		roots = append(roots, *root.constituents[0])
	}

	return roots
}

func (ex *treeExtracter) extract() {

	wordCount := len(ex.chart.words)
	rule := parse.NewGrammarRule([]string{ parse.PosTypeRelation, parse.PosTypeRelation }, []string{"gamma", "s"}, [][]string{{"G"}, {"S"}}, mentalese.RelationSet{})
	completedGammaState := newChartState(0, rule, [][]string{{""}, {""}}, 2, 0, wordCount)
	rootNode := &ParseTreeNode{
		category:         "gamma",
		constituents:     []*ParseTreeNode{},
		form:             "",
		rule:             rule,
		nameInformations: nil,
	}

	ex.trees = append(ex.trees, rootNode)

	tree := treeInProgress{
		root: rootNode,
		path: []workingStep{
		{
			states:      []chartState{completedGammaState},
			nodes:       []*ParseTreeNode{ rootNode },
			stateIndex: 0,
		},
	}}

	ex.next(tree)
}

// walk through the parse-tree-in-progress, one step at a time
func (ex *treeExtracter) next(tree treeInProgress) {

	newTree, done := tree.advance()
	if done {
		return
	}

	ex.forkOnCompletedStates(newTree)
}

func (ex *treeExtracter) forkOnCompletedStates(tree treeInProgress) {

	state := tree.peek().getCurrentState()

	//variants := ex.findCompletedStates(state)
	//
	//newTrees := ex.forkTrees(tree, len(variants))
	//
	//for i, variant := range variants {
	//	newTree := newTrees[i]
	//	ex.addChildren(newTree, variant)
	//}

	ex.addChildren(tree, state)
}

func (ex *treeExtracter) addChildren(tree treeInProgress, parentState chartState) {

	allChildStates := ex.findCompletedChildStates(parentState)

	if len(allChildStates) == 0 {

		ex.next(tree)

	} else {

		newTrees := ex.forkTrees(tree, len(allChildStates))

		for i, childStates := range allChildStates {

			newTree := newTrees[i]

			childNodes := []*ParseTreeNode{}
			parentNode := newTree.peek().getCurrentNode()

			for _, childState := range childStates {
				childNode := ex.createNode(childState)
				childNodes = append(childNodes, childNode)
				parentNode.constituents = append(parentNode.constituents, childNode)
			}

			step := workingStep{
				states:     childStates,
				nodes:      childNodes,
				stateIndex: 0,
			}

			newTree = newTree.push(step)

			ex.next(newTree)
		}
	}
}

// create `count` clones of `tree`; the first tree is just the original
func (ex *treeExtracter) forkTrees(tree treeInProgress, count int) []treeInProgress {

	tips := []treeInProgress{}

	for i := 0; i < count; i++ {
		if i == 0 {
			tips = append(tips, tree)
		} else {
			newTip := tree.clone()
			tips = append(tips, newTip)
			ex.trees = append(ex.trees, newTip.root)
		}
	}

	return tips
}

func (ex *treeExtracter) createNode(state chartState) *ParseTreeNode {

	node := &ParseTreeNode{
		category: state.rule.GetAntecedent(),
		constituents: []*ParseTreeNode{},
		form: "",
		rule: state.rule,
		nameInformations: state.nameInformations,
	}

	if len(state.parentIds) == 0 {
		node.form = state.rule.GetConsequent(0)
	}

	return node
}

// Find completed versions of `state`
func (ex *treeExtracter) findCompletedStates(state chartState) []chartState {

	startWordIndex := state.startWordIndex
	endWordIndex := state.endWordIndex

	completedStates := []chartState{}

	for _, aState := range ex.chart.states[endWordIndex] {
		if !aState.isIncomplete() &&
			aState.rule.Equals(state.rule) &&
			aState.startWordIndex == startWordIndex &&
			aState.endWordIndex == endWordIndex {

			completedStates = append(completedStates, aState)
		}
	}

	return completedStates
}

func (ex *treeExtracter) findCompletedChildStates(state chartState) [][]chartState {

	allChildStates := [][]chartState{}

	rows, found := ex.chart.children[state.Canonical()]
	if found {

		for _, row := range rows {
			children := []chartState{}
			for _, stateId := range row {
				state1 := ex.stateIndex[stateId]
				children = append(children, state1)
			}
			allChildStates = append(allChildStates, children)
		}

	}

	return allChildStates
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


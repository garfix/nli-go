package parse

// no backtracking! uses custom stacks

type treeExtracter struct {
	trees []*ParseTreeNode
	chart *chart
}

func ExtractTreeRoots(chart *chart) []ParseTreeNode {

	extracter := &treeExtracter{
		trees: []*ParseTreeNode{},
		chart: chart,
	}

	extracter.extract()

	// the sentence node is the first child
	roots := []ParseTreeNode{}
	for _, root := range extracter.trees {
		if len(root.Constituents) > 0 {
			roots = append(roots, *root.Constituents[0])
		}
	}

	return roots
}

func (ex *treeExtracter) extract() {

	completedGammaState := ex.chart.buildCompleteGammaState()

	rootNode := &ParseTreeNode{
		Category:     "gamma",
		Constituents: nil,
		Form:         "",
		Rule:         completedGammaState.rule,
	}

	ex.trees = append(ex.trees, rootNode)

	tree := treeInProgress{
		root: rootNode,
		path: []workingStep{
		{
			states:      []chartState{completedGammaState},
			nodes:       []*ParseTreeNode{rootNode },
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

	ex.addChildren(newTree)
}

func (ex *treeExtracter) addChildren(tree treeInProgress) {

	parentState := tree.peek().getCurrentState()

	allChildStates, found := ex.chart.completed[parentState.BasicForm()]
	if !found {

		ex.next(tree)

	} else {

		newTrees := ex.forkTrees(tree, len(allChildStates))

		for i, childStates := range allChildStates {

			newTree := newTrees[i]
			parentNode := newTree.peek().getCurrentNode()

			childNodes := []*ParseTreeNode{}
			for _, childState := range childStates {
				childNodes = append(childNodes, ex.createNode(childState))
			}
			parentNode.Constituents = childNodes

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
// the new trees are registered with the tree extractor
func (ex *treeExtracter) forkTrees(tree treeInProgress, count int) []treeInProgress {

	tips := []treeInProgress{}

	for i := 0; i < count; i++ {
		if i == 0 {
			tips = append(tips, tree)
		} else {
			newTip := tree.clone()
			tips = append(tips, newTip)

			// register new tree
			ex.trees = append(ex.trees, newTip.root)
		}
	}

	return tips
}

// creates a single parse tree node
func (ex *treeExtracter) createNode(state chartState) *ParseTreeNode {

	form := ""
	if state.isTerminal() {
		form = state.rule.GetConsequent(0)
	}

	return &ParseTreeNode{
		Category:     state.rule.GetAntecedent(),
		Constituents: []*ParseTreeNode{},
		Form:         form,
		Rule:         state.rule,
	}
}

// Returns the word that could not be parsed (or ""), and the index of the last completed word
func FindLastCompletedWordIndex(chart *chart) (int, string) {

	nextWord := ""
	lastIndex := -1

	// find the last completed nextWord

	for i := len(chart.states) - 1; i >= 0; i-- {
		states := chart.states[i]
		for _, state := range states {
			if state.isComplete() {

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


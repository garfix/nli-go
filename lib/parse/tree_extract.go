package parse

import (
	"nli-go/lib/mentalese"
	"strings"
)

// no backtracking! uses custom stacks

type treeExtracter struct {
	trees []*mentalese.ParseTreeNode
	chart *chart
}

func ExtractTreeRoots(chart *chart) []*mentalese.ParseTreeNode {

	extracter := &treeExtracter{
		trees: []*mentalese.ParseTreeNode{},
		chart: chart,
	}

	extracter.extract()

	// the sentence node is the first child
	roots := []*mentalese.ParseTreeNode{}
	for _, root := range extracter.trees {
		if len(root.Constituents) > 0 {
			roots = append(roots, root.Constituents[0])
		}
	}

	return roots
}

func (ex *treeExtracter) extract() {

	completedGammaState := ex.chart.buildCompleteGammaState()

	rootNode := &mentalese.ParseTreeNode{
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
				states:     []chartState{completedGammaState},
				nodes:      []*mentalese.ParseTreeNode{rootNode},
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

			childNodes := []*mentalese.ParseTreeNode{}
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
func (ex *treeExtracter) createNode(state chartState) *mentalese.ParseTreeNode {

	form := ""
	if state.isTerminal() {
		form = state.rule.GetConsequent(0)
	}

	return &mentalese.ParseTreeNode{
		Category:     state.rule.GetAntecedent(),
		Constituents: []*mentalese.ParseTreeNode{},
		Form:         form,
		Rule:         state.rule,
	}
}

// Returns the word that could not be parsed (or ""), and the index of the last completed word
func FindUnknownWord(chart *chart) string {

	nextWord := ""
	lastUnderstoodIndex := -1

	// find the last completed nextWord

	for i := len(chart.states) - 1; i >= 0; i-- {
		states := chart.states[i]
		for _, state := range states {
			if state.isComplete() {
				if state.endWordIndex > lastUnderstoodIndex {
					lastUnderstoodIndex = state.endWordIndex - 1
				}
			}
		}
	}

	if lastUnderstoodIndex+1 < len(chart.words) {
		nextWord = chart.words[lastUnderstoodIndex+1]
	} else {
		nextWord = strings.Join(chart.words, " ")
	}

	return nextWord
}

package parse


type workingStep struct {
	states     	[]chartState
	nodes       []*ParseTreeNode
	stateIndex 	int
}

func (step workingStep) getCurrentState() chartState {
	return step.states[step.stateIndex - 1]
}

func (step workingStep) getCurrentNode() *ParseTreeNode {
	return step.nodes[step.stateIndex - 1]
}

type treeInProgress struct {
	root *ParseTreeNode
	path []workingStep
}

func (tip treeInProgress) clone() treeInProgress {

	newRoot, aMap := tip.cloneTree(tip.root)

	newSteps := []workingStep{}
	for _, step := range tip.path {
		newNodes := []*ParseTreeNode{}
		for _, node := range step.nodes {
			newNode, _ := aMap[node]
			newNodes = append(newNodes, newNode)
		}
		newStep := workingStep{
			states:     step.states,
			nodes:      newNodes,
			stateIndex: step.stateIndex,
		}
		newSteps = append(newSteps, newStep)
	}

	newStack := treeInProgress{
		root: newRoot,
		path: newSteps,
	}

	return newStack
}


func (tip *treeInProgress) cloneTree(tree *ParseTreeNode) (*ParseTreeNode, map[*ParseTreeNode]*ParseTreeNode) {

	aMap := map[*ParseTreeNode]*ParseTreeNode{}
	newTree := tip.cloneNodeWithMap(tree, &aMap)

	return newTree, aMap
}

func (tip *treeInProgress) cloneNodeWithMap(node *ParseTreeNode, aMap *map[*ParseTreeNode]*ParseTreeNode) *ParseTreeNode {

	constituents := []*ParseTreeNode{}
	for _, constituent := range node.Constituents {
		clone := tip.cloneNodeWithMap(constituent, aMap)
		constituents = append(constituents, clone)
	}

	newNode := &ParseTreeNode{
		Category:     node.Category,
		Constituents: constituents,
		Form:         node.Form,
		Rule:         node.Rule,
	}

	(*aMap)[node] = newNode

	return newNode
}

func (tip treeInProgress) advance() (treeInProgress, bool) {

	newTip := tip
	done := true

	if len(newTip.path) > 0 {
		step := &newTip.path[len(newTip.path)-1]
		if step.stateIndex < len(step.states) {
			step.stateIndex++
		} else {
			return newTip.pop().advance()
		}
		done = false
	}

	return newTip, done
}

func (tip treeInProgress) peek() workingStep {
	if len(tip.path) == 0 {
		panic("empty stack!")
	} else {
		return tip.path[len(tip.path) - 1]
	}
}

func (tip treeInProgress) push(step workingStep) treeInProgress {
	newStack := tip
	newStack.path = append(newStack.path, step)
	return newStack
}

func (tip treeInProgress) pop() treeInProgress {
	newStack := tip
	newStack.path = newStack.path[0:len(newStack.path) - 1]
	return newStack
}

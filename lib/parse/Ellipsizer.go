package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type Ellipsizer struct {
	sentences []*mentalese.ParseTreeNode
	log       *common.SystemLog
}

func NewEllipsizer(Sentences []*mentalese.ParseTreeNode, log *common.SystemLog) *Ellipsizer {
	return &Ellipsizer{
		sentences: Sentences,
		log:       log,
	}
}

// Returns a copy of `tree` where the `ellipsis` directions are included
func (e *Ellipsizer) Ellipsize(tree *mentalese.ParseTreeNode) (*mentalese.ParseTreeNode, bool) {

	// quick check if this is necessary
	if !e.hasEllipsis(tree) {
		return tree, true
	}

	biDirTree := CreateBidirectionalParseTree(tree)

	variableMapping := &map[string]string{}
	newTree, ok := e.ellipsizeNode(biDirTree, variableMapping)
	updatedTree := e.replaceVariables(newTree, variableMapping)

	return updatedTree, ok
}

func (e *Ellipsizer) replaceVariables(node *mentalese.ParseTreeNode, variableMapping *map[string]string) *mentalese.ParseTreeNode {

	newNode := *node

	for fromVar, toVar := range *variableMapping {
		newNode = newNode.ReplaceVariable(fromVar, toVar)
	}

	return &newNode
}

// quick check
func (e *Ellipsizer) hasEllipsis(node *mentalese.ParseTreeNode) bool {
	present := len(node.Rule.Ellipsis) > 0
	for _, constituent := range node.Constituents {
		present = present || e.hasEllipsis(constituent)
	}
	return present
}

// Handle ellipsis in a single node
func (e *Ellipsizer) ellipsizeNode(node *BidirectionalParseTreeNode, variableMapping *map[string]string) (*mentalese.ParseTreeNode, bool) {
	ok := true

	// add original constituents, having been ellipsized
	newConstituents := []*mentalese.ParseTreeNode{}
	for _, child := range node.children {
		newConstituent, success := e.ellipsizeNode(child, variableMapping)
		if !success {
			ok = false
		}
		newConstituents = append(newConstituents, newConstituent)
	}

	// add new constituents, from ellipsis
	ellipsisConstituents := e.createEllipsisConstituents(node, variableMapping)
	if len(node.source.Rule.Ellipsis) != len(ellipsisConstituents) {
		ok = false
	}

	newSource := node.source.PartialCopy()
	newSource.Constituents = newConstituents
	for i, ellipsisConstituent := range ellipsisConstituents {
		categoryPath := node.source.Rule.Ellipsis[i]
		lastNode := categoryPath[len(categoryPath)-1]
		newSource.Constituents = append(newSource.Constituents, ellipsisConstituent)
		newSource.Rule.PositionTypes = append(newSource.Rule.PositionTypes, mentalese.PosTypeRelation)
		newSource.Rule.SyntacticCategories = append(newSource.Rule.SyntacticCategories, lastNode.Category)
		newSource.Rule.EntityVariables = append(newSource.Rule.EntityVariables, lastNode.Variables)
	}

	return &newSource, ok
}

// Create syntactic / semantic constituents to replace the ellipsis
func (e *Ellipsizer) createEllipsisConstituents(node *BidirectionalParseTreeNode, variableMapping *map[string]string) []*mentalese.ParseTreeNode {
	ellipsisConstituents := []*mentalese.ParseTreeNode{}

	source := node.source
	for _, categoryPath := range source.Rule.Ellipsis {
		newConstituent := e.processCategoryPath(node, categoryPath)
		if newConstituent != nil {
			ellipsisConstituents = append(ellipsisConstituents, newConstituent)
			e.mapVariables(categoryPath[len(categoryPath)-1].Variables, newConstituent.Rule.GetAntecedentVariables(), variableMapping)
		}
	}

	return ellipsisConstituents
}

func (e *Ellipsizer) mapVariables(pathNodeVariables []string, antecedentVariables []string, variableMapping *map[string]string) {
	for i, pathNodeVariable := range pathNodeVariables {
		antecedentVariable := antecedentVariables[i]
		(*variableMapping)[pathNodeVariable] = antecedentVariable
	}
}

// Handle a single path
func (e *Ellipsizer) processCategoryPath(currentNode *BidirectionalParseTreeNode, path mentalese.CategoryPath) *mentalese.ParseTreeNode {
	antecedentNode := e.step(currentNode, path)

	// Create an error if no node could be found (or a clarification request?)
	if antecedentNode == nil {
		e.log.AddError("No antecedent found for ellipsis")
	}

	return antecedentNode
}

// Move a single step on the ellipsis category path
func (e *Ellipsizer) step(currentNode *BidirectionalParseTreeNode, path mentalese.CategoryPath) *mentalese.ParseTreeNode {
	if len(path) == 0 {
		return currentNode.source
	}

	pathNode := path[0]
	restPath := path[1:]
	nextNodes := e.navigate(currentNode, pathNode)

	for _, nextNode := range nextNodes {
		result := e.step(nextNode, restPath)
		// a single result is enough
		if result != nil {
			return result
		}
	}

	return nil
}

// Find all possible next nodes given a single node in the ellipsis path
func (e *Ellipsizer) navigate(currentNode *BidirectionalParseTreeNode, direction mentalese.CategoryPathNode) []*BidirectionalParseTreeNode {

	newNodes := []*BidirectionalParseTreeNode{}

	switch direction.GetNodeType() {
	case mentalese.NodeTypePrevSentence:
		newNodes = e.navigatePrevSentence(currentNode)
	case mentalese.NodeTypeParent:
		newNodes = e.navigateParent(currentNode, direction.GetCategory())
	case mentalese.NodeTypeNextSibling:
		newNodes = e.navigateNextSibling(currentNode, direction.GetCategory())
	case mentalese.NodeTypePrevSibling:
		newNodes = e.navigatePrevSibling(currentNode, direction.GetCategory())
	case mentalese.NodeTypeSibling:
		newNodes = e.navigateSibling(currentNode, direction.GetCategory())
	case mentalese.NodeTypeChild:
		newNodes = e.navigateChild(currentNode, direction.GetCategory(), direction.DoesAllowIndirect())
	default:
		e.log.AddError("Node type not recognized: " + direction.GetNodeType())
	}

	return newNodes
}

// form: ..
// meaning: move up a node
func (e *Ellipsizer) navigateParent(currentNode *BidirectionalParseTreeNode, category string) []*BidirectionalParseTreeNode {

	currentNode = currentNode.parent

	if currentNode == nil {
		return []*BidirectionalParseTreeNode{}
	}

	if category != "" {
		if currentNode.source.Category != category {
			return e.navigateParent(currentNode.parent, category)
		}
	}

	return []*BidirectionalParseTreeNode{currentNode}
}

func (e *Ellipsizer) navigatePrevSentence(currentNode *BidirectionalParseTreeNode) []*BidirectionalParseTreeNode {
	var newNode *BidirectionalParseTreeNode = nil
	if len(e.sentences) > 1 {
		sentence := e.sentences[len(e.sentences)-2]
		// todo: this `prev` always goes to the last sentence
		newNode = CreateBidirectionalParseTree(sentence)
	} else {
		return []*BidirectionalParseTreeNode{}
	}
	return []*BidirectionalParseTreeNode{newNode}
}

func (e *Ellipsizer) navigatePrevSibling(currentNode *BidirectionalParseTreeNode, category string) []*BidirectionalParseTreeNode {
	newNodes := []*BidirectionalParseTreeNode{}

	parent := currentNode.parent
	if parent == nil {
		return []*BidirectionalParseTreeNode{}
	}

	active := false
	for i := len(parent.children) - 1; i >= 0; i-- {
		child := parent.children[i]
		if child == currentNode {
			active = true
		} else if active {
			if category != "" && child.source.Category != category {
				continue
			}
			newNodes = append(newNodes, child)
		}
	}

	return newNodes
}

func (e *Ellipsizer) navigateNextSibling(currentNode *BidirectionalParseTreeNode, category string) []*BidirectionalParseTreeNode {
	newNodes := []*BidirectionalParseTreeNode{}

	parent := currentNode.parent
	if parent == nil {
		return []*BidirectionalParseTreeNode{}
	}

	active := false
	for i := 0; i < len(parent.children); i++ {
		child := parent.children[i]
		if child == currentNode {
			active = true
		} else if active {
			if category != "" && child.source.Category != category {
				continue
			}
			newNodes = append(newNodes, child)
		}
	}

	return newNodes
}

func (e *Ellipsizer) navigateSibling(currentNode *BidirectionalParseTreeNode, category string) []*BidirectionalParseTreeNode {
	newNodes := []*BidirectionalParseTreeNode{}

	parent := currentNode.parent
	if parent == nil {
		return []*BidirectionalParseTreeNode{}
	}

	for i := 0; i < len(parent.children); i++ {
		child := parent.children[i]
		if child == currentNode {
			continue
		}
		if category != "" && child.source.Category != category {
			continue
		}
		newNodes = append(newNodes, child)
	}

	return newNodes
}

func (e *Ellipsizer) navigateChild(currentNode *BidirectionalParseTreeNode, category string, allowIndirect bool) []*BidirectionalParseTreeNode {
	newNodes := []*BidirectionalParseTreeNode{}

	for i := 0; i < len(currentNode.children); i++ {
		child := currentNode.children[i]
		if child == currentNode {
			continue
		}
		if category != "" && child.source.Category != category {
			if allowIndirect {
				newChildNodes := e.navigateChild(child, category, allowIndirect)
				newNodes = append(newNodes, newChildNodes...)
			}
			continue
		}
		newNodes = append(newNodes, child)
	}

	return newNodes
}

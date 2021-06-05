package parse

import "nli-go/lib/mentalese"

type Ellipsizer struct {

}

func NewEllipsizer() *Ellipsizer {
	return &Ellipsizer{}
}

// Returns a copy of `tree` where the `ellipsis` directions are included
func (e *Ellipsizer) Ellipsize(tree mentalese.ParseTreeNode) mentalese.ParseTreeNode {

	// quick check if this is necessary
	if !e.hasEllipsis(tree) { return tree }

	biDirTree := CreateBidirectionalParseTree(tree)
	biDirTree = biDirTree

	newTree := e.ellipsizeNode(biDirTree)

	return *newTree
}

func (e *Ellipsizer) hasEllipsis(node mentalese.ParseTreeNode) bool {
	present := len(node.Rule.Ellipsis) > 0
	for _, constituent := range node.Constituents {
		present = present || e.hasEllipsis(*constituent)
	}
	return present
}

func (e *Ellipsizer) ellipsizeNode(node *BidirectionalParseTreeNode) *mentalese.ParseTreeNode {
	newSource := node.source

	// add original constituents, having been ellipsized
	newConstituents := []*mentalese.ParseTreeNode{}
	for _, child := range node.children {
		newConstituents = append(newConstituents, e.ellipsizeNode(child))
	}

	// add new constituents, from ellipsis
	ellipsisConstituents := e.createEllipsisConstituents(node)
	newConstituents = append(newConstituents, ellipsisConstituents...)

	newSource.Constituents = newConstituents

	return &newSource
}

func (e *Ellipsizer) createEllipsisConstituents(node *BidirectionalParseTreeNode) []*mentalese.ParseTreeNode {
	ellipsisConstituents := []*mentalese.ParseTreeNode{}

	//source := node.source
	//for _, categoryPath := range source.Rule.Ellipsis {
	//	newConstituent := e.processCategoryPath(node, categoryPath)
	//	ellipsisConstituents = append(ellipsisConstituents, newConstituent)
	//}

	return ellipsisConstituents
}

//func (e *Ellipsizer) processCategoryPath(node *BidirectionalParseTreeNode, path CategoryPath) *ParseTreeNode {
//	// Create an error if no node could be found (or a clarification request?)
//}
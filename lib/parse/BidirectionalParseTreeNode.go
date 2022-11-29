package parse

import "nli-go/lib/mentalese"

// A simple wrapper around a parse tree node, adding a parent pointer
type BidirectionalParseTreeNode struct {
	source   *mentalese.ParseTreeNode
	parent   *BidirectionalParseTreeNode
	children []*BidirectionalParseTreeNode
}

func CreateBidirectionalParseTree(root *mentalese.ParseTreeNode) *BidirectionalParseTreeNode {
	return createBidirectionalParseTreeNode(nil, root)
}

func createBidirectionalParseTreeNode(parent *BidirectionalParseTreeNode, source *mentalese.ParseTreeNode) *BidirectionalParseTreeNode {
	node := BidirectionalParseTreeNode{
		source:   source,
		parent:   parent,
		children: []*BidirectionalParseTreeNode{},
	}

	for _, constituent := range source.Constituents {
		node.children = append(node.children, createBidirectionalParseTreeNode(&node, constituent))
	}

	return &node
}

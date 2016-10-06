package natlang

type SimpleParseTreeNode struct {
	SyntacticCategory string
	Word              string
	Children          []SimpleParseTreeNode
}

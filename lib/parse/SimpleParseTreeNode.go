package parse

type SimpleParseTreeNode struct {
	SyntacticCategory string
	Word              string
	Children          []SimpleParseTreeNode
}

package example2

type SimpleParseTreeNode struct {
	SyntacticCategory string
	Word              string
	Children          []SimpleParseTreeNode
}

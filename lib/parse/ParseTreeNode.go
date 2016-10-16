package parse

type ParseTreeNode struct {
	SyntacticCategory string
	Word              string
	Children          []ParseTreeNode
}

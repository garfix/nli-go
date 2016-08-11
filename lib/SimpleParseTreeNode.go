package lib

type SimpleParseTreeNode struct {
    SyntacticCategory string
    Word              string
    Children          []SimpleParseTreeNode
}

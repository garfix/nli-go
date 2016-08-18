package example1

type SimpleParseTreeNode struct {
    SyntacticCategory string
    Word              string
    Children          []SimpleParseTreeNode
}

package api

type Parser interface {
	Parse(words []string) []ParseTreeNode
}

type ParseTreeNode interface {
	IsLeafNode() bool
	GetConstituents() []*ParseTreeNode
	String() string
}

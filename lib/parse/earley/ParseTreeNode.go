package earley

import (
	"nli-go/lib/parse"
)

type ParseTreeNode struct {
	category     string
	constituents []ParseTreeNode
	form         string
	rule         parse.GrammarRule
}

func (node ParseTreeNode) IsLeafNode() bool {
	return len(node.constituents) == 0
}

func (node ParseTreeNode) String() string {

	body := ""

	if node.form != "" {
		body = node.form
	} else {
		sep := ""
		for _, child := range node.constituents {
			body += sep + child.String()
			sep = " "
		}
	}

	return "[" + node.category + " " + body + "]"
}

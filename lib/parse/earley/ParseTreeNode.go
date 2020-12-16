package earley

import (
	"nli-go/lib/api"
	"nli-go/lib/parse"
)

type ParseTreeNode struct {
	category     string
	constituents []*ParseTreeNode
	form         string
	rule         parse.GrammarRule
}

func (node ParseTreeNode) IsLeafNode() bool {
	return len(node.constituents) == 0
}

func (node ParseTreeNode) GetConstituents() []*api.ParseTreeNode {

	s := []api.ParseTreeNode{}

	for i, c := range node.constituents {
		q := *c
		s[i] = q
	}

	t := []*api.ParseTreeNode{}
	for i, c := range s {
		q := &c
		t[i] = q
	}

	return t
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

package parse

import "strings"

type ParseTreeNode struct {
	category     string
	constituents []*ParseTreeNode
	form         string
	rule         GrammarRule
}

func (node ParseTreeNode) IsLeafNode() bool {
	return len(node.constituents) == 0
}

func (node ParseTreeNode) GetConstituents() []*ParseTreeNode {
	return node.constituents
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

func (node ParseTreeNode) IndentedString(level int) string {

	prefix := ""
	if level > 0 {
		prefix = strings.Repeat("| ", level - 1) + "+- "
	}
	body := "\n" + prefix + node.category + " "

	if node.form != "" {
		body += node.form
	} else {
		for _, child := range node.constituents {
			body += child.IndentedString(level + 1)
		}
	}

	return body
}

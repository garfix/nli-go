package parse

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

func (node ParseTreeNode) IndentedString(indent string) string {

	body := ""

	if indent == "" {
		body = node.category + "\n"
	}

	for i, child := range node.constituents {
		if child.form != "" {
			body += indent + "+- " + child.category + " '" + child.form + "'\n"
			continue
		}

		body += indent + "+- " + child.category + "\n"
		newIndent := indent
		if i < len(node.constituents) - 1 {
			newIndent += "|  "
		} else {
			newIndent += "   "
		}
		body += child.IndentedString(newIndent)
	}

	return body
}

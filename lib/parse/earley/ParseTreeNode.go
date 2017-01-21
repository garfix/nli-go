package earley

type ParseTreeNode struct {
	category     string
	constituents []ParseTreeNode
	form         string
}

func (node ParseTreeNode) String() string {

	body := ""

	if node.form != "" {
		body = node.form
	} else {
		sep := ""
		for _, child := range node.constituents {
			body += sep + child.String()
			sep = ", "
		}
	}

	return "[" + node.category + " " + body + "]"
}
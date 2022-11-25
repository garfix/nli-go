package mentalese

type ParseTreeNode struct {
	Category     string
	Constituents []*ParseTreeNode
	Form         string
	Rule         GrammarRule
}

func (node *ParseTreeNode) IsLeafNode() bool {
	return len(node.Constituents) == 0
}

func (node *ParseTreeNode) GetConstituents() []*ParseTreeNode {
	return node.Constituents
}

func (node *ParseTreeNode) PartialCopy() ParseTreeNode {
	return ParseTreeNode{
		Category:     node.Category,
		Constituents: []*ParseTreeNode{},
		Form:         node.Form,
		Rule:         node.Rule.Copy(),
	}
}

func (node *ParseTreeNode) ReplaceVariable(fromVar string, toVar string) ParseTreeNode {
	constituents := []*ParseTreeNode{}
	for _, constituent := range node.Constituents {
		newConstituent := constituent.ReplaceVariable(fromVar, toVar)
		constituents = append(constituents, &newConstituent)
	}
	rule := node.Rule.Copy()
	rule = rule.ReplaceVariable(fromVar, toVar)
	return ParseTreeNode{
		Category:     node.Category,
		Constituents: constituents,
		Form:         node.Form,
		Rule:         rule,
	}
}

func (node *ParseTreeNode) GetVariablesRecursive() []string {
	variables := node.Rule.EntityVariables[0]
	for _, child := range node.Constituents {
		variables = append(variables, child.GetVariablesRecursive()...)
	}
	return variables
}

func (node *ParseTreeNode) String() string {

	body := ""

	if node.Form != "" {
		body = node.Form
	} else {
		sep := ""
		for _, child := range node.Constituents {
			body += sep + child.String()
			sep = " "
		}
	}

	return "[" + node.Category + " " + body + "]"
}

func (node *ParseTreeNode) IndentedString(indent string) string {

	body := ""

	if indent == "" {
		body = node.Category + "\n"
	}

	for i, child := range node.Constituents {
		if child.Form != "" {
			body += indent + "+- " + child.Category + " '" + child.Form + "'\n"
			continue
		}

		body += indent + "+- " + child.Category + "\n"
		newIndent := indent
		if i < len(node.Constituents)-1 {
			newIndent += "|  "
		} else {
			newIndent += "   "
		}
		body += child.IndentedString(newIndent)
	}

	return body
}

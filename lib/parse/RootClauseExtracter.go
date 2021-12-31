package parse

import "nli-go/lib/mentalese"

type RootClauseExtracter struct {
}

func NewRootClauseExtracter() *RootClauseExtracter {
	return &RootClauseExtracter{}
}

func (e *RootClauseExtracter) Extract(sentence *mentalese.ParseTreeNode) []*mentalese.ParseTreeNode {

	rootClauses := e.findRootClauses(sentence)

	if len(rootClauses) == 0 {
		rootClauses = []*mentalese.ParseTreeNode{sentence}
	}

	return rootClauses
}

func (e *RootClauseExtracter) findRootClauses(node *mentalese.ParseTreeNode) []*mentalese.ParseTreeNode {

	rootClauses := []*mentalese.ParseTreeNode{}

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == mentalese.TagRootClause {
			variable := tag.Arguments[0].TermValue
			for i, entityVariable := range node.Rule.EntityVariables {
				if i == 0 {
					continue
				}
				if entityVariable[0] == variable {
					rootClauses = append(rootClauses, node.Constituents[i-1])
				}
			}
		}
	}

	for _, child := range node.Constituents {
		rootClauses = append(rootClauses, e.findRootClauses(child)...)
	}

	return rootClauses
}

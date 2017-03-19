package earley

import (
	"nli-go/lib/parse"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

// The relationizer turns a parse tree into a relation set
type Relationizer struct {
	lexicon         *parse.Lexicon
	senseBuilder    parse.SenseBuilder
}

func NewRelationizer(lexicon *parse.Lexicon) Relationizer {
	return Relationizer{
		lexicon: lexicon,
		senseBuilder: parse.NewSenseBuilder(),
	}
}

func (relationizer Relationizer) Relationize(rootNode ParseTreeNode) mentalese.RelationSet {
	return relationizer.extractSenseFromNode(rootNode, relationizer.senseBuilder.GetNewVariable("Sentence"))
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable contains the actual variable used for the antecedent (for example: E1)
func (relationizer Relationizer) extractSenseFromNode(node ParseTreeNode, antecedentVariable string) mentalese.RelationSet {

	common.LogTree("extractSenseFromNode", node, antecedentVariable)

	relations := mentalese.RelationSet{}

	if node.IsLeafNode() {

		// leaf state rule: category -> word
		lexItem, _ := relationizer.lexicon.GetLexItem(node.form, node.category)
		lexItemRelations := relationizer.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relations = append(relations, lexItemRelations...)

	} else {

		variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariable, node.rule.EntityVariables)
		parentRelations := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)
		relations = append(relations, parentRelations...)

		// create relations for each of the children
		for i, childNode := range node.constituents {

			consequentVariable := variableMap[node.rule.EntityVariables[i + 1]]
			childRelations := relationizer.extractSenseFromNode(childNode, consequentVariable)
			relations = append(relations, childRelations...)
		}
	}

	common.LogTree("extractSenseFromNode", relations)

	return relations
}

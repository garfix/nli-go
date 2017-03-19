package earley

import (
	"nli-go/lib/parse"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

// The relationizer turns a parse tree into a relation set
// It also subsumes the range and quantifier relation sets inside its quantification relation
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

	relationSet := mentalese.RelationSet{}

	if node.IsLeafNode() {

		// leaf state rule: category -> word
		lexItem, _ := relationizer.lexicon.GetLexItem(node.form, node.category)
		lexItemRelations := relationizer.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relationSet = append(relationSet, lexItemRelations...)

	} else {

		variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariable, node.rule.EntityVariables)
		parentRelations := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)
		relationSet = append(relationSet, parentRelations...)

		// create relations for each of the children
		childSets := []mentalese.RelationSet{}
		for i, childNode := range node.constituents {

			consequentVariable := variableMap[node.rule.EntityVariables[i + 1]]
			childRelations := relationizer.extractSenseFromNode(childNode, consequentVariable)
			childSets = append(childSets, childRelations)
		}

		relationSet = relationizer.processChildRelations(relationSet, childSets, node.rule)
	}

	common.LogTree("extractSenseFromNode", relationSet)

	return relationSet
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) processChildRelations(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule parse.GrammarRule) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}
	extractedSetIndexes := map[int]bool{}

	for _, parentRelation := range parentSet {

		// special case
		if parentRelation.Predicate == "quantification!!" {

			qRelation := mentalese.Relation{}
			qRelation.Predicate = parentRelation.Predicate
			lastVariable := ""

			for _, argument := range parentRelation.Arguments {
				if argument.IsVariable() {
					lastVariable = argument.TermValue
					qRelation.Arguments = append(qRelation.Arguments, argument);
				} else if argument.IsRelationSet() {
					index, found := rule.GetConsequentIndexByVariable(lastVariable)
					if found {
						extractedSetIndexes[index] = true
						childSet := childSets[index]
						relationSetArgument := mentalese.Term{ TermType: mentalese.Term_relationSet, TermValueRelationSet: childSet }
						qRelation.Arguments = append(qRelation.Arguments, relationSetArgument);
					} else {
						panic("Relation set placeholder should be preceded by a variable from the rule")
					}
				}
			}

			newSet = append(newSet, qRelation)

		} else {
			newSet = append(newSet, parentRelation)
		}
	}

	for i, childSet := range childSets {

		// skip the child sets that were used in quantifications
		_, found := extractedSetIndexes[i]
		if !found {
			newSet = append(newSet, childSet...)
		}
	}

	return newSet
}
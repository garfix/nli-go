package earley

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
)

// The relationizer turns a parse tree into a relation set
// It also subsumes the range and quantifier relation sets inside its quantification relation
type Relationizer struct {
	lexicon      *parse.Lexicon
	senseBuilder parse.SenseBuilder
	log          *common.SystemLog
}

func NewRelationizer(lexicon *parse.Lexicon, log *common.SystemLog) Relationizer {
	return Relationizer{
		lexicon:      lexicon,
		senseBuilder: parse.NewSenseBuilder(),
		log:          log,
	}
}

func (relationizer Relationizer) Relationize(rootNode ParseTreeNode) mentalese.RelationSet {
	sense, _ := relationizer.extractSenseFromNode(rootNode, relationizer.senseBuilder.GetNewVariable("Sentence"))
	return sense
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (relationizer Relationizer) extractSenseFromNode(node ParseTreeNode, antecedentVariable string) (mentalese.RelationSet, bool) {

	relationizer.log.StartDebug("extractSenseFromNode", antecedentVariable, node.rule, node.rule.Sense)

	relationSet := mentalese.RelationSet{}
	var makeConstant = false

	if node.IsLeafNode() {

		// leaf state rule: category -> word
		lexItem, _, isRegExp := relationizer.lexicon.GetLexItem(node.form, node.category)
		lexItemRelations := relationizer.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relationSet = lexItemRelations

		// if the variable is assigned a constant now depends on it being a regexp; maybe this should be more explicit
		if isRegExp {
			makeConstant = true
		}

	} else {

		variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariable, node.rule.EntityVariables)

		// create relations for each of the children
		boundChildSets := []mentalese.RelationSet{}
		for i, childNode := range node.constituents {

			entityVariable := node.rule.EntityVariables[i+1]
			consequentVariable := variableMap[entityVariable]
			childRelations, makeConstant := relationizer.extractSenseFromNode(childNode, consequentVariable.TermValue)
			boundChildSets = append(boundChildSets, childRelations)

			if makeConstant {
				_, err := strconv.Atoi(childNode.form)
				if err == nil {
					variableMap[entityVariable] = mentalese.NewNumber(childNode.form)
				} else {
					variableMap[entityVariable] = mentalese.NewString(childNode.form)
				}
			}
		}

		boundParentSet := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)
		relationSet = relationizer.combineParentsAndChildren(boundParentSet, boundChildSets, node.rule)
	}

	relationizer.log.EndDebug("extractSenseFromNode", relationSet)

	return relationSet, makeConstant
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) combineParentsAndChildren(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule parse.GrammarRule) mentalese.RelationSet {

	relationizer.log.StartDebug("processChildRelations", parentSet, childSets, rule)

	newSet := parentSet.Copy()
	extractedSetIndexes := map[int]bool{}

	for i, childSet := range childSets {

		// skip the child sets that were used in quantifications
		_, found := extractedSetIndexes[i]
		if !found {
			newSet = append(newSet, childSet...)
		}
	}

	relationizer.log.EndDebug("processChildRelations", newSet)

	return newSet
}

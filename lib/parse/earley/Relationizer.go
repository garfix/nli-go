package earley

import (
	"nli-go/lib/parse"
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
	"fmt"
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
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (relationizer Relationizer) extractSenseFromNode(node ParseTreeNode, antecedentVariable string) mentalese.RelationSet {

	common.LogTree("extractSenseFromNode", antecedentVariable, node.rule, node.rule.Sense)
	common.LogTree("")

	relationSet := mentalese.RelationSet{}

	if node.IsLeafNode() {

		// leaf state rule: category -> word
		lexItem, _ := relationizer.lexicon.GetLexItem(node.form, node.category)
		lexItemRelations := relationizer.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relationSet = lexItemRelations

	} else {

		variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariable, node.rule.EntityVariables)

		// create relations for each of the children
		boundChildSets := []mentalese.RelationSet{}
		for i, childNode := range node.constituents {

			consequentVariable := variableMap[node.rule.EntityVariables[i + 1]]
			childRelations := relationizer.extractSenseFromNode(childNode, consequentVariable)
			boundChildSets = append(boundChildSets, childRelations)
		}

		boundParentSet := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)
		relationSet = relationizer.combineParentsAndChildren(boundParentSet, boundChildSets, node.rule)
	}

	common.LogTree("")
	common.LogTree("extractSenseFromNode", relationSet)

	return relationSet
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) combineParentsAndChildren(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule parse.GrammarRule) mentalese.RelationSet {

	common.LogTree("processChildRelations", parentSet, childSets, rule)

	newSet := mentalese.RelationSet{}
	extractedSetIndexes := map[int]bool{}
	compoundRelation := mentalese.Relation{}

	for i, parentRelation := range parentSet {

		// special case
		if parentRelation.Predicate == "quantification" {

			compoundRelation, extractedSetIndexes = relationizer.doQuantification(parentRelation, i, childSets, rule, extractedSetIndexes)
			newSet = append(newSet, compoundRelation)

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

	common.LogTree("processChildRelations", newSet)

	return newSet
}

// Example:
// relation = quantification(E1, [], D1, [])
// extractedSetIndexes = []
// childSets = [ [], [isa(E1, dog)], [], [isa(D1, every)] ]
// rule = np(E1) -> dp(D1) nbar(E1);
func (relationizer Relationizer) doQuantification(actualRelation mentalese.Relation, childIndex int, childSets []mentalese.RelationSet, rule parse.GrammarRule, extractedSetIndexes map[int]bool) (mentalese.Relation, map[int]bool) {

	common.LogTree("doQuantification", actualRelation, childSets, rule)

	formalRelation := rule.Sense[childIndex]
	lastVariable := ""

	for i, formalArgument := range formalRelation.Arguments {
		if formalArgument.IsVariable() {
			lastVariable = formalArgument.TermValue
		} else if formalArgument.IsRelationSet() {
			index, found := rule.GetConsequentIndexByVariable(lastVariable)
			if found {
				actualArgument := actualRelation.Arguments[i]
				extractedSetIndexes[index] = true
				subSet := append(actualArgument.TermValueRelationSet, childSets[index]...)
				relationSetArgument := mentalese.Term{ TermType: mentalese.Term_relationSet, TermValueRelationSet: subSet }
				actualRelation.Arguments[i] = relationSetArgument;
			} else {
				panic(fmt.Sprintf("Relation set placeholder should be preceded by a variable from the rule  %v %s", rule, lastVariable))
			}
		}
	}

	common.LogTree("doQuantification", actualRelation, extractedSetIndexes)

	return actualRelation, extractedSetIndexes
}
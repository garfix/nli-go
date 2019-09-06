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

	extractedSetIndexes := map[int]bool{}
	compoundRelation := mentalese.Relation{}

	// handle sem(1) sem(2)
	newSet1 := mentalese.RelationSet{}
	for i, parentRelation := range parentSet {
		compoundRelation, extractedSetIndexes = relationizer.includeChildSenses(parentRelation, i, childSets, rule, extractedSetIndexes)
		newSet1 = append(newSet1, compoundRelation)
	}

	// handle sem(parent)
	newSet2, extractedSetIndexes := relationizer.raiseChildSenses(newSet1, childSets, extractedSetIndexes)

	// skip the child sets that were incorporated in other relations
	for i, childSet := range childSets {
		_, found := extractedSetIndexes[i]
		if !found {
			newSet2 = append(newSet2, childSet...)
		}
	}

	relationizer.log.EndDebug("processChildRelations", newSet2)

	return newSet2
}

func (relationizer Relationizer) raiseChildSenses(parentRelations mentalese.RelationSet, childSets []mentalese.RelationSet, extractedSetIndexes map[int]bool) (mentalese.RelationSet, map[int]bool) {

	newParentRelations := parentRelations

	for s, childSet := range childSets {
		for r, childRelation := range childSet {
			for a, argument := range childRelation.Arguments {
				if argument.IsRelationSet() {
					for _, argumentRelation := range argument.TermValueRelationSet {
						if argumentRelation.Predicate == mentalese.PredicateSem && argumentRelation.Arguments[0].TermValue == mentalese.AtomParent {

							// the the sem of P is replaced by this quant
							newParentRelations = childSet.Copy()
							// the argument 'scope' in the quant of C is replaced by the current sem of P
							newParentRelations[r].Arguments[a] = mentalese.NewRelationSet(parentRelations)

							extractedSetIndexes[s] = true
						}
					}
				}
			}
		}
	}

	return newParentRelations, extractedSetIndexes
}

// Example:
// relation = quantification(E1, [], D1, [])
// extractedSetIndexes = []
// childSets = [ [], [isa(E1, dog)], [], [isa(D1, every)] ]
// rule = np(E1) -> dp(D1) nbar(E1);
func (relationizer Relationizer) includeChildSenses(parentRelation mentalese.Relation, childIndex int, childSets []mentalese.RelationSet, rule parse.GrammarRule, extractedSetIndexes map[int]bool) (mentalese.Relation, map[int]bool) {

	relationizer.log.StartDebug("includeChildSenses", parentRelation, childSets, rule)

	ruleRelation := rule.Sense[childIndex]

	for i, formalArgument := range ruleRelation.Arguments {
		if formalArgument.IsRelationSet() {
			firstRelation := formalArgument.TermValueRelationSet[0]
			if firstRelation.Predicate == mentalese.PredicateSem {
				index, err := strconv.Atoi(firstRelation.Arguments[0].TermValue)
				if err == nil {
					index = index - 1
					extractedSetIndexes[index] = true
					subSet := childSets[index]
					relationSetArgument := mentalese.Term{TermType: mentalese.TermRelationSet, TermValueRelationSet: subSet}
					parentRelation.Arguments[i] = relationSetArgument
				}
			}
		}
	}

	relationizer.log.EndDebug("includeChildSenses", parentRelation, extractedSetIndexes)

	return parentRelation, extractedSetIndexes
}
package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

// The relationizer turns a parse tree into a relation set
// It also subsumes the range and quantifier relation sets inside its quantification relation
type Relationizer struct {
	senseBuilder SenseBuilder
	log          *common.SystemLog
}

func NewRelationizer(variableGenerator *mentalese.VariableGenerator, log *common.SystemLog) *Relationizer {
	return &Relationizer{
		senseBuilder: NewSenseBuilder(variableGenerator),
		log:          log,
	}
}

func (relationizer Relationizer) Relationize(rootNode mentalese.ParseTreeNode, rootVariables []string) mentalese.RelationSet {
	if rootVariables == nil {
		rootVariables = []string{relationizer.senseBuilder.variableGenerator.GenerateVariable("Sentence").TermValue}
	}
	sense, constantBinding := relationizer.extractSenseFromNode(rootNode, rootVariables)
	sense = sense.BindSingle(constantBinding)
	return sense
}

func (relationizer Relationizer) ExtractTags(node mentalese.ParseTreeNode) mentalese.RelationSet {

	tags := node.Rule.Tag

	for _, childNode := range node.Constituents {
		tags = append(tags, relationizer.ExtractTags(*childNode)...)
	}

	return tags
}

func (relationizer Relationizer) ExtractIntents(node mentalese.ParseTreeNode) mentalese.RelationSet {

	intents := node.Rule.Intent

	for _, childNode := range node.Constituents {
		intents = append(intents, relationizer.ExtractIntents(*childNode)...)
	}

	return intents
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (relationizer Relationizer) extractSenseFromNode(node mentalese.ParseTreeNode, antecedentVariables []string) (mentalese.RelationSet, mentalese.Binding) {

	constantBinding := mentalese.NewBinding()
	relationSet := mentalese.RelationSet{}

	// create relations for each of the children
	boundChildSets := []mentalese.RelationSet{}
	for i, childNode := range node.Constituents {

		consequentVariables := node.Rule.GetConsequentVariables(i)

		childRelations, childConstantBinding := relationizer.extractSenseFromNode(*childNode, consequentVariables)
		boundChildSets = append(boundChildSets, childRelations)
		constantBinding = constantBinding.Merge(childConstantBinding)

		if node.Rule.GetConsequentPositionType(i) == mentalese.PosTypeRegExp {
			constantBinding.Set(antecedentVariables[0], mentalese.NewTermString(childNode.Form))
		}
	}

	boundParentSet := relationizer.senseBuilder.CreateGrammarRuleRelations(node.Rule.Sense)
	relationSet = relationizer.combineParentsAndChildren(boundParentSet, boundChildSets, node.Rule)

	return relationSet, constantBinding
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) combineParentsAndChildren(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule mentalese.GrammarRule) mentalese.RelationSet {

	referencedChildrenIndexes := []int{}
	compoundRelations := mentalese.RelationSet{}

	// process sem(Cat, Index)
	combination := mentalese.RelationSet{}
	for _, parentRelation := range parentSet {
		compoundRelations, referencedChildrenIndexes = relationizer.includeChildSenses(parentRelation, childSets, referencedChildrenIndexes, rule)
		combination = append(combination, compoundRelations...)
	}

	// prepend simple children
	restChildrenRelations := mentalese.RelationSet{}
	for i, childSet := range childSets {
		if !common.IntArrayContains(referencedChildrenIndexes, i) {
			restChildrenRelations = append(restChildrenRelations, childSet...)
		}
	}
	combination = append(restChildrenRelations, combination...)

	return combination
}

// replaces `sem(N)` in parentRelation
func (relationizer Relationizer) includeChildSenses(parentRelation mentalese.Relation, childSets []mentalese.RelationSet, childIndexes []int, rule mentalese.GrammarRule) (mentalese.RelationSet, []int) {

	newParentRelationSet := mentalese.RelationSet{}

	if parentRelation.Predicate == mentalese.PredicateSem {

		cat := parentRelation.Arguments[0].TermValue
		catIndex, err := strconv.Atoi(parentRelation.Arguments[1].TermValue)
		if err == nil {
			index := rule.FindConsequentIndex(cat, catIndex)
			if index != -1 {
				childIndexes = append(childIndexes, index)
				newParentRelationSet = childSets[index].Copy()
				for i := range newParentRelationSet {
					if parentRelation.Negate {
						newParentRelationSet[i].Negate = !newParentRelationSet[i].Negate
					}
				}
			} else {
				relationizer.log.AddError("$" + cat + strconv.Itoa(catIndex) + " not found in grammar rule")
			}
		}

	} else {

		newParentRelation := parentRelation.Copy()

		for i, formalArgument := range parentRelation.Arguments {
			if formalArgument.IsRelationSet() {
				newSet, newChildIndexes := relationizer.includeChildSensesInSet(formalArgument.TermValueRelationSet, childSets, childIndexes, rule)
				childIndexes = newChildIndexes
				newParentRelation.Arguments[i] = mentalese.NewTermRelationSet(newSet)
			} else if formalArgument.IsRule() {
				newRule := mentalese.Rule{}
				replaced := mentalese.RelationSet{}
				replaced, childIndexes = relationizer.includeChildSenses(formalArgument.TermValueRule.Goal, childSets, childIndexes, rule)
				newRule.Goal = replaced[0]
				replaced, childIndexes = relationizer.includeChildSensesInSet(formalArgument.TermValueRule.Pattern, childSets, childIndexes, rule)
				newRule.Pattern = replaced
				newParentRelation.Arguments[i] = mentalese.NewTermRule(newRule)
			} else if formalArgument.IsList() {
				panic("not implemented")
			}
		}

		newParentRelationSet = mentalese.RelationSet{newParentRelation}
	}

	return newParentRelationSet, childIndexes
}

func (relationizer Relationizer) includeChildSensesInSet(parentRelations mentalese.RelationSet, childSets []mentalese.RelationSet, childIndexes []int, rule mentalese.GrammarRule) (mentalese.RelationSet, []int) {

	newSet := mentalese.RelationSet{}
	for _, relation := range parentRelations {
		replacedChild, newChildIndexes := relationizer.includeChildSenses(relation, childSets, childIndexes, rule)
		childIndexes = newChildIndexes
		newSet = append(newSet, replacedChild...)
	}

	return newSet, childIndexes
}

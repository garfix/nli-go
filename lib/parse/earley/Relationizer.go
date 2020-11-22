package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
)

// The relationizer turns a parse tree into a relation set
// It also subsumes the range and quantifier relation sets inside its quantification relation
type Relationizer struct {
	senseBuilder parse.SenseBuilder
	log          *common.SystemLog
}

func NewRelationizer(log *common.SystemLog) *Relationizer {
	return &Relationizer{
		senseBuilder: parse.NewSenseBuilder(),
		log:          log,
	}
}

func (relationizer Relationizer) Relationize(rootNode ParseTreeNode, nameResolver *central.NameResolver) (mentalese.RelationSet, mentalese.Binding) {
	rootEntityVariable := relationizer.senseBuilder.GetNewVariable("Sentence")
	sense, nameBinding, constantBinding := relationizer.extractSenseFromNode(rootNode, nameResolver, []string{ rootEntityVariable } )
	sense = sense.BindSingle(constantBinding)
	return sense, nameBinding
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (relationizer Relationizer) extractSenseFromNode(node ParseTreeNode, nameResolver *central.NameResolver, antecedentVariables []string) (mentalese.RelationSet, mentalese.Binding, mentalese.Binding) {

	nameBinding := mentalese.NewBinding()
	constantBinding := mentalese.NewBinding()
	relationSet := mentalese.RelationSet{}

	if len(node.nameInformations) > 0 {
		firstAntecedentVariable := antecedentVariables[0]
		resolvedNameInformations := nameResolver.Resolve(node.nameInformations)
		for _, nameInformation := range resolvedNameInformations {
			nameBinding.Set(firstAntecedentVariable, mentalese.NewTermId(nameInformation.SharedId, nameInformation.EntityType))
		}
	}

	variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariables, node.rule.GetAntecedentVariables(), node.rule.GetAllConsequentVariables())

	// create relations for each of the children
	boundChildSets := []mentalese.RelationSet{}
	for i, childNode := range node.constituents {

		consequentVariables := node.rule.GetConsequentVariables(i)

		mappedConsequentVariables := []string{}
		for _, consequentVariable := range consequentVariables {
			mappedConsequentVariables = append(mappedConsequentVariables, variableMap[consequentVariable].TermValue)
		}

		childRelations, childNameBinding, childConstantBinding := relationizer.extractSenseFromNode(*childNode, nameResolver, mappedConsequentVariables)
		nameBinding = nameBinding.Merge(childNameBinding)
		boundChildSets = append(boundChildSets, childRelations)
		constantBinding = constantBinding.Merge(childConstantBinding)

		if node.rule.GetConsequentPositionType(i) == parse.PosTypeRegExp {
			constantBinding.Set(antecedentVariables[0], mentalese.NewTermString(childNode.form))
		}
	}

	variableMap = relationizer.senseBuilder.ExtendVariableMap(node.rule.Sense, variableMap)

	boundParentSet := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)

	relationSet = relationizer.combineParentsAndChildren(boundParentSet, boundChildSets, node.rule)

	return relationSet, nameBinding, constantBinding
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) combineParentsAndChildren(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule parse.GrammarRule) mentalese.RelationSet {

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
func (relationizer Relationizer) includeChildSenses(parentRelation mentalese.Relation, childSets []mentalese.RelationSet, childIndexes []int, rule parse.GrammarRule) (mentalese.RelationSet, []int) {

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
					if !parentRelation.Positive {
						newParentRelationSet[i].Positive = !newParentRelationSet[i].Positive
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

		newParentRelationSet = mentalese.RelationSet{ newParentRelation}
	}

	return newParentRelationSet, childIndexes
}

func (relationizer Relationizer) includeChildSensesInSet(parentRelations mentalese.RelationSet, childSets []mentalese.RelationSet, childIndexes []int, rule parse.GrammarRule) (mentalese.RelationSet, []int) {

	newSet := mentalese.RelationSet{}
	for _, relation := range parentRelations {
		replacedChild, newChildIndexes := relationizer.includeChildSenses(relation, childSets, childIndexes, rule)
		childIndexes = newChildIndexes
		newSet = append(newSet, replacedChild...)
	}

	return newSet, childIndexes
}
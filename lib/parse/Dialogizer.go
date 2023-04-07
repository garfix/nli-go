package parse

import "nli-go/lib/mentalese"

// Replaces the entity variables in a parse tree with discourse entities (discourse-wide variables)
type Dialogizer struct {
	variableGenerator *mentalese.VariableGenerator
	senseBuilder      SenseBuilder
}

func NewDialogizer(variableGenerator *mentalese.VariableGenerator) *Dialogizer {
	return &Dialogizer{
		variableGenerator: variableGenerator,
		senseBuilder:      NewSenseBuilder(variableGenerator),
	}
}

func (d *Dialogizer) Dialogize(node *mentalese.ParseTreeNode, rootVariables []string) *mentalese.ParseTreeNode {

	if rootVariables == nil {
		rootVariables = []string{d.variableGenerator.GenerateVariable("Sentence").TermValue}
	}

	return d.dialogizeNode(node, rootVariables)
}

func (d *Dialogizer) dialogizeNode(node *mentalese.ParseTreeNode, actualAntecedents []string) *mentalese.ParseTreeNode {

	dialogizedRule := mentalese.GrammarRule{}

	variableMap := d.senseBuilder.CreateVariableMap(actualAntecedents, node.Rule.GetAntecedentVariables(), node.Rule.GetAllConsequentVariables())

	variableMap = d.senseBuilder.ExtendVariableMap(node.Rule.Sense, variableMap)

	// add constituents
	newConstituents := []*mentalese.ParseTreeNode{}
	for i, child := range node.Constituents {

		childAntecedentVariables := node.Rule.GetConsequentVariables(i)

		childActualAntecedents := []string{}
		for _, antecedentVariable := range childAntecedentVariables {
			childActualAntecedents = append(childActualAntecedents, variableMap[antecedentVariable].TermValue)
		}

		newConstituent := d.dialogizeNode(child, childActualAntecedents)
		newConstituents = append(newConstituents, newConstituent)
	}

	dialogizedRule = d.dialogizeRule(node.Rule, variableMap)

	newSource := mentalese.ParseTreeNode{
		Category:     node.Category,
		Constituents: newConstituents,
		Form:         node.Form,
		Rule:         dialogizedRule,
	}
	newSource.Constituents = newConstituents

	return &newSource
}

func (d *Dialogizer) dialogizeRule(rule mentalese.GrammarRule, variableMap map[string]mentalese.Term) mentalese.GrammarRule {
	binding := mentalese.NewBinding()
	for key, value := range variableMap {
		binding.Set(key, value)
	}

	newRule := rule.BindSimple(binding)

	return newRule
}

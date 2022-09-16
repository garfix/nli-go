package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type AnaphoraResolver2 struct {
	dialogContext *DialogContext
	meta          *mentalese.Meta
	messenger     api.ProcessMessenger
}

func NewAnaphoraResolver2(dialogContext *DialogContext, meta *mentalese.Meta, messenger api.ProcessMessenger) *AnaphoraResolver2 {
	return &AnaphoraResolver2{
		dialogContext: dialogContext,
		meta:          meta,
		messenger:     messenger,
	}
}

func (resolver *AnaphoraResolver2) Resolve(root *mentalese.ParseTreeNode, request mentalese.RelationSet, binding mentalese.Binding) (*mentalese.ParseTreeNode, mentalese.RelationSet, mentalese.BindingSet, string) {

	newBindings := mentalese.InitBindingSet(binding)
	collection := NewAnaphoraResolverCollection()

	resolvedRoot := resolver.resolveNode(root, binding, collection)

	resolvedRequest := request

	return resolvedRoot, resolvedRequest, newBindings, collection.output
}

func (resolver *AnaphoraResolver2) resolveNode(node *mentalese.ParseTreeNode, binding mentalese.Binding, collection *AnaphoraResolverCollection) *mentalese.ParseTreeNode {

	for _, childNode := range node.GetConstituents() {
		resolver.resolveNode(childNode, binding, collection)
	}

	variables := node.Rule.GetAntecedentVariables()
	for _, variable := range variables {
		//tags := resolver.dialogContext.EntityTags.GetTagPredicates(variable)
		tags := node.Rule.Tag
		for _, tag := range tags {
			resolvedVariable := variable
			if tag.Predicate == mentalese.TagSortalReference {
			}
			if tag.Predicate == mentalese.TagReference {
				resolvedVariable = resolver.reference(variable, node, binding, collection)
				if resolvedVariable != variable {
					collection.AddReference(variable, mentalese.NewTermVariable(resolvedVariable))
				}
			}
			if tag.Predicate == mentalese.TagLabeledReference {
			}
			if tag.Predicate == mentalese.TagReflectiveReference {
			}
			resolvedVariable = resolvedVariable
		}
	}

	return node
}

func (resolver *AnaphoraResolver2) reference(variable string, node *mentalese.ParseTreeNode, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	set := resolver.dialogContext.EntityDefinitions.Get(variable) //node.Rule.Sense
	resolvedVariable := variable

	// if the variable has been bound already, don't try to look for a reference
	_, found := resolver.dialogContext.EntityBindings.Get(variable)
	if found {
		return variable
	}

	found, referentVariable, referentValue := resolver.findReferent(variable, set, binding)
	if found {
		if referentVariable != "" {
			collection.AddReplacement(variable, referentVariable)
			resolvedVariable = referentVariable
		} else {
			collection.AddReference(variable, referentValue)
		}
	} else {

		newBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(binding))
		if newBindings.GetLength() > 1 {
			// ask the user which of the specified entities he/she means
			collection.output = "I don't understand which one you mean"
		}
	}

	return resolvedVariable
}

func (resolver *AnaphoraResolver2) findReferent(variable string, set mentalese.RelationSet, binding mentalese.Binding) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		// there may be 1..n groups (bindings)
		referentVariable := group.Variable

		if !resolver.dialogContext.CheckAgreement(variable, referentVariable) {
			continue
		}

		// if there's 1 group and its id = "", it is unbound
		isBound := group.values[0].Id != ""

		if isBound {
			// empty set ("it")
			if len(set) == 0 {
				found = true
				foundVariable = referentVariable
				break
			}
		}

		for _, referent := range group.values {

			if referent.Id == "" {
				continue
			}

			b := mentalese.NewBinding()
			value := mentalese.NewTermId(referent.Id, referent.Sort)
			b.Set(variable, value)

			refBinding := binding.Merge(b)
			testRangeBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(refBinding))

			if testRangeBindings.GetLength() > 0 {
				found = true
				if len(group.values) == 1 {
					foundVariable = referentVariable
				} else {
					foundTerm = value
				}
				goto end
			}
		}

	}

end:

	return found, foundVariable, foundTerm
}

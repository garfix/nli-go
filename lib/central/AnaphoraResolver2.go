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

	println("---")
	println(request.String())

	newBindings := mentalese.InitBindingSet(binding)
	collection := NewAnaphoraResolverCollection()

	resolvedRoot := resolver.resolveNode(root, binding, collection)

	resolvedRequest := request
	resolvedRequest = resolver.replaceOneAnaphora(resolvedRequest, collection.oneAnaphors)

	// binding the reference variable to one of the values of its referent (when the referent is a group and we need just one element from it)
	for fromVariable, value := range collection.values {
		binding.Set(fromVariable, value)
	}

	// update the binding by replacing the variables
	for fromVariable, toVariable := range collection.replacements {
		// replace the other variables in the set
		resolvedRequest = resolvedRequest.ReplaceTerm(mentalese.NewTermVariable(fromVariable), mentalese.NewTermVariable(toVariable))
		value, found := resolver.dialogContext.EntityBindings.Get(toVariable)
		if found {
			if value.IsList() {
				tempBindings := mentalese.NewBindingSet()
				for _, item := range value.TermValueList {
					for _, b := range newBindings.GetAll() {
						newBinding := b.Copy()
						newBinding.Set(toVariable, item)
						tempBindings.Add(newBinding)
					}
				}
				newBindings = tempBindings
			} else {
				newBindings.SetAll(toVariable, value)
			}
		}
	}

	println(resolvedRequest.String())
	println(resolvedRoot.String())

	return resolvedRoot, resolvedRequest, newBindings, collection.output
}

func (resolver *AnaphoraResolver2) replaceOneAnaphora(set mentalese.RelationSet, replacements map[string]mentalese.RelationSet) mentalese.RelationSet {
	newSet := set
	for variable, definition := range replacements {
		relation := mentalese.NewRelation(false, mentalese.PredicateReferenceSlot, []mentalese.Term{mentalese.NewTermVariable(variable)})
		newSet = newSet.ReplaceRelation(relation, definition)
	}
	return newSet
}

func (resolver *AnaphoraResolver2) resolveNode(node *mentalese.ParseTreeNode, binding mentalese.Binding, collection *AnaphoraResolverCollection) *mentalese.ParseTreeNode {

	for _, relation := range node.Rule.Sense {
		if relation.Predicate == mentalese.PredicateReferenceSlot {
			variable := relation.Arguments[0].TermValue
			found, referentVariable := resolver.sortalReference(variable)
			if found {
				oneAnaphor := resolver.dialogContext.EntityDefinitions.Get(referentVariable).
					ReplaceTerm(mentalese.NewTermVariable(referentVariable), mentalese.NewTermVariable(variable))
				collection.AddOneAnaphor(variable, oneAnaphor)
			}
		} else if relation.Predicate == mentalese.PredicateQuant {
			variable := relation.Arguments[1].TermValue
			resolver.dialogContext.ClauseList.GetLastClause().AddEntity(variable)
		}
	}

	variables := node.Rule.GetAntecedentVariables()
	for _, variable := range variables {
		tags := node.Rule.Tag
		for _, tag := range tags {
			resolvedVariable := variable
			if tag.Predicate == mentalese.TagSortalReference {
			}
			if tag.Predicate == mentalese.TagReference {
				resolvedVariable = resolver.reference(variable, binding, collection)
				if resolvedVariable != variable {
					collection.AddReference(variable, mentalese.NewTermVariable(resolvedVariable))
				}
			}
			if tag.Predicate == mentalese.TagLabeledReference {
				label := tag.Arguments[1].TermValue
				resolvedVariable = resolver.labeledReference(variable, label, binding, collection)
			}
			if tag.Predicate == mentalese.TagReflectiveReference {
			}
			resolvedVariable = resolvedVariable
		}
	}

	for _, childNode := range node.GetConstituents() {
		resolver.resolveNode(childNode, binding, collection)
	}

	return node
}

func (resolver *AnaphoraResolver2) reference(variable string, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

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

func (resolver *AnaphoraResolver2) labeledReference(variable string, label string, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	aLabel, found := resolver.dialogContext.EntityLabels.GetLabel(label)
	if found {
		resolver.dialogContext.EntityLabels.IncreaseActivation(label)
		// use the reference of the existing label
		referencedVariable := aLabel.label
		collection.AddReference(variable, mentalese.NewTermVariable(referencedVariable))
		return referencedVariable
	} else {
		referencedVariable := resolver.reference(variable, binding, collection)
		if referencedVariable != variable {
			// create a new label
			resolver.dialogContext.EntityLabels.SetLabel(label, variable)
		}
		return referencedVariable
	}
}

func (resolver *AnaphoraResolver2) sortalReference(variable string) (bool, string) {

	//sortRelationSet := mentalese.RelationSet{}
	found := false
	foundVariable := ""

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		//sort := ""
		//
		//// if their are multiple values, their sorts should match
		//for _, ref := range group.values {
		//	if sort == "" {
		//		sort = ref.Sort
		//	} else if sort != ref.Sort {
		//		sort = ""
		//		break
		//	}
		//}
		//
		//if sort == "" {
		//	continue
		//}
		//
		//sortInfo, found := resolver.meta.GetSortProperty(sort)
		//if !found {
		//	continue
		//}
		//
		//if sortInfo.Entity.Equals(mentalese.RelationSet{}) {
		//	continue
		//}
		//
		//sortRelationSet = sortInfo.Entity.ReplaceTerm(mentalese.NewTermVariable(mentalese.IdVar), mentalese.NewTermVariable(variable))
		foundVariable = group.Variable
		found = true
		break
	}

	return found, foundVariable
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

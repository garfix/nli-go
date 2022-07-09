package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type AnaphoraResolver struct {
	dialogContext *DialogContext
	meta          *mentalese.Meta
	messenger     api.ProcessMessenger
}

func NewAnaphoraResolver(dialogContext *DialogContext, meta *mentalese.Meta, messenger api.ProcessMessenger) *AnaphoraResolver {
	return &AnaphoraResolver{
		dialogContext: dialogContext,
		meta:          meta,
		messenger:     messenger,
	}
}

func (resolver *AnaphoraResolver) Resolve(set mentalese.RelationSet, binding mentalese.Binding) (mentalese.RelationSet, mentalese.BindingSet, string) {

	collection := NewAnaphoraResolverCollection()

	// extend the set with one one-anaphora resolutions, replace variables, and collect the other matches
	newSet := resolver.resolveSet(set, binding, collection)

	// binding the reference variable to one of the values of its referent (when the referent is a group and we need just one element from it)
	for fromVariable, value := range collection.values {
		binding.Set(fromVariable, value)
	}

	newBindings := mentalese.InitBindingSet(binding)

	// update the binding by replacing the variables
	for fromVariable, toVariable := range collection.replacements {
		// replace the other variables in the set
		newSet = newSet.ReplaceTerm(mentalese.NewTermVariable(fromVariable), mentalese.NewTermVariable(toVariable))
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

	//fmt.Println(newSet)
	//fmt.Println(newBindings)

	return newSet, newBindings, collection.output
}

func (resolver *AnaphoraResolver) resolveSet(set mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}

	for _, relation := range set {

		newRelation := relation

		if relation.Predicate == mentalese.PredicateQuant {
			newRelation = resolver.resolveQuant(relation, binding, collection)
		} else if relation.Predicate == mentalese.PredicateEventReference {
			resolver.resolveEvent(relation, binding, collection)
		} else {
			newRelation = resolver.resolveArguments(relation, binding, collection)
		}

		newSet = append(newSet, newRelation)
	}

	return newSet
}

func (resolver *AnaphoraResolver) resolveArguments(relation mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) mentalese.Relation {

	newRelation := relation.Copy()

	for i, argument := range relation.Arguments {
		if argument.IsRelationSet() {
			newRelation.Arguments[i].TermValueRelationSet = resolver.resolveSet(argument.TermValueRelationSet, binding, collection)
		}
	}

	return newRelation
}

func (resolver *AnaphoraResolver) resolveQuant(quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) mentalese.Relation {
	rangeVar := quant.Arguments[1].TermValue

	resolvedVariable := rangeVar

	tags := resolver.dialogContext.EntityTags.GetTagPredicates(rangeVar)
	if common.StringArrayContains(tags, mentalese.TagSortalReference) {
		sortRelationSet := resolver.sortalReference(quant)
		quant = quant.Copy()
		quant.Arguments[2] = mentalese.NewTermRelationSet(append(sortRelationSet, quant.Arguments[2].TermValueRelationSet...))
	}
	if common.StringArrayContains(tags, mentalese.TagReference) {
		resolvedVariable = resolver.reference(quant, binding, collection)

		if rangeVar != resolvedVariable {
			quant = quant.ReplaceTerm(mentalese.NewTermVariable(rangeVar), mentalese.NewTermVariable(resolvedVariable))
			resolver.dialogContext.ReplaceVariable(rangeVar, resolvedVariable)
		}
	}
	if common.StringArrayContains(tags, mentalese.TagLabeledReference) {
		labels := resolver.dialogContext.EntityTags.GetTagsByPredicate(rangeVar, mentalese.TagLabeledReference)
		resolvedVariable = resolver.labeledReference(labels[0], quant, binding, collection)

		if rangeVar != resolvedVariable {
			quant = quant.ReplaceTerm(mentalese.NewTermVariable(rangeVar), mentalese.NewTermVariable(resolvedVariable))
			resolver.dialogContext.ReplaceVariable(rangeVar, resolvedVariable)
		}
	}
	if common.StringArrayContains(tags, mentalese.TagReflectiveReference) {
		resolvedVariable = resolver.reflectiveReference(quant, binding, collection)

		if rangeVar != resolvedVariable {
			quant = quant.ReplaceTerm(mentalese.NewTermVariable(rangeVar), mentalese.NewTermVariable(resolvedVariable))
			resolver.dialogContext.ReplaceVariable(rangeVar, resolvedVariable)
		}
	}

	resolver.dialogContext.ClauseList.GetLastClause().AddEntity(resolvedVariable)

	return quant
}

func (resolver *AnaphoraResolver) labeledReference(labelRelation mentalese.Relation, quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	variable := quant.Arguments[1].TermValue
	label := labelRelation.Arguments[1].TermValue

	aLabel, found := resolver.dialogContext.EntityLabels.GetLabel(label)
	if found {
		resolver.dialogContext.EntityLabels.IncreaseActivation(label)
		// use the reference of the existing label
		referencedVariable := aLabel.label
		collection.AddReference(variable, mentalese.NewTermVariable(referencedVariable))
		return referencedVariable
	} else {
		referencedVariable := resolver.reference(quant, binding, collection)
		if referencedVariable != variable {
			// create a new label
			resolver.dialogContext.EntityLabels.SetLabel(label, variable)
		}
		return referencedVariable
	}
}

func (resolver *AnaphoraResolver) reference(quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	variable := quant.Arguments[1].TermValue
	set := quant.Arguments[2].TermValueRelationSet
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

func (resolver *AnaphoraResolver) reflectiveReference(quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	variable := quant.Arguments[1].TermValue
	resolvedVariable := variable

	// if the variable has been bound already, don't try to look for a reference
	_, found := resolver.dialogContext.EntityBindings.Get(variable)
	if found {
		return variable
	}

	// presuming the target of the reflection is the previously handled entity
	// (because it lives in the parent quant)
	clause := resolver.dialogContext.ClauseList.GetLastClause()
	if len(clause.ResolvedEntities) > 0 {
		lastVariable := clause.ResolvedEntities[len(clause.ResolvedEntities)-1]
		if resolver.dialogContext.CheckAgreement(variable, lastVariable) {
			collection.AddReplacement(variable, lastVariable)
			resolvedVariable = lastVariable
		}
	}

	return resolvedVariable
}

func (resolver *AnaphoraResolver) sortalReference(quant mentalese.Relation) mentalese.RelationSet {

	sortRelationSet := mentalese.RelationSet{}

	variable := quant.Arguments[1].TermValue

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		sort := ""

		// if their are multiple values, their sorts should match
		for _, ref := range group.values {
			if sort == "" {
				sort = ref.Sort
			} else if sort != ref.Sort {
				sort = ""
				break
			}
		}

		if sort == "" {
			continue
		}

		sortInfo, found := resolver.meta.GetSortProperty(sort)
		if !found {
			continue
		}

		if sortInfo.Entity.Equals(mentalese.RelationSet{}) {
			continue
		}

		sortRelationSet = sortInfo.Entity.ReplaceTerm(mentalese.NewTermVariable(mentalese.IdVar), mentalese.NewTermVariable(variable))
		break
	}

	return sortRelationSet
}

func (resolver *AnaphoraResolver) findReferent(variable string, set mentalese.RelationSet, binding mentalese.Binding) (bool, string, mentalese.Term) {

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

func (resolver *AnaphoraResolver) resolveEvent(event mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) {

	variable := event.Arguments[0].TermValue

	_, found := binding.Get(variable)
	if found {
		return
	}

	tags := resolver.dialogContext.EntityTags.GetTagPredicates(variable)
	resolvedVariable := variable
	if common.StringArrayContains(tags, mentalese.TagLabeledReference) {
		labels := resolver.dialogContext.EntityTags.GetTagsByPredicate(variable, mentalese.TagLabeledReference)
		resolvedVariable = resolver.labeledEventReference(labels[0], event, binding, collection)

		if variable != resolvedVariable {
			event = event.ReplaceTerm(mentalese.NewTermVariable(variable), mentalese.NewTermVariable(resolvedVariable))
			resolver.dialogContext.ReplaceVariable(variable, resolvedVariable)
			return
		}
	}

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		// there may be 1..n groups (bindings)
		referentVariable := group.Variable

		if len(group.values) > 1 {
			continue
		}

		sort := group.values[0].Sort
		if sort == mentalese.SortEvent {
			println(variable, " => ", referentVariable)
			collection.AddReplacement(variable, referentVariable)
			break
		}
	}
}

func (resolver *AnaphoraResolver) labeledEventReference(labelRelation mentalese.Relation, event mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	variable := event.Arguments[0].TermValue
	label := labelRelation.Arguments[1].TermValue

	aLabel, found := resolver.dialogContext.EntityLabels.GetLabel(label)
	if found {
		resolver.dialogContext.EntityLabels.IncreaseActivation(label)
		// use the reference of the existing label
		referencedVariable := aLabel.label
		collection.AddReference(variable, mentalese.NewTermVariable(referencedVariable))
		return referencedVariable
	} else {
		//referencedVariable := resolver.reference(event, binding, collection)
		//if referencedVariable != variable {
		//	// create a new label
		//	resolver.dialogContext.EntityLabels.SetLabel(label, variable)
		//}
		//return referencedVariable
		resolver.dialogContext.EntityLabels.SetLabel(label, variable)
		return variable
	}
}

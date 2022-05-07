package central

import (
	"nli-go/lib/api"
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

	resolver.resolveSet(set, binding, collection)

	newBindings := mentalese.InitBindingSet(binding)

	// replace the reference variable by the variable of its referent
	for fromVariable, toVariable := range collection.replacements {
		set = set.ReplaceTerm(mentalese.NewTermVariable(fromVariable), mentalese.NewTermVariable(toVariable))
		value, found := resolver.dialogContext.DiscourseEntities.Get(toVariable)
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

	// binding the reference variable to one of the values of its referent (when the referent is a group and we need just one element from it)
	for fromVariable, value := range collection.references {
		resolver.dialogContext.DiscourseEntities.Set(fromVariable, value)
		binding.Set(fromVariable, value)
	}

	//println(set.String())
	//println(newBindings.String())

	return set, newBindings, collection.output
}

func (resolver *AnaphoraResolver) resolveSet(set mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection) {

	for _, relation := range set {
		if relation.Predicate == mentalese.PredicateQuant {
			resolver.resolveQuant(relation, binding, collection)
		} else {
			resolver.resolveArguments(relation, binding, collection)
		}
	}
}

func (resolver *AnaphoraResolver) resolveArguments(relation mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) {

	for _, argument := range relation.Arguments {
		if argument.IsRelationSet() {
			resolver.resolveSet(argument.TermValueRelationSet, binding, collection)
		}
	}
}

func (resolver *AnaphoraResolver) resolveQuant(quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) (mentalese.Relation, mentalese.Binding, string) {

	output := ""
	newBinding := binding
	rangeVar := quant.Arguments[1].TermValue
	newQuant := quant

	tags := resolver.dialogContext.TagList.GetTags(rangeVar)

	for _, tag := range tags {
		switch tag.Predicate {
		case mentalese.TagReference:
			resolver.reference(quant, binding, collection)
		}
	}

	return newQuant, newBinding, output
}

func (resolver *AnaphoraResolver) reference(quant mentalese.Relation, binding mentalese.Binding, collection *AnaphoraResolverCollection) {

	variable := quant.Arguments[1].TermValue
	set := quant.Arguments[2].TermValueRelationSet

	//println("reference? " + set.String())

	found, referentVariable, referentValue := resolver.findReferent(variable, set, binding)
	if found {
		if referentVariable != "" {
			collection.AddReplacement(variable, referentVariable)
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
}

func (resolver *AnaphoraResolver) findReferent(variable string, set mentalese.RelationSet, binding mentalese.Binding) (bool, string, mentalese.Term) {

	//newBindings := mentalese.NewBindingSet()

	//unscopedSense := request.UnScope()

	//if resolver.dialogContext.DiscourseEntities.ContainsVariable(variable) {
	//	value := resolver.dialogContext.DiscourseEntities.MustGet(variable)
	//	newBindings := mentalese.NewBindingSet()
	//	if value.IsList() {
	//		for _, item := range value.TermValueList {
	//			newBinding := mentalese.NewBinding()
	//			newBinding.Set(variable, item)
	//			newBindings.Add(newBinding)
	//		}
	//	} else {
	//		newBinding := mentalese.NewBinding()
	//		newBinding.Set(variable, value)
	//		newBindings.Add(newBinding)
	//	}
	//
	//	return newBindings
	//}

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		// there may be 1..n groups (bindings)
		referentVariable := group[0].Variable

		// if there's 1 group and its id = "", it is unbound
		isBound := group[0].Id != ""

		//if resolver.isReflexive(unscopedSense, variable, ref) {
		//	continue
		//}

		if isBound {
			// empty set ("it")
			if len(set) == 0 {
				println(" NO SET " + referentVariable)
				found = true
				foundVariable = referentVariable
				break
			}
		}

		//if !resolver.quickAcceptabilityCheck(variable, ref.Sort, set) {
		//	continue
		//}

		for _, referent := range group {

			if referent.Id == "" {
				continue
			}

			b := mentalese.NewBinding()
			value := mentalese.NewTermId(referent.Id, referent.Sort)
			b.Set(variable, value)

			refBinding := binding.Merge(b)
			testRangeBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(refBinding))

			if testRangeBindings.GetLength() > 0 {

				println(" => " + referent.String() + " " + set.String())

				found = true
				if len(group) == 1 {
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

// checks if a (irreflexive) pronoun does not refer to another element in a same relation
func (base *AnaphoraResolver) isReflexive(unscopedSense mentalese.RelationSet, referenceVariable string, antecedent EntityReference) bool {

	antecedentvariable := antecedent.Variable

	if antecedentvariable == "" {
		return false
	}

	reflexive := false
	for _, relation := range unscopedSense {
		ref := false
		ante := false
		for _, argument := range relation.Arguments {
			if argument.IsVariable() {
				if argument.TermValue == antecedentvariable {
					ante = true
				}
				if argument.TermValue == referenceVariable {
					ref = true
				}
			}
		}
		if ref && ante {
			reflexive = true
		}
	}

	return reflexive
}

func (resolver *AnaphoraResolver) quickAcceptabilityCheck(variable string, sort string, relations mentalese.RelationSet) bool {

	accepted := false

	for _, relation := range relations {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {
				argumentEntityType := resolver.meta.GetSort(relation.Predicate, i)

				if argumentEntityType == "" || resolver.meta.MatchesSort(argumentEntityType, sort) {
					accepted = true
					break
				}
			}
		}
	}

	return accepted
}

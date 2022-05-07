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

func (resolver *AnaphoraResolver) Resolve(set mentalese.RelationSet, binding mentalese.Binding) (mentalese.RelationSet, mentalese.Binding, string) {

	collection := NewAnaphoraResolverCollection()

	resolver.resolveSet(set, binding, collection)

	for fromVariable, toVariable := range collection.replacements {
		set = set.ReplaceTerm(mentalese.NewTermVariable(fromVariable), mentalese.NewTermVariable(toVariable))
		value, found := resolver.dialogContext.DiscourseEntities.Get(toVariable)
		if found {
			binding.Set(toVariable, value)
		}
	}

	return set, binding, collection.output
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

	referentVariable := resolver.findReferent(variable, set, binding)
	if referentVariable != "" {
		//println("reference!")
		//println(variable + " " + referentVariable)
		// replace
		collection.AddReplacement(variable, referentVariable)
	} else {

		newBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(binding))
		if newBindings.GetLength() > 1 {
			// ask the user which of the specified entities he/she means
			collection.output = "I don't understand which one you mean"
		}
	}
}

func (resolver *AnaphoraResolver) findReferent(variable string, set mentalese.RelationSet, binding mentalese.Binding) string {

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

	foundVariable := ""

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
			b.Set(variable, mentalese.NewTermId(referent.Id, referent.Sort))

			refBinding := binding.Merge(b)
			testRangeBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(refBinding))

			if testRangeBindings.GetLength() > 0 {
				println(" => " + referent.String() + " " + set.String())
				foundVariable = referentVariable
				goto end
			}
		}

	}

end:

	return foundVariable
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

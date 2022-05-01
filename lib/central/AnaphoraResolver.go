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

func (resolver *AnaphoraResolver) Resolve(request mentalese.RelationSet, binding mentalese.Binding) (mentalese.RelationSet, mentalese.Binding, string) {

	newRelations := mentalese.RelationSet{}
	newBinding := binding.Copy()
	output := ""

	for _, relation := range request {
		newRelation := relation
		if relation.Predicate == mentalese.PredicateQuant {
			newRelation, newBinding, output = resolver.resolveQuant(relation, newBinding, request)
			if output != "" {
				break
			}
		} else {
			newRelation, newBinding, output = resolver.resolveArguments(relation, newBinding)
			if output != "" {
				break
			}
		}
		newRelations = append(newRelations, newRelation)
	}

	return newRelations, newBinding, output
}

func (resolver *AnaphoraResolver) resolveArguments(relation mentalese.Relation, binding mentalese.Binding) (mentalese.Relation, mentalese.Binding, string) {

	newRelation := relation.Copy()
	newBinding := binding.Copy()
	output := ""

	for i, argument := range relation.Arguments {
		if argument.IsRelationSet() {
			newArgument := mentalese.RelationSet{}
			newArgument, newBinding, output = resolver.Resolve(argument.TermValueRelationSet, newBinding)
			newRelation.Arguments[i] = mentalese.NewTermRelationSet(newArgument)
			if output != "" {
				break
			}
		}
	}

	return newRelation, newBinding, output
}

func (resolver *AnaphoraResolver) resolveQuant(quant mentalese.Relation, binding mentalese.Binding, request mentalese.RelationSet) (mentalese.Relation, mentalese.Binding, string) {

	output := ""
	rangeVar := quant.Arguments[1].TermValue
	rangeRelations := quant.Arguments[2].TermValueRelationSet

	tags := resolver.dialogContext.TagList.GetTags(rangeVar)

	for _, tag := range tags {
		switch tag.Predicate {
		case mentalese.TagReference:
			output = resolver.reference(rangeVar, rangeRelations, request, binding)
		}
	}

	return quant, binding, output
}

func (resolver *AnaphoraResolver) reference(variable string, set mentalese.RelationSet, request mentalese.RelationSet, binding mentalese.Binding) string {

	output := ""

	if !resolver.doBackReference(variable, set, request, binding).IsEmpty() {

	} else {

		if set[0].Predicate == "x" { //mentalese.PredicateDefiniteBackReference {
			newBindings := resolver.messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(binding))
			if newBindings.GetLength() > 1 {
				// ask the user which of the specified entities he/she means
				output = "I don't understand which one you mean1"
			}
		}
	}

	return output
}

func (resolver *AnaphoraResolver) doBackReference(variable string, set mentalese.RelationSet, request mentalese.RelationSet, binding mentalese.Binding) mentalese.BindingSet {

	newBindings := mentalese.NewBindingSet()

	unscopedSense := request.UnScope()

	if resolver.dialogContext.DiscourseEntities.ContainsVariable(variable) {
		value := resolver.dialogContext.DiscourseEntities.MustGet(variable)
		newBindings := mentalese.NewBindingSet()
		if value.IsList() {
			for _, item := range value.TermValueList {
				newBinding := mentalese.NewBinding()
				newBinding.Set(variable, item)
				newBindings.Add(newBinding)
			}
		} else {
			newBinding := mentalese.NewBinding()
			newBinding.Set(variable, value)
			newBindings.Add(newBinding)
		}

		return newBindings
	}

	for _, group := range resolver.dialogContext.GetAnaphoraQueue() {

		ref := group[0]

		newBindings1 := mentalese.NewBindingSet()
		for _, r1 := range group {
			b := mentalese.NewBinding()
			b.Set(variable, mentalese.NewTermId(r1.Id, r1.Sort))

			refBinding := binding.Merge(b)
			newBindings1.Add(refBinding)
		}

		if resolver.isReflexive(unscopedSense, variable, ref) {
			continue
		}

		// empty set ("it")
		if len(set) == 0 {
			newBindings = newBindings1
			break
		}

		if !resolver.quickAcceptabilityCheck(variable, ref.Sort, set) {
			continue
		}

		testRangeBindings := mentalese.BindingSet{}

		if set[0].Predicate == mentalese.PredicateDefiniteBackReference {

		} else {

			testRangeBindings = resolver.messenger.ExecuteChildStackFrame(set, newBindings1)

			if testRangeBindings.GetLength() == 1 {
				newBindings = testRangeBindings
				break
			}

		}

	}

	return newBindings
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

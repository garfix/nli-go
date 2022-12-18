package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type AnaphoraResolver struct {
	clauseList        *mentalese.ClauseList
	entityBindings    *mentalese.EntityBindings
	entityTags        *mentalese.TagList
	entitySorts       *mentalese.EntitySorts
	entityLabels      *mentalese.EntityLabels
	entityDefinitions *mentalese.EntityDefinitions
	meta              *mentalese.Meta
	messenger         api.ProcessMessenger
}

func NewAnaphoraResolver(clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entityTags *mentalese.TagList, entitySorts *mentalese.EntitySorts, entityLabels *mentalese.EntityLabels, entityDefinitions *mentalese.EntityDefinitions, meta *mentalese.Meta, messenger api.ProcessMessenger) *AnaphoraResolver {
	return &AnaphoraResolver{
		clauseList:        clauseList,
		entityBindings:    entityBindings,
		entityTags:        entityTags,
		entitySorts:       entitySorts,
		entityLabels:      entityLabels,
		entityDefinitions: entityDefinitions,
		meta:              meta,
		messenger:         messenger,
	}
}

func (resolver *AnaphoraResolver) collectCoArguments(set mentalese.RelationSet, collection *AnaphoraResolverCollection) {

	for _, relation := range set {
		if relation.Predicate == mentalese.PredicateCheck || relation.Predicate == mentalese.PredicateDo {
			body := relation.Arguments[mentalese.CheckBodyIndex].TermValueRelationSet

			for _, bodyRelation1 := range body {
				for _, argument1 := range bodyRelation1.Arguments {

					if !argument1.IsVariable() {
						continue
					}

					for _, bodyRelation2 := range body {
						for _, argument2 := range bodyRelation2.Arguments {

							if !argument2.IsVariable() {
								continue
							}

							if argument1.TermValue != argument2.TermValue {
								collection.AddCoArgument(argument1.TermValue, argument2.TermValue)
							}
						}
					}
				}
			}
		}

		for _, argument := range relation.Arguments {
			if argument.IsRelationSet() {
				resolver.collectCoArguments(argument.TermValueRelationSet, collection)
			}
		}
	}
}

func (resolver *AnaphoraResolver) Resolve(root *mentalese.ParseTreeNode, request mentalese.RelationSet, binding mentalese.Binding) (*mentalese.ParseTreeNode, mentalese.RelationSet, mentalese.BindingSet, string) {

	//println("---")
	// println(request.String())

	newBindings := mentalese.InitBindingSet(binding)
	collection := NewAnaphoraResolverCollection()

	resolver.collectCoArguments(request, collection)

	resolver.resolveNode(root, binding, collection)

	resolvedRequest := request
	resolvedRequest = resolver.replaceOneAnaphora(resolvedRequest, collection.oneAnaphors)

	// binding the reference variable to one of the values of its referent (when the referent is a group and we need just one element from it)
	for fromVariable, value := range collection.values {
		binding.Set(fromVariable, value)
	}

	// update the binding by replacing the variables
	for fromVariable, toVariable := range collection.replacements {

		resolver.clauseList.GetLastClause().ReplaceVariable(fromVariable, toVariable)

		// replace the other variables in the set
		resolvedRequest = resolvedRequest.ReplaceTerm(mentalese.NewTermVariable(fromVariable), mentalese.NewTermVariable(toVariable))
		value, found := resolver.entityBindings.Get(toVariable)
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

	// println(resolvedRequest.String())
	//println(resolvedRoot.String())
	// println(newBindings.String())

	resolvedTree := resolver.clauseList.GetLastClause().ParseTree

	return resolvedTree, resolvedRequest, newBindings, collection.output
}

func (resolver *AnaphoraResolver) replaceOneAnaphora(set mentalese.RelationSet, replacements map[string]mentalese.RelationSet) mentalese.RelationSet {
	newSet := set
	for variable, definition := range replacements {
		relation := mentalese.NewRelation(false, mentalese.PredicateReferenceSlot, []mentalese.Term{mentalese.NewTermVariable(variable)})
		newSet = resolver.ReplaceRelation(newSet, relation, definition)
	}
	return newSet
}

// Replaces all occurrences from from to to
func (resolver *AnaphoraResolver) ReplaceRelation(source mentalese.RelationSet, placeholder mentalese.Relation, replacement mentalese.RelationSet) mentalese.RelationSet {

	target := mentalese.RelationSet{}
	placeholderFound := false

	for _, relation := range source {
		if relation.Equals(placeholder) {
			placeholderFound = true
			break
		}
	}

	peacefulReplacement := replacement
	if placeholderFound {
		sourceSorts := resolver.findSorts(source)
		peacefulReplacement = mentalese.RelationSet{}
		for _, relation := range replacement {
			if !resolver.conflicts(relation, sourceSorts) {
				peacefulReplacement = append(peacefulReplacement, relation)
			}
		}
	}

	for _, relation := range source {
		if relation.Equals(placeholder) {
			target = append(target, peacefulReplacement...)
		} else {
			newArguments := []mentalese.Term{}
			for _, argument := range relation.Arguments {
				newArgument := argument
				if argument.IsRelationSet() {
					newArgument = mentalese.NewTermRelationSet(resolver.ReplaceRelation(argument.TermValueRelationSet, placeholder, replacement))
				}
				newArguments = append(newArguments, newArgument)
			}

			newRelation := mentalese.NewRelation(relation.Negate, relation.Predicate, newArguments)
			target = append(target, newRelation)
		}
	}

	return target
}

func (resolver *AnaphoraResolver) findSorts(set mentalese.RelationSet) []string {

	sorts := []string{}

	for _, relation := range set {
		sorts = append(sorts, resolver.findSortsSingle(relation)...)
	}

	return sorts
}

func (resolver *AnaphoraResolver) findSortsSingle(relation mentalese.Relation) []string {

	sorts := []string{}

	isa := mentalese.NewRelation(false, mentalese.PredicateHasSort, []mentalese.Term{
		mentalese.NewTermAtom(relation.GetPredicateWithoutNamespace()),
		mentalese.NewTermVariable("Type"),
	})
	bindings := resolver.messenger.ExecuteChildStackFrame(mentalese.RelationSet{isa}, mentalese.InitBindingSet(mentalese.NewBinding()))
	for _, binding := range bindings.GetAll() {
		sort := binding.MustGet("Type").TermValue
		sorts = append(sorts, sort)
	}

	return sorts
}

func (resolver *AnaphoraResolver) conflicts(relation mentalese.Relation, sorts []string) bool {
	relationSorts := resolver.findSortsSingle(relation)
	conflicts := false
	for _, relationSort := range relationSorts {
		if common.StringArrayContains(sorts, relationSort) {
			conflicts = true
			break
		}
	}
	return conflicts
}

func (resolver *AnaphoraResolver) resolveNode(node *mentalese.ParseTreeNode, binding mentalese.Binding, collection *AnaphoraResolverCollection) {

	variables := node.Rule.GetAntecedentVariables()
	for _, variable := range variables {
		tags := node.Rule.Tag
		for _, tag := range tags {
			resolvedVariable := variable
			if tag.Predicate == mentalese.TagReference {
				resolvedVariable = resolver.reference(variable, binding, collection)
				if resolvedVariable != variable {
					collection.AddReplacement(variable, resolvedVariable)
				}
			}
			if tag.Predicate == mentalese.TagLabeledReference {
				label := tag.Arguments[1].TermValue
				condition := tag.Arguments[2].TermValueRelationSet
				resolvedVariable = resolver.labeledReference(variable, label, condition, binding, collection)
				if resolvedVariable != variable {
					collection.AddReplacement(variable, resolvedVariable)
				}
			}
			// if tag.Predicate == mentalese.TagReflectiveReference {
			// }
		}
	}

	for _, relation := range node.Rule.Sense {
		if relation.Predicate == mentalese.PredicateReferenceSlot {
			variable := relation.Arguments[0].TermValue
			found, referentVariable := resolver.sortalReference(variable)
			if found {
				oneAnaphor := resolver.entityDefinitions.Get(referentVariable).
					ReplaceTerm(mentalese.NewTermVariable(referentVariable), mentalese.NewTermVariable(variable))
				collection.AddOneAnaphor(variable, oneAnaphor)

			}
		} else if relation.Predicate == mentalese.PredicateQuant {
			variable := relation.Arguments[1].TermValue
			resolver.clauseList.GetLastClause().AddEntity(variable)
		}
	}

	for _, childNode := range node.GetConstituents() {
		resolver.resolveNode(childNode, binding, collection)
	}
}

func (resolver *AnaphoraResolver) reference(variable string, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	set := resolver.entityDefinitions.Get(variable) //node.Rule.Sense
	resolvedVariable := variable

	// if the variable has been bound already, don't try to look for a reference
	_, found := resolver.entityBindings.Get(variable)
	if found {
		return variable
	}

	found, referentVariable, referentValue := resolver.findReferent(variable, set, binding, collection)
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
		} else {
			// non-anaphoric reference ("the red cube", found in the scene)
			// newBinding := newBindings.Get(0)
			// value, found := newBinding.Get(variable)
			// if found {
			// 	collection.AddReference(variable, value)
			// 	resolver.entityBindings.Set(variable, value)
			// }
		}
	}

	return resolvedVariable
}

func (resolver *AnaphoraResolver) labeledReference(variable string, label string, condition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection) string {

	aLabel, found := resolver.entityLabels.GetLabel(label)
	if found {
		resolver.entityLabels.IncreaseActivation(label)
		// use the reference of the existing label
		referencedVariable := aLabel.GetVariable()
		return referencedVariable
	} else {
		referencedVariable := resolver.reference(variable, binding, collection)
		if referencedVariable != variable {

			conditionBindings := resolver.messenger.ExecuteChildStackFrame(condition, mentalese.InitBindingSet(binding))
			if conditionBindings.GetLength() > 0 {
				// create a new label
				resolver.entityLabels.SetLabel(label, referencedVariable)
			} else {
				return variable
			}

		}
		return referencedVariable
	}
}

func (resolver *AnaphoraResolver) sortalReference(variable string) (bool, string) {

	found := false
	foundVariable := ""

	queue := GetAnaphoraQueue(resolver.clauseList, resolver.entityBindings, resolver.entitySorts)
	for _, group := range queue {

		foundVariable = group.Variable

		// go:reference_slot() is inside a quant with the same variable
		if foundVariable == variable {
			continue
		}

		definition := resolver.entityDefinitions.Get(foundVariable)
		if definition.IsEmpty() {
			continue
		}

		typeFound := false
		for _, relation := range definition {
			if relation.Predicate == mentalese.PredicateHasSort {
				typeFound = true
			}
		}
		if !typeFound {
			continue
		}

		found = true
		break
	}

	return found, foundVariable
}

func (resolver *AnaphoraResolver) findReferent(variable string, set mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}

	groups := GetAnaphoraQueue(resolver.clauseList, resolver.entityBindings, resolver.entitySorts)
	for _, group := range groups {

		// there may be 1..n groups (bindings)
		referentVariable := group.Variable

		// the entity itself is in the queue
		if referentVariable == variable {
			continue
		}

		agreementChecker := NewAgreementChecker()

		agree, _, _ := agreementChecker.CheckForCategoryConflictBetween(variable, referentVariable, resolver.entityTags)
		if !agree {
			continue
		}

		// disallow reflective references
		if collection.IsCoArgument(variable, referentVariable) {
			continue
		}

		// if there's 1 group and its id = "", it is unbound
		isBound := group.values[0].Id != ""

		sameSentence := group.SentenceDistance == 0

		if isBound || sameSentence {
			// empty set ("it")
			if len(set) == 0 {
				if group.values[0].Sort != "" {
					found = true
					foundVariable = referentVariable
					break
				}
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

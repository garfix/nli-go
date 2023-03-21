package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type AnaphoraResolver struct {
	log               *common.SystemLog
	clauseList        *mentalese.ClauseList
	entityBindings    *mentalese.EntityBindings
	entityTags        *mentalese.TagList
	entitySorts       *mentalese.EntitySorts
	entityLabels      *mentalese.EntityLabels
	entityDefinitions *mentalese.EntityDefinitions
	meta              *mentalese.Meta
	messenger         api.ProcessMessenger
	sortFinder        SortFinder
}

func NewAnaphoraResolver(log *common.SystemLog, clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entityTags *mentalese.TagList, entitySorts *mentalese.EntitySorts, entityLabels *mentalese.EntityLabels, entityDefinitions *mentalese.EntityDefinitions, meta *mentalese.Meta, messenger api.ProcessMessenger) *AnaphoraResolver {
	return &AnaphoraResolver{
		log:               log,
		clauseList:        clauseList,
		entityBindings:    entityBindings,
		entityTags:        entityTags,
		entitySorts:       entitySorts,
		entityLabels:      entityLabels,
		entityDefinitions: entityDefinitions,
		meta:              meta,
		messenger:         messenger,
		sortFinder:        NewSortFinder(meta, messenger),
	}
}

func (resolver *AnaphoraResolver) Resolve(root *mentalese.ParseTreeNode, request mentalese.RelationSet, binding mentalese.Binding) (*mentalese.ParseTreeNode, mentalese.RelationSet, mentalese.BindingSet, string) {

	// println("---")
	// println(request.String())
	// println(binding.String())

	// prepare
	collection := NewAnaphoraResolverCollection()
	NewCoArgumentCollector().collectCoArguments(request, collection)

	resolver.resolveNode(root, binding, collection)

	// post process
	resolvedRequest, newBindings := resolver.processCollection(request, binding, collection)
	resolvedTree := resolver.clauseList.GetLastClause().ParseTree

	// println("---")
	// println(resolvedRequest.IndentedString("\t"))
	// println(newBindings.String())
	// println(resolvedRoot.String())

	return resolvedTree, resolvedRequest, newBindings, collection.output
}

func (resolver *AnaphoraResolver) processCollection(request mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection) (mentalese.RelationSet, mentalese.BindingSet) {

	newBindings := mentalese.InitBindingSet(binding)

	resolvedRequest := resolver.replaceOneAnaphora(request, collection.oneAnaphors)

	// binding the reference variable to one of the values of its referent (when the referent is a group and we need just one element from it)
	for fromVariable, value := range collection.values {
		binding.Set(fromVariable, value)
		resolver.entityBindings.Set(fromVariable, value)
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

	return resolvedRequest, newBindings
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
		reflective := false
		tags := node.Rule.Tag
		for _, tag := range tags {
			if tag.Predicate == mentalese.TagReflective {
				reflective = true
			}
		}
		for _, tag := range tags {

			resolvedVariable := variable

			if tag.Predicate == mentalese.TagReference {
				referenceSort := tag.Arguments[1].TermValue
				resolvedVariable = resolver.reference(variable, referenceSort, binding, collection, reflective)
				if resolvedVariable != variable {
					collection.AddReplacement(variable, resolvedVariable)
				}
			}

			if tag.Predicate == mentalese.TagLabeledReference {
				label := tag.Arguments[1].TermValue
				referenceSort := tag.Arguments[2].TermValue
				resolvedVariable = resolver.labeledReference(variable, referenceSort, label, binding, collection, reflective)
				if resolvedVariable != variable {
					collection.AddReplacement(variable, resolvedVariable)
				}
			}
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

func (resolver *AnaphoraResolver) reference(variable string, referenceSort string, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) string {

	entityDefinition := resolver.entityDefinitions.Get(variable)
	resolvedVariable := variable

	// if the variable has been bound already, don't try to look for a reference
	_, found := resolver.entityBindings.Get(variable)
	if found {
		return variable
	}

	found, referentVariable, referentValue := resolver.findAnaphoricReferent(variable, referenceSort, entityDefinition, binding, collection, reflective)
	if found {
		// found anaphoric referent (within dialog)
		if referentVariable != "" {
			collection.AddReplacement(variable, referentVariable)
			resolvedVariable = referentVariable
		} else {
			collection.AddReference(variable, referentValue)
		}
	} else {
		if len(entityDefinition) == 0 {
			collection.output = "I don't understand what you are referring to"
		} else {
			// try to find a non-anaphoric referent (outside the dialog), in the scene
			newBindings := resolver.messenger.ExecuteChildStackFrame(entityDefinition, mentalese.InitBindingSet(binding))
			if newBindings.GetLength() == 0 {
				collection.output = "I don't understand what you are referring to"
			} else if newBindings.GetLength() > 1 {
				// ask the user which of the specified entities he/she means
				collection.output = "I don't understand which one you mean"
			}
		}
	}

	return resolvedVariable
}

func (resolver *AnaphoraResolver) labeledReference(variable string, referenceSort string, label string, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) string {

	aLabel, found := resolver.entityLabels.GetLabel(label)
	if found {
		// check if this referent is still acceptable
		oldVariable := aLabel.GetVariable()
		oldSort := resolver.entitySorts.GetSort(oldVariable)
		oldId, _ := resolver.entityBindings.Get(oldVariable)
		entityDefinition := resolver.entityDefinitions.Get(variable)
		match, _, _ := resolver.matchReferenceToReferent(variable, referenceSort, oldVariable, oldSort, oldId.TermValue, entityDefinition, binding, collection, reflective)
		if match {
			resolver.entityLabels.IncreaseActivation(label)
			// use the reference of the existing label
			referencedVariable := aLabel.GetVariable()
			return referencedVariable
		}
	}

	referencedVariable := resolver.reference(variable, referenceSort, binding, collection, reflective)
	if referencedVariable != variable {
		resolver.entityLabels.SetLabel(label, referencedVariable)
	}
	return referencedVariable
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

func (resolver *AnaphoraResolver) findAnaphoricReferent(variable string, referenceSort string, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}

	groups := GetAnaphoraQueue(resolver.clauseList, resolver.entityBindings, resolver.entitySorts)
	for _, group := range groups {

		// check if the referent is a group
		if len(group.values) == 1 {

			referent := group.values[0]

			found, foundVariable, foundTerm = resolver.matchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, entityDefinition, binding, collection, reflective)
			if found {
				goto end
			}

		} else {

			// try to match an element in the group
			for _, referent := range group.values {

				found, foundVariable, foundTerm = resolver.matchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, entityDefinition, binding, collection, reflective)
				if found {
					goto end
				}
			}
		}

	}

end:

	if found {
		if foundVariable != "" {
			resolver.log.AddProduction("ref", "accept "+foundVariable+"\n")
		} else {
			resolver.log.AddProduction("ref", "accept "+foundTerm.String()+"\n")
		}
	} else {
		resolver.log.AddProduction("ref", "reject all\n")
	}

	return found, foundVariable, foundTerm
}

func (resolver *AnaphoraResolver) matchReferenceToReferent(variable string, referenceSort string, referentVariable string, referentSort string, referentId string, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}
	agree := false
	mostSpecificFound := false

	resolver.log.AddProduction("\nresolving", variable+"\n")

	// the entity itself is in the queue
	// should not be possible
	if referentVariable == variable {
		resolver.log.AddProduction("ref", referentVariable+" equals "+variable+"\n")
		goto end
	}

	if referentSort == "" {
		resolver.log.AddProduction("ref", referentVariable+" has no sort\n")
		goto end
	}
	_, mostSpecificFound = resolver.sortFinder.getMostSpecific(referenceSort, referentSort)
	if !mostSpecificFound {
		resolver.log.AddProduction("ref", referentVariable+" ("+referentSort+") does not have common sort with "+variable+" ("+referenceSort+")\n")
		goto end
	}

	agree, _, _ = NewAgreementChecker().CheckForCategoryConflictBetween(variable, referentVariable, resolver.entityTags)
	if !agree {
		resolver.log.AddProduction("ref", referentVariable+" does not agree with "+variable+"\n")
		goto end
	}

	if reflective {
		if !collection.IsCoArgument(variable, referentVariable) {
			resolver.log.AddProduction("ref", referentVariable+" is not co-argument "+variable+"\n")
			goto end
		}
	} else {
		if collection.IsCoArgument(variable, referentVariable) {
			resolver.log.AddProduction("ref", referentVariable+" is co-argument of "+variable+"\n")
			goto end
		}
	}

	// is this a definite reference?
	if len(entityDefinition) == 0 {
		// no: we're done
		found = true
		foundVariable = referentVariable
	} else {
		// yes, it is a definite reference
		// a definite reference can only be checked against an id
		if referentId == "" {
			resolver.log.AddProduction("ref", referentVariable+" has no id "+entityDefinition.String()+"\n")
			goto end
		} else {
			b := mentalese.NewBinding()
			value := mentalese.NewTermId(referentId, referentSort)
			b.Set(variable, value)

			refBinding := binding.Merge(b)
			testRangeBindings := resolver.messenger.ExecuteChildStackFrame(entityDefinition, mentalese.InitBindingSet(refBinding))
			if testRangeBindings.GetLength() > 0 {
				// found: bind the reference variable to the id of the referent
				// (don't replace variable)
				found = true
				foundTerm = value
				goto end
			} else {
				resolver.log.AddProduction("ref", referentVariable+" could not be bound\n")
				goto end
			}
		}
	}

end:

	return found, foundVariable, foundTerm
}

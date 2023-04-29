package central

import (
	"fmt"
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
	messenger         api.ProcessMessenger
	referentFinder    *ReferentFinder
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
		messenger:         messenger,
		referentFinder:    NewReferentFinder(log, meta, messenger, clauseList, entityBindings, entityDefinitions, entityTags, entitySorts),
	}
}

func (resolver *AnaphoraResolver) Resolve(root *mentalese.ParseTreeNode, request mentalese.RelationSet, binding mentalese.Binding) (*mentalese.ParseTreeNode, mentalese.RelationSet, mentalese.BindingSet, string, string) {

	// println("---")
	// println(request.IndentedString("  "))
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

	resolver.log.AddProduction("Resolved request", resolvedRequest.IndentedString("  "))

	return resolvedTree, resolvedRequest, newBindings, collection.output, collection.remark
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

	for _, relation := range node.Rule.Sense.UnScope() {
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

	referents, ambiguous := resolver.referentFinder.FindAnaphoricReferents(variable, referenceSort, entityDefinition, binding, collection, reflective, false)
	if len(referents) > 0 {
		referentVariable := referents[0].Variable
		referentValue := referents[0].Term

		// found anaphoric referent (within dialog)
		if referentVariable != "" {
			collection.AddReplacement(variable, referentVariable)
			resolvedVariable = referentVariable
		} else {
			collection.AddReference(variable, referentValue)
		}

		if ambiguous {
			fmt.Printf("\n===\n%v\n%v\n", variable, referents)
			// collection.remark = "AMBIGUOUS!"
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
		match, _, _ := resolver.referentFinder.MatchReferenceToReferent(variable, referenceSort, oldVariable, oldSort, oldId.TermValue, 0, entityDefinition, binding, collection, reflective)
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

	referents, _ := resolver.referentFinder.FindAnaphoricReferents(
		variable,
		mentalese.SortEntity, mentalese.RelationSet{}, mentalese.NewBinding(), NewAnaphoraResolverCollection(), false, true)

	foundVariable := ""
	found := len(referents) > 0
	if found {
		foundVariable = referents[0].Variable
		resolver.log.AddProduction("ref", variable+" resolves to "+foundVariable+"\n")
	} else {
		resolver.log.AddProduction("ref", variable+" has no referent\n")
	}

	return found, foundVariable
}

package central

import (
	"nli-go/lib/mentalese"
)

type EntityDefinitionsExtracter struct {
	dialogContext *DialogContext
}

func NewEntityDefinitionsExtracter(dialogContext *DialogContext) *EntityDefinitionsExtracter {
	return &EntityDefinitionsExtracter{
		dialogContext: dialogContext,
	}
}

func (e *EntityDefinitionsExtracter) Extract(set mentalese.RelationSet) {
	for _, relation := range set {
		if relation.Predicate == mentalese.PredicateQuant {
			definition := e.removeSelfReferences(relation.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet)
			e.AddDefinition(relation.Arguments[mentalese.QuantRangeVariableIndex].TermValue, definition)
		}
		for _, argument := range relation.Arguments {
			if argument.IsRelationSet() {
				e.Extract(argument.TermValueRelationSet)
			}
		}
	}
}

func (e *EntityDefinitionsExtracter) removeSelfReferences(set mentalese.RelationSet) mentalese.RelationSet {
	newSet := mentalese.RelationSet{}

	for _, relation := range set {
		if !e.containsSelfReference(relation) {
			newSet = append(newSet, relation)
		}
	}

	return newSet
}

func (e *EntityDefinitionsExtracter) containsSelfReference(relation mentalese.Relation) bool {

	contains := relation.Predicate == mentalese.PredicateReferenceSlot

	for _, argument := range relation.Arguments {
		if argument.IsRelationSet() {
			for _, child := range argument.TermValueRelationSet {
				contains = contains || e.containsSelfReference(child)
			}
		}
	}

	return contains
}

func (e *EntityDefinitionsExtracter) AddDefinition(variable string, definition mentalese.RelationSet) {
	e.dialogContext.EntityDefinitions.Add(variable, definition)
}

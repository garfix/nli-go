package central

import "nli-go/lib/mentalese"

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
			e.AddDefinition(relation.Arguments[mentalese.QuantRangeVariableIndex].TermValue, relation.Arguments[mentalese.QuantRangeSetIndex].TermValueRelationSet)
		}
		for _, argument := range relation.Arguments {
			if argument.IsRelationSet() {
				e.Extract(argument.TermValueRelationSet)
			}
		}
	}
}

func (e *EntityDefinitionsExtracter) AddDefinition(variable string, definition mentalese.RelationSet) {
	e.dialogContext.EntityDefinitions.Add(variable, definition)
}

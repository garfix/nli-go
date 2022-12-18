package central

import "nli-go/lib/mentalese"

type CoArgumentCollector struct {
}

func NewCoArgumentCollector() *CoArgumentCollector {
	return &CoArgumentCollector{}
}

func (c *CoArgumentCollector) collectCoArguments(set mentalese.RelationSet, collection *AnaphoraResolverCollection) {

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
				c.collectCoArguments(argument.TermValueRelationSet, collection)
			}
		}
	}
}

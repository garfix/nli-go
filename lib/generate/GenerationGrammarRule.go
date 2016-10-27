package generate

import "nli-go/lib/mentalese"

type GenerationGrammarRule struct {
	Antecedent          mentalese.Relation
	Consequents         mentalese.RelationSet
	Condition           mentalese.RelationSet
}

func (rule GenerationGrammarRule) String() string {

	return rule.Antecedent.String() + " :- " + rule.Consequents.String() + " { " + rule.Condition.String() + " }"
}
package mentalese

type Rule struct {
	Goal    Relation
	Pattern RelationSet
}

func (rule Rule) BindSingle(binding Binding) Rule {
	return Rule{
		Goal: rule.Goal.BindSingle(binding),
		Pattern: rule.Pattern.BindSingle(binding),
	}
}
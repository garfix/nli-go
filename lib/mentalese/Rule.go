package mentalese

type Rule struct {
	Goal    Relation
	Pattern RelationSet
}

func (rule Rule) ImportBinding(binding Binding) Rule {
	return Rule{
		Goal: rule.Goal.BindSingleRelationSingleBinding(binding),
		Pattern: rule.Pattern.ImportBinding(binding),
	}
}
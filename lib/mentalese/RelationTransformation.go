package mentalese

type RelationTransformation struct {
	Pattern     RelationSet
	Replacement RelationSet
	Condition   RelationSet
}

func (transformation RelationTransformation) BindSingle(binding Binding) RelationTransformation {
	return RelationTransformation{
		Pattern: transformation.Pattern.BindSingle(binding),
		Replacement: transformation.Replacement.BindSingle(binding),
		Condition: transformation.Condition.BindSingle(binding),
	}
}
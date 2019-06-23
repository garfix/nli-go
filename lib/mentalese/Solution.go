package mentalese

type Solution struct {
	Condition   RelationSet
	Transformations []RelationTransformation
	NoResults   ResultHandler
	SomeResults ResultHandler
}

func (solution Solution) BindSingle(binding Binding) Solution {

	boundTransformations := []RelationTransformation{}
	for _, transformation := range solution.Transformations {
		boundTransformations = append(boundTransformations, transformation.BindSingle(binding))
	}

	return Solution{
		Condition: solution.Condition.BindSingle(binding),
		Transformations: boundTransformations,
		NoResults: solution.NoResults.Bind(binding),
		SomeResults: solution.SomeResults.Bind(binding),
	}
}
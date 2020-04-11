package mentalese

type Solution struct {
	Condition   RelationSet
	Transformations []Rule
	Result Term
	Responses   []ResultHandler
}

func (solution Solution) BindSingle(binding Binding) Solution {

	boundTransformations := []Rule{}
	for _, transformation := range solution.Transformations {
		boundTransformations = append(boundTransformations, transformation.BindSingle(binding))
	}

	newResponses := []ResultHandler{}
	for _, response := range solution.Responses {
		newResponses = append(newResponses, response.Bind(binding))
	}

	return Solution{
		Condition: solution.Condition.BindSingle(binding),
		Transformations: boundTransformations,
		Result: solution.Result.Bind(binding),
		Responses: newResponses,
	}
}
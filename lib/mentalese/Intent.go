package mentalese

type Intent struct {
	Condition       RelationSet
	Transformations []Rule
	Responses       []ResultHandler
}

func (intent Intent) BindSingle(binding Binding) Intent {

	boundTransformations := []Rule{}
	for _, transformation := range intent.Transformations {
		boundTransformations = append(boundTransformations, transformation.BindSingle(binding))
	}

	newResponses := []ResultHandler{}
	for _, response := range intent.Responses {
		newResponses = append(newResponses, response.Bind(binding))
	}

	return Intent{
		Condition:       intent.Condition.BindSingle(binding),
		Transformations: boundTransformations,
		Responses:       newResponses,
	}
}

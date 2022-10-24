package mentalese

type Intent struct {
	Condition RelationSet
	Responses []ResultHandler
}

func (intent Intent) BindSingle(binding Binding) Intent {
	newResponses := []ResultHandler{}
	for _, response := range intent.Responses {
		newResponses = append(newResponses, response.Bind(binding))
	}

	return Intent{
		Condition: intent.Condition.BindSingle(binding),
		Responses: newResponses,
	}
}

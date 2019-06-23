package mentalese

type ResultHandler struct {
	Preparation RelationSet
	Answer RelationSet
}

func (handler ResultHandler) Bind(binding Binding) ResultHandler {
	return ResultHandler{
		Preparation: handler.Preparation.BindSingle(binding),
		Answer: handler.Answer.BindSingle(binding),
	}
}
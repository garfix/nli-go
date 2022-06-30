package knowledge

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type GenerationFunctionBase struct {
	KnowledgeBaseCore
	matcher *central.RelationMatcher
	state   *mentalese.GenerationState
	log     *common.SystemLog
}

func NewGenerationFunctionBase(name string, state *mentalese.GenerationState, log *common.SystemLog) *GenerationFunctionBase {
	return &GenerationFunctionBase{
		log:               log,
		KnowledgeBaseCore: KnowledgeBaseCore{name},
		state:             state,
		matcher:           central.NewRelationMatcher(log),
	}
}

func (base *GenerationFunctionBase) GetFunctions() map[string]api.SimpleFunction {
	return map[string]api.SimpleFunction{
		mentalese.PredicateAlreadyGenerated: base.alreadyGenerated,
	}
}

func (base *GenerationFunctionBase) alreadyGenerated(messenger api.SimpleMessenger, input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(input, "v", base.log) {
		return mentalese.NewBinding(), false
	}

	id := bound.Arguments[0]

	if !base.state.IsGenerated(id) {
		return mentalese.NewBinding(), false
	} else {
		return binding, true
	}
}

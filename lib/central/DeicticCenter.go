package central

import "nli-go/lib/mentalese"

type DeicticCenter struct {
	Binding mentalese.Binding
}

const DeixisTime = "time"

func NewDeicticCenter() *DeicticCenter {
	return &DeicticCenter{
		Binding: mentalese.NewBinding(),
	}
}

func (center *DeicticCenter) Initialize() {
	center.Binding = mentalese.NewBinding()
}

func (center *DeicticCenter) SetTime(time mentalese.RelationSet) {
	center.Binding.Set(DeixisTime, mentalese.NewTermRelationSet(time))
}

func (center *DeicticCenter) GetTime() mentalese.RelationSet {
	time, found := center.Binding.Get(DeixisTime)
	if found {
		return time.TermValueRelationSet
	} else {
		return mentalese.RelationSet{}
	}
}
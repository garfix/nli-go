package central

import "nli-go/lib/mentalese"

type DeicticCenter struct {
	Binding mentalese.Binding
}

const DeixisTime = "time"
const DeixisCenter = "center"

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

func (center *DeicticCenter) SetCenter(variable string) {
	center.Binding.Set(DeixisCenter, mentalese.NewTermVariable(variable))
}

func (center *DeicticCenter) GetCenter() string {
	c, found := center.Binding.Get(DeixisCenter)
	if found {
		return c.TermValue
	} else {
		return ""
	}
}

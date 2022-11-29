package mentalese

type DeicticCenter struct {
	Binding Binding
}

const DeixisTime = "time"
const DeixisCenter = "center"

func NewDeicticCenter() *DeicticCenter {
	return &DeicticCenter{
		Binding: NewBinding(),
	}
}

func (center *DeicticCenter) Initialize() {
	center.Binding = NewBinding()
}

func (center *DeicticCenter) Copy() *DeicticCenter {
	return &DeicticCenter{
		Binding: center.Binding.Copy(),
	}
}

func (center *DeicticCenter) SetTime(time RelationSet) {
	center.Binding.Set(DeixisTime, NewTermRelationSet(time))
}

func (center *DeicticCenter) GetTime() RelationSet {
	time, found := center.Binding.Get(DeixisTime)
	if found {
		return time.TermValueRelationSet
	} else {
		return RelationSet{}
	}
}

func (center *DeicticCenter) SetCenter(variable string) {
	center.Binding.Set(DeixisCenter, NewTermVariable(variable))
}

func (center *DeicticCenter) GetCenter() string {
	c, found := center.Binding.Get(DeixisCenter)
	if found {
		return c.TermValue
	} else {
		return ""
	}
}

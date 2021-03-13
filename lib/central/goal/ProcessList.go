package goal

import "nli-go/lib/mentalese"

type ProcessList struct {
	list []*Process
}

func NewProcessList() *ProcessList {
	return &ProcessList{
		list: []*Process{},
	}
}

func (p *ProcessList) Initialize() {
	p.list = []*Process{}
}

func (p *ProcessList) GetProcesses() []*Process {
	return p.list
}

func (p *ProcessList) GetProcess(goalId string) *Process {
	for _, process := range p.list {
		if process.GoalId == goalId {
			return process
		}
	}

	return nil
}

func (p *ProcessList) GetOrCreateProcess(goalId string, goalSet mentalese.RelationSet) *Process {
	for _, process := range p.list {
		if process.GoalId == goalId {
			return process
		}
	}

	process := NewProcess(goalId, goalSet, mentalese.InitBindingSet(mentalese.NewBinding()))
	p.list = append(p.list, process)

	return process
}

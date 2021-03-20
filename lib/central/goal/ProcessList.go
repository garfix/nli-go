package goal

import "nli-go/lib/mentalese"

type ProcessList struct {
	List []*Process
}

func NewProcessList() *ProcessList {
	return &ProcessList{
		List: []*Process{},
	}
}

func (p *ProcessList) Initialize() {
	p.List = []*Process{}
}

func (p *ProcessList) GetProcesses() []*Process {
	return p.List
}

func (p *ProcessList) GetProcess(goalId string) *Process {
	for _, process := range p.List {
		if process.GoalId == goalId {
			return process
		}
	}

	return nil
}

func (p *ProcessList) GetOrCreateProcess(goalId string, goalSet mentalese.RelationSet) *Process {
	for _, process := range p.List {
		if process.GoalId == goalId {
			return process
		}
	}

	process := NewProcess(goalId, goalSet, mentalese.InitBindingSet(mentalese.NewBinding()))
	p.List = append(p.List, process)

	return process
}

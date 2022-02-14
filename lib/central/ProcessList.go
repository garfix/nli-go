package central

import "nli-go/lib/mentalese"

type ProcessList struct {
	List           []*Process
	messageManager *MessageManager
}

func NewProcessList(messageManager *MessageManager) *ProcessList {
	return &ProcessList{
		List:           []*Process{},
		messageManager: messageManager,
	}
}

func (p *ProcessList) IsEmpty() bool {
	return len(p.List) == 0
}

func (p *ProcessList) Add(process *Process) {
	p.List = append(p.List, process)
}

func (p *ProcessList) Remove(process *Process) {
	newList := []*Process{}

	for _, item := range p.List {
		if item != process {
			newList = append(newList, item)
		}
	}

	p.List = newList
}

func (p *ProcessList) AddListener(l MessageListener) []mentalese.RelationSet {
	return p.messageManager.AddListener(l)
}

func (p *ProcessList) RemoveListener(l MessageListener) {
	p.messageManager.RemoveListener(l)
}

func (p *ProcessList) NotifyListeners(message mentalese.RelationSet) {
	p.messageManager.NotifyListeners(message)
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

func (p *ProcessList) RemoveProcess(goalId string) {

	newList := []*Process{}

	for _, process := range p.List {
		if process.GoalId != goalId {
			newList = append(newList, process)
		}
	}

	p.List = newList
}

func (p *ProcessList) GetOrCreateProcess(goalId string, goalSet mentalese.RelationSet) *Process {
	for _, process := range p.List {
		if process.GoalId == goalId {
			return process
		}
	}

	process := p.CreateProcess(goalId, goalSet, mentalese.InitBindingSet(mentalese.NewBinding()))
	p.List = append(p.List, process)

	return process
}

func (p *ProcessList) CreateProcess(goalId string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return NewProcess(goalId, goalSet, bindings)
}

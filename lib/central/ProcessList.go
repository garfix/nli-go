package central

import "nli-go/lib/mentalese"

const LANGUAGE_PROCESS = "language process"
const ROBOT_PROCESS = "robot process"
const SIMPLE_PROCESS = "simple process"

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

func (p *ProcessList) IsProcessTypeActive(processType string) bool {
	if processType == SIMPLE_PROCESS {
		return false
	}

	for _, p := range p.List {
		// each process type can occur only once
		if p.ProcessType == processType {
			return true
		}
	}

	return false
}

func (p *ProcessList) IsEmpty() bool {
	return len(p.List) == 0
}

func (p *ProcessList) Add(process *Process) bool {
	if p.IsProcessTypeActive(process.ProcessType) {
		return false
	}

	p.List = append(p.List, process)
	return true
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

func (p *ProcessList) GetOrCreateProcess(processType string, goalId string, goalSet mentalese.RelationSet) *Process {
	for _, process := range p.List {
		if process.GoalId == goalId {
			return process
		}
	}

	process := p.CreateProcess(processType, goalId, goalSet, mentalese.InitBindingSet(mentalese.NewBinding()))
	p.List = append(p.List, process)

	return process
}

func (p *ProcessList) CreateProcess(processType string, goalId string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return NewProcess(processType, goalId, goalSet, bindings)
}

package central

import "nli-go/lib/mentalese"

const LANGUAGE_PROCESS = "language"
const ROBOT_PROCESS = "robot"
const SIMPLE_PROCESS = "simple"

type MessageListener func()

type ProcessList struct {
	List      []*Process
	listeners []MessageListener
}

func NewProcessList() *ProcessList {
	return &ProcessList{
		List: []*Process{},
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

	if len(p.List) == 0 {
		p.NotifyListeners()
	}
}

func (p *ProcessList) AddListener(l MessageListener) {
	p.listeners = append(p.listeners, l)
}

func (p *ProcessList) NotifyListeners() {
	for _, l := range p.listeners {
		l()
	}
}

func (p *ProcessList) Initialize() {
	p.List = []*Process{}
}

func (p *ProcessList) GetProcesses() []*Process {
	return p.List
}

func (p *ProcessList) GetProcessByType(processType string) *Process {
	for _, process := range p.List {
		if process.ProcessType == processType {
			return process
		}
	}

	return nil
}

func (p *ProcessList) CreateProcess(processType string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return NewProcess(processType, goalSet, bindings)
}

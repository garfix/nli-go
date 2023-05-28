package central

const RESOURCE_LANGUAGE = "language"
const RESOURCE_ROBOT = "robot"
const NO_RESOURCE = "no-resource"

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
	if processType == NO_RESOURCE {
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

func (p *ProcessList) GetProcessByType(processType string) *Process {
	for _, process := range p.List {
		if process.ProcessType == processType {
			return process
		}
	}

	return nil
}

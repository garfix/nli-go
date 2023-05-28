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

func (p *ProcessList) IsResourceActive(resource string) bool {
	if resource == NO_RESOURCE {
		return false
	}

	for _, p := range p.List {
		// each process type can occur only once
		if p.resource == resource {
			return true
		}
	}

	return false
}

func (p *ProcessList) IsEmpty() bool {
	return len(p.List) == 0
}

func (p *ProcessList) Add(process *Process) bool {
	if p.IsResourceActive(process.resource) {
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

func (p *ProcessList) GetProcessByResource(resource string) *Process {
	for _, process := range p.List {
		if process.resource == resource {
			return process
		}
	}

	return nil
}

package central

import (
	"nli-go/lib/mentalese"
)

type Process struct {
	GoalId     string
	Stack      []*StackFrame
	Slots      map[string]mentalese.Term
	WaitingFor mentalese.RelationSet
}

func NewProcess(goalId string, goalSet mentalese.RelationSet, bindings mentalese.BindingSet) *Process {
	return &Process{
		GoalId: goalId,
		Stack: []*StackFrame{
			NewStackFrame(goalSet, bindings),
		},
		Slots:      map[string]mentalese.Term{},
		WaitingFor: nil,
	}
}

func (p *Process) SetWaitingFor(set mentalese.RelationSet) {
	p.WaitingFor = set
}

func (p *Process) GetWaitingFor() mentalese.RelationSet {
	return p.WaitingFor
}

func (p *Process) PushFrame(frame *StackFrame) {
	p.Stack = append(p.Stack, frame)
}

func (p *Process) Clear() {
	p.Stack = []*StackFrame{}
}

func (p *Process) EmptyRelationCheck() {
	frame := p.GetLastFrame()
	if frame.Relations.IsEmpty() {
		p.advanceFrame(frame)
	}
}

func (p *Process) Advance() {

	frame := p.GetLastFrame()
	frame.Cursor = NewStackFrameCursor()

	p.AdvanceHandler()
}

func (p *Process) AdvanceHandler() {

	frame := p.GetLastFrame()
	frame.HandlerIndex++

	if frame.HandlerIndex >= frame.HandlerCount {
		p.AdvanceBinding()
	}
}

func (p *Process) AdvanceBinding() {

	frame := p.GetLastFrame()
	frame.HandlerIndex = 0
	frame.InBindingIndex++

	if frame.InBindingIndex >= frame.InBindings.GetLength() {
		p.advanceRelation(frame)
	}
}

func (p *Process) advanceRelation(frame *StackFrame) {

	frame.InBindings = frame.OutBindings
	frame.InBindingIndex = 0

	frame.OutBindings = mentalese.NewBindingSet()

	frame.RelationIndex++

	if frame.InBindings.IsEmpty() {
		// process failed due to no result bindings
		p.advanceFrame(frame)
	} else if frame.IsDone() {
		p.advanceFrame(frame)
	}
}

func (p *Process) advanceFrame(frame *StackFrame) {

	p.PopFrame()

	// transfer child bindings to parent
	resultBindings := frame.InBindings
	newLastFrame := p.GetLastFrame()
	if newLastFrame != nil {
		newLastFrame.Cursor.ChildFrameResultBindings.AddMultiple(resultBindings)
	}
}

// prepare the active binding to be fed to a function
func (p *Process) GetPreparedBinding(f *StackFrame) mentalese.Binding {

	binding := f.GetCurrentInBinding()
	relation := f.GetCurrentRelation()

	binding = p.addMutableVariables(relation, binding)

	// filter out only the variables needed by the relation
	binding = binding.FilterVariablesByName(relation.GetVariableNames())

	return binding
}

func (p *Process) addMutableVariables(relation mentalese.Relation, binding mentalese.Binding) mentalese.Binding {
	scope := p.GetContextScope()
	if scope != nil {
		mutables := scope.Cursor.MutableVariableValues
		filtered := mutables.FilterVariablesByName(relation.GetVariableNames())
		return binding.Merge(filtered)
	} else {
		return binding
	}
}

func (p *Process) AddMutableVariablesMultiple(relation mentalese.Relation, bindings mentalese.BindingSet) mentalese.BindingSet {
	scope := p.GetContextScope()
	if scope != nil {
		mutables := scope.Cursor.MutableVariableValues
		filtered := mutables.FilterVariablesByName(relation.GetVariableNames())
		newBindings := mentalese.NewBindingSet()
		for _, binding := range bindings.GetAll() {
			newBindings.Add(binding.Merge(filtered))
		}
		return newBindings
	} else {
		return bindings
	}
}

func (p *Process) CreateMessenger(processRunner *ProcessRunner, process *Process) *Messenger {
	frame := p.GetLastFrame()

	return NewMessenger(processRunner, process, frame.Cursor, p.Slots)
}

func (p *Process) ProcessMessenger(messenger *Messenger, currentFrame *StackFrame) *StackFrame {

	outBindings := messenger.GetOutBindings()

	for slot, value := range messenger.newSlots {
		p.Slots[slot] = value
	}

	// assign the outbound mutable variables
	//relationVariables1 := currentFrame.GetCurrentRelation().GetVariableNames()
	//mutBindings := outBindings.FilterMutableVariables()
	//for _, v := range relationVariables1 {
	//	val, found := mutBindings.Get(v)
	//	if found {
	//		p.GetCurrentScope().Cursor.MutableVariableValues.Set(v, val)
	//	}
	//}

	outBindingsWithoutMutables := outBindings.RemoveMutableVariables()

	processedOutBindings := mentalese.NewBindingSet()
	for _, outBinding := range outBindingsWithoutMutables.GetAll() {
		relationVariables := currentFrame.GetCurrentRelation().GetVariableNames()

		// filter out temporary variables
		cleanBinding := outBinding.FilterVariablesByName(relationVariables)
		// make sure the original values are present
		cleanBinding = cleanBinding.Merge(currentFrame.GetCurrentInBinding())
		processedOutBindings.Add(cleanBinding)
	}

	outBindings2 := mentalese.NewBindingSet()
	currentFrame, outBindings2 = p.executeProcessInstructions(messenger, currentFrame, processedOutBindings)

	currentFrame.AddOutBindings(currentFrame.GetCurrentInBinding(), outBindings2)

	//if messenger.GetChildFrame() != nil {
	//	p.PushFrame(messenger.GetChildFrame())
	//}

	return currentFrame
}

func (p *Process) executeProcessInstructions(messenger *Messenger, currentFrame *StackFrame, outBindings mentalese.BindingSet) (*StackFrame, mentalese.BindingSet) {

	for instruction, _ := range messenger.GetProcessInstructions() {
		switch instruction {
		//case mentalese.ProcessInstructionLet:
		//	p.AddMutableVariable(value)
		case mentalese.ProcessInstructionBreak:
			currentFrame = p.executeBreak(currentFrame, outBindings, false)
			outBindings = mentalese.NewBindingSet() //currentFrame.InBindings
		case mentalese.ProcessInstructionCancel:
			currentFrame = p.executeBreak(currentFrame, outBindings, true)
			outBindings = mentalese.NewBindingSet()
		case mentalese.ProcessInstructionReturn:
			currentFrame = p.executeReturn(currentFrame, outBindings)
			outBindings = mentalese.NewBindingSet() //currentFrame.InBindings
		}
	}

	return currentFrame, outBindings
}

func (p *Process) SetMutableVariable(variable string, value mentalese.Term) {
	scope := p.GetCurrentScope()
	if scope != nil {
		scope.Cursor.MutableVariableValues.Set(variable, value)
	}
}

func (p *Process) GetCurrentScope() *StackFrame {

	var scope *StackFrame = nil
	i := len(p.Stack) - 1

	for i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			return frame
		}

		i--
	}

	return scope
}

func (p *Process) GetContextScope() *StackFrame {

	var scope *StackFrame = nil
	i := len(p.Stack) - 2

	for i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			return frame
		}

		i--
	}

	return scope
}

func (p *Process) executeBreak(currentFrame *StackFrame, bindings mentalese.BindingSet, cancel bool) *StackFrame {
	done := false
	i := len(p.Stack) - 1

	for !done && i >= 0 {

		frame := p.Stack[i]
		frameType := frame.Cursor.GetType()

		if frameType == mentalese.FrameTypeLoop {
			frame.Cursor.SetPhase(PhaseInterrupted)
			if cancel {
				//frame.Cursor.SetPhase(PhaseCanceled)
			} else {
				frame.Cursor.ChildFrameResultBindings.AddMultiple(bindings)
				//frame.Cursor.SetPhase(PhaseBreaked)
			}
			done = true
		} else {
			frame.Cursor.SetPhase(PhaseInterrupted)
		}

		i--
	}

	return currentFrame
}

func (p *Process) executeReturn(currentFrame *StackFrame, bindings mentalese.BindingSet) *StackFrame {
	done := false
	i := len(p.Stack) - 1

	for !done && i >= 0 {
		frame := p.Stack[i]

		if frame.Cursor.GetType() == mentalese.FrameTypeScope {
			//frame.Cursor.SetPhase(PhaseInterrupted)
			frame.Cursor.ChildFrameResultBindings.AddMultiple(bindings)
			done = true
		} else {
			frame.Cursor.SetPhase(PhaseInterrupted)
		}

		i--
	}

	return currentFrame
}

func (p *Process) ProcessMessengerMultipleBindings(messenger *Messenger, frame *StackFrame) {

	outBindings := messenger.GetOutBindings()
	outBindingsWithoutMutables := outBindings.RemoveMutableVariables()

	// add bindings without variable validation
	frame.OutBindings.AddMultiple(outBindingsWithoutMutables)

	// skip the bindings
	frame.InBindingIndex = frame.InBindings.GetLength() - 1
}

func (p *Process) GetCursor() *StackFrameCursor {
	frame := p.GetLastFrame()
	if frame == nil {
		return nil
	}
	return frame.Cursor
}

func (p *Process) GetLastFrame() *StackFrame {
	if len(p.Stack) == 0 {
		return nil
	} else {
		frame := p.Stack[len(p.Stack)-1]
		return frame
	}
}

func (p *Process) GetBeforeLastFrame() *StackFrame {
	if len(p.Stack) < 2 {
		return nil
	} else {
		return p.Stack[len(p.Stack)-2]
	}
}

func (p *Process) PopFrame() {
	p.Stack = p.Stack[0 : len(p.Stack)-1]
}

func (p *Process) IsDone() bool {
	return len(p.Stack) == 0
}

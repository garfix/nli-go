package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
)

type ProcessRunner struct {
	solver *ProblemSolver
	log    *common.SystemLog
	list   *ProcessList
}

func NewProcessRunner(list *ProcessList, solver *ProblemSolver, log *common.SystemLog) *ProcessRunner {
	return &ProcessRunner{
		solver: solver,
		log:    log,
		list:   list,
	}
}

func (p *ProcessRunner) StartProcess(processType string, relationSet mentalese.RelationSet, binding mentalese.Binding) bool {
	process := p.list.CreateProcess(processType, "", relationSet, mentalese.InitBindingSet(binding))
	if p.list.IsProcessTypeActive(processType) {
		return false
	}
	go p.StartProcessNow(process)
	return true
}

func (p *ProcessRunner) StartProcessNow(process *Process) {
	p.list.Add(process)
	p.RunProcessLevel(process, 0)
	p.list.Remove(process)
	// notify listeners we're done
	p.list.NotifyListeners(mentalese.RelationSet{})
}

func (p *ProcessRunner) SendMessage(message mentalese.Request) {
	process := p.list.GetProcessByType(message.ProcessType)
	if process != nil {
		process.GetChannel() <- message
	}
}

func (p *ProcessRunner) RunRelationSet(processType string, relationSet mentalese.RelationSet) mentalese.BindingSet {
	bindings := mentalese.InitBindingSet(mentalese.NewBinding())
	return p.RunRelationSetWithBindings(processType, relationSet, bindings)
}

func (p *ProcessRunner) PushAndRun(process *Process, relations mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {
	level := len(process.Stack)
	process.PushFrame(NewStackFrame(relations, bindings))
	return p.RunProcessLevel(process, level)
}

func (p *ProcessRunner) RunRelationSetWithBindings(processType string, relationSet mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {
	process := p.list.CreateProcess(processType, "", relationSet, bindings)
	frame := process.Stack[0]
	p.RunProcessLevel(process, 0)
	// note: frame has already been deleted; frame is now just the last reference
	return frame.InBindings
}

func (p *ProcessRunner) RunProcessLevel(process *Process, level int) mentalese.BindingSet {
	for len(process.Stack) > level {
		hasStopped := p.step(process)
		if hasStopped {
			return mentalese.NewBindingSet()
		}
	}

	if level == 0 {
		return mentalese.NewBindingSet()
	} else {

		resultBindings := process.Stack[len(process.Stack)-1].Cursor.ChildFrameResultBindings
		process.Stack[len(process.Stack)-1].Cursor.ChildFrameResultBindings = mentalese.NewBindingSet()
		return resultBindings
	}
}

func (p *ProcessRunner) step(process *Process) bool {
	currentFrame := process.GetLastFrame()
	hasStopped := false

	debug := p.before(process, currentFrame, len(process.Stack))
	p.log.AddDebug("frame", debug)

	messenger := process.CreateMessenger(p, process)
	relation := currentFrame.GetCurrentRelation()
	outBindings := mentalese.NewBindingSet()

	_, found := p.solver.multiBindingFunctions[relation.Predicate]
	if found {

		preparedBindings := process.AddMutableVariablesMultiple(relation, currentFrame.InBindings)
		outBindings, _ = p.solver.SolveMultipleBindings(messenger, relation, preparedBindings)
		messenger.AddOutBindings(outBindings)
		process.ProcessMessengerMultipleBindings(messenger, currentFrame)

	} else {

		preparedBinding := process.GetPreparedBinding(currentFrame)

		handler := p.PrepareHandler(relation, currentFrame, process)
		if handler == nil {
			return true
		} else {
			preparedRelation := p.evaluateArguments(process, relation, preparedBinding)
			outBindings = handler(messenger, preparedRelation, preparedBinding)
			messenger.AddOutBindings(outBindings)
			currentFrame = process.ProcessMessenger(messenger, currentFrame)
		}
	}

	debug = p.after(process, currentFrame, outBindings, len(process.Stack))
	p.log.AddDebug("frame", debug)

	if messenger.GetCursor().GetState() == StateInterrupted {

		debug = p.breaked(len(process.Stack))
		p.log.AddDebug("frame", debug)
	}

	if currentFrame == process.GetLastFrame() {
		process.Advance()
	} else {
		process.EmptyRelationCheck()
	}

	return hasStopped
}

func (p *ProcessRunner) removeMutableVariables(bindings mentalese.BindingSet) mentalese.BindingSet {
	return bindings.RemoveMutableVariables()
}

func (p *ProcessRunner) evaluateArguments(process *Process, relation mentalese.Relation, binding mentalese.Binding) mentalese.Relation {
	newRelation := relation

	for i, argument := range relation.Arguments {
		if argument.IsRelationSet() && len(argument.TermValueRelationSet) == 1 {
			firstRelation := argument.TermValueRelationSet[0]
			for j, arg := range firstRelation.Arguments {
				if arg.IsAtom() && arg.TermValue == mentalese.AtomReturnValue {
					newRelation = newRelation.Copy()
					newRelation.Arguments[i] = p.evaluateFunction(process, firstRelation, j, binding)
					break
				}
			}
		}
	}

	return newRelation
}

func (p *ProcessRunner) evaluateFunction(process *Process, relation mentalese.Relation, returnVariableIndex int, binding mentalese.Binding) mentalese.Term {
	variable := p.solver.variableGenerator.GenerateVariable("ReturnVal")
	newRelation := relation.Copy()
	newRelation.Arguments[returnVariableIndex] = variable
	resultBindings := p.PushAndRun(process, mentalese.RelationSet{newRelation}, mentalese.InitBindingSet(binding))
	if resultBindings.GetLength() == 0 {
		return mentalese.NewTermAtom(mentalese.AtomNone)
	} else {
		returnValue, found := resultBindings.Get(0).Get(variable.TermValue)
		if found {
			return returnValue
		} else {
			return mentalese.NewTermAtom(mentalese.AtomNone)
		}
	}
}

func (p *ProcessRunner) PrepareHandler(relation mentalese.Relation, frame *StackFrame, process *Process) api.RelationHandler {

	handlers := p.solver.GetHandlers(relation)

	frame.HandlerCount = len(handlers)

	if frame.HandlerIndex >= len(handlers) {
		// there may just be no handlers, or handlers could have been removed from the knowledge bases
		if frame.HandlerIndex == 0 {
			p.log.AddError("Predicate not supported by any knowledge base: " + relation.String())
		}
		return nil
	}

	return handlers[frame.HandlerIndex]
}

func (p *ProcessRunner) before(process *Process, frame *StackFrame, stackDepth int) string {

	padding := strings.Repeat("  ", stackDepth)

	child := ""
	childBindings := frame.Cursor.ChildFrameResultBindings
	if !childBindings.IsEmpty() {
		child = " from child: " + childBindings.String()
	}

	prepared := ""
	if frame.InBindings.GetLength() > 0 {
		prepared = process.GetPreparedBinding(frame).String()
	}

	fromChild := ""
	handlerIndex := strconv.Itoa(frame.HandlerIndex)

	text := fromChild + frame.Relations[frame.RelationIndex].String() + ":" + handlerIndex + "  " + prepared + " " + child
	return padding + text
}

func (p *ProcessRunner) after(process *Process, frame *StackFrame, bindings mentalese.BindingSet, stackDepth int) string {
	padding := strings.Repeat("  ", stackDepth)
	debug := padding + "╰ " + bindings.String()
	if process.GetLastFrame() != frame {
		debug = ""
	}
	return debug
}

func (p *ProcessRunner) breaked(stackDepth int) string {
	padding := strings.Repeat("  ", stackDepth)
	return padding + "* breaked"
}

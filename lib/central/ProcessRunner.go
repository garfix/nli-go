package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
	"strings"
)

type ProcessRunner struct {
	solver *ProblemSolverAsync
	log  *common.SystemLog
	list ProcessList
}

func NewProcessRunner(solver *ProblemSolverAsync, log *common.SystemLog) *ProcessRunner {
	return &ProcessRunner{
		solver: solver,
		log:    log,
	}
}

func (p *ProcessRunner) RunRelationSet(relationSet mentalese.RelationSet) mentalese.BindingSet {
	bindings := mentalese.InitBindingSet(mentalese.NewBinding())
	return p.RunRelationSetWithBindings(relationSet, bindings)
}

func (p *ProcessRunner) RunRelationSetWithBindings(relationSet mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.BindingSet {
	process := NewProcess("", relationSet, bindings)
	frame := process.Stack[0]
	p.RunProcess(process)
	// note: frame has already been deleted; frame is now just the last reference
	return frame.InBindings
}

func (p *ProcessRunner) RunProcess(process *Process) {
	for !process.IsDone() {
		hasStopped := p.step(process)
		if hasStopped {
			break
		}
	}
}

func (p *ProcessRunner) step(process *Process) bool {
	currentFrame := process.GetLastFrame()
	hasStopped := false

	debug := p.before(process, currentFrame, len(process.Stack))

	messenger := process.CreateMessenger(p)
	relation := currentFrame.GetCurrentRelation()

	_, found := p.solver.multiBindingFunctions[relation.Predicate]
	if found {

		preparedBindings := currentFrame.InBindings
		outBindings, _ := p.solver.SolveMultipleBindings(messenger, relation, preparedBindings)
		messenger.AddOutBindings(outBindings)
		process.ProcessMessengerMultipleBindings(messenger, currentFrame)

	} else {

		preparedBinding := process.GetPreparedBinding(currentFrame)

		handler := p.PrepareHandler(relation, currentFrame, process)
		if handler == nil {
			return true
		} else {
			preparedRelation := p.evaluateArguments(relation, preparedBinding)

			//if currentFrame.callback {
			//	currentFrame.callback()
			//} else {

			outBindings := handler(messenger, preparedRelation, preparedBinding)

			//}
			//
			//if outBindings == ignoreThis {
			//	push new frames
			//	currentFrameOfZo.callback = messenger.callback
			//	return true
			//}

			messenger.AddOutBindings(outBindings)
			currentFrame, hasStopped = process.ProcessMessenger(messenger, currentFrame)
		}
	}

	debug += p.after(process, currentFrame)
	p.log.AddDebug("frame", debug)

	if currentFrame == process.GetLastFrame() {
		process.Advance()
	} else {
		process.EmptyRelationCheck()
	}

	return hasStopped
}

func (p *ProcessRunner) evaluateArguments(relation mentalese.Relation, binding mentalese.Binding) mentalese.Relation {
	newRelation := relation

	for i, argument := range relation.Arguments {
		if argument.IsRelationSet() && len(argument.TermValueRelationSet) == 1 {
			firstRelation := argument.TermValueRelationSet[0]
			for j, arg := range firstRelation.Arguments {
				if arg.IsAtom() && arg.TermValue == mentalese.AtomReturnValue {
					newRelation = newRelation.Copy()
					newRelation.Arguments[i] = p.evaluateFunction(firstRelation, j, binding)
					break
				}
			}
		}
	}

	return newRelation
}

func (p *ProcessRunner) evaluateFunction(relation mentalese.Relation, returnVariableIndex int, binding mentalese.Binding) mentalese.Term {
	variable := p.solver.variableGenerator.GenerateVariable("ReturnVal")
	newRelation := relation.Copy()
	newRelation.Arguments[returnVariableIndex] = variable
	resultBindings := p.RunRelationSetWithBindings(mentalese.RelationSet{newRelation}, mentalese.InitBindingSet(binding))
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
			p.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
			process.Clear()
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
	if len(frame.Cursor.State) > 0 {
		fromChild = "╰ "
	}

	handlerIndex := strconv.Itoa(frame.HandlerIndex)

	text := fromChild + frame.Relations[frame.RelationIndex].String() + ":" + handlerIndex + "  " + prepared + " " + child
	return padding + text
}

func (p *ProcessRunner) after(process *Process, frame *StackFrame) string {
	debug := ": " + frame.OutBindings.String()
	if process.GetLastFrame() != frame {
		debug = ""
	}
	return debug
}

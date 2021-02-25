package central

import (
	"nli-go/lib/api"
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strings"
)

type ProcessRunner struct {
	solver *ProblemSolverAsync
	log *common.SystemLog
	list goal.ProcessList
}

func NewProcessRunner(solver *ProblemSolverAsync, log *common.SystemLog) *ProcessRunner {
	return &ProcessRunner{
		solver: solver,
		log: log,
		list: goal.ProcessList{},
	}
}

func (p ProcessRunner) RunProcess(goalId int, goalSet mentalese.RelationSet) {
	process := p.list.GetProcess(goalId, goalSet)

	for !process.IsDone() {
		p.step(process)
	}
}

func (p ProcessRunner) step(process *goal.Process) {
	currentFrame := process.GetLastFrame()

	debug := p.before(process, currentFrame, len(process.Stack))

	messenger := process.CreateMessenger()
	relation := currentFrame.GetCurrentRelation()

	_, found := p.solver.solver.index.multiBindingFunctions[relation.Predicate]
	if found {

		preparedBindings := currentFrame.InBindings
		p.solver.SolveMultipleBindings(messenger, relation, preparedBindings)
		process.ProcessMessengerMultipleBindings(messenger, currentFrame)

	} else {

		if relation.Predicate == mentalese.PredicateLet {
			p.createMutableVariable(process, relation)
		}

		preparedBinding := process.GetPreparedBinding(currentFrame)
		//p.solver.SolveSingleRelationSingleBinding(messenger, relation, preparedBinding)

		handler := p.PrepareHandler(relation.Predicate, currentFrame)
		if handler == nil {
			// there may just be no handlers, or handlers could have been removed from the knowledge bases
			if currentFrame.HandlerIndex == 0 {
				p.log.AddError("Predicate not supported by any knowledge base: " + relation.Predicate)
				process.Clear()
				return
			}
		} else {
			outBindings := handler(messenger, relation, preparedBinding)
			messenger.AddOutBindings(outBindings)
			process.ProcessMessenger(messenger, currentFrame)
		}
	}

	debug += p.after(process, currentFrame)
	p.log.AddDebug("frame", debug)

	// if the relation has not pushed a new frame, then it is done processing
	if currentFrame == process.GetLastFrame() {
		process.Advance()
	} else {
		process.EmptyRelationCheck()
	}
}

func (p ProcessRunner) PrepareHandler(predicate string, frame *goal.StackFrame) api.RelationHandler {

	handlers := p.solver.GetHandlers(predicate)

	frame.HandlerCount = len(handlers)

	if frame.HandlerIndex >= len(handlers) {
		return nil
	}

	return handlers[frame.HandlerIndex]
}

func (p ProcessRunner) createMutableVariable(process *goal.Process, relation mentalese.Relation) {

	if len(relation.Arguments) != 2 {
		p.log.AddError("`let should have two arguments`")
		return
	}

	variableTerm := relation.Arguments[0]

	if !variableTerm.IsVariable() {
		p.log.AddError("First argument of `let` should be a variable")
		return
	}

	process.AddMutableVariable(variableTerm.TermValue)
}

func (p ProcessRunner) before(process *goal.Process, frame *goal.StackFrame, stackDepth int) string {

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

	text := frame.Relations[frame.RelationIndex].String() + "  " + prepared + " " + child
	return padding + text
}

func (p ProcessRunner) after(process *goal.Process, frame *goal.StackFrame) string {
	debug := ": " + frame.OutBindings.String()
	if process.GetLastFrame() != frame {
		debug = ""
	}
	return debug
}
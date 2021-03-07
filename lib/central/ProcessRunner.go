package central

import (
	"nli-go/lib/api"
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
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

func (p *ProcessRunner) GetProcesses() []*goal.Process {
	return p.list.GetProcesses()
}

func (p *ProcessRunner) GetProcessByGoalId(goalId string) *goal.Process {
	return p.list.GetProcess(goalId)
}

func (p *ProcessRunner) RunNow(goalSet mentalese.RelationSet) mentalese.BindingSet {
	process := goal.NewProcess("", goalSet)
	frame := process.Stack[0]
	p.runProcessNow(process)
	// note: frame has already been deleted; frame is now just the last reference
	return frame.InBindings
}

func (p *ProcessRunner) RunProcess(goalId string, goalSet mentalese.RelationSet) {
	process := p.list.GetOrCreateProcess(goalId, goalSet)
	p.runProcessNow(process)
}

func (p *ProcessRunner) runProcessNow(process *goal.Process) {
	for !process.IsDone() {
		hasStopped := p.step(process)
		if hasStopped {
			break
		}
	}
}

func (p *ProcessRunner) step(process *goal.Process) bool {
	currentFrame := process.GetLastFrame()
	hasStopped := false

	debug := p.before(process, currentFrame, len(process.Stack))

	messenger := process.CreateMessenger()
	relation := currentFrame.GetCurrentRelation()

	_, found := p.solver.index.multiBindingFunctions[relation.Predicate]
	if found {

		preparedBindings := currentFrame.InBindings
		outBindings, _ := p.solver.SolveMultipleBindings(messenger, relation, preparedBindings)
		messenger.AddOutBindings(outBindings)
		process.ProcessMessengerMultipleBindings(messenger, currentFrame)

	} else {

		preparedBinding := process.GetPreparedBinding(currentFrame)

		handler := p.PrepareHandler(relation, currentFrame, process)
		if handler == nil {
			return hasStopped
		} else {
			outBindings := handler(messenger, relation, preparedBinding)
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

func (p *ProcessRunner) PrepareHandler(relation mentalese.Relation, frame *goal.StackFrame, process *goal.Process) api.RelationHandler {

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

func (p *ProcessRunner) before(process *goal.Process, frame *goal.StackFrame, stackDepth int) string {

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

func (p *ProcessRunner) after(process *goal.Process, frame *goal.StackFrame) string {
	debug := ": " + frame.OutBindings.String()
	if process.GetLastFrame() != frame {
		debug = ""
	}
	return debug
}
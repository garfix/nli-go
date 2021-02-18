package central

import (
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

	p.debug(currentFrame, len(process.Stack))

	messenger := process.CreateMessenger()
	relation := currentFrame.GetCurrentRelation()

	_, found := p.solver.solver.index.multiBindingFunctions[relation.Predicate]
	if found {

		preparedBindings := currentFrame.InBindings
		p.solver.SolveMultipleBindings(messenger, relation, preparedBindings)
		process.ProcessMessengerMultipleBindings(messenger, currentFrame)

	} else {

		preparedBinding := currentFrame.GetPreparedBinding()
		p.solver.SolveSingleRelationSingleBinding(messenger, relation, preparedBinding)

		process.ProcessMessenger(messenger, currentFrame)

	}

	// if the relation has not pushed a new frame, then it is done processing
	if currentFrame == process.GetLastFrame() {
		process.Advance()
	} else {
		process.EmptyRelationCheck()
	}
}

func (p ProcessRunner) debug(frame *goal.StackFrame, stackDepth int) {

	padding := strings.Repeat("  ", stackDepth)

	child := ""
	childBindings := frame.Cursor.ChildFrameResultBindings
	if !childBindings.IsEmpty() {
		child = " from child: " + childBindings.String()
	}

	prepared := ""
	if frame.InBindings.GetLength() > 0 {
		prepared = frame.GetPreparedBinding().String()
	}

	text := frame.Relations[frame.RelationIndex].String() + "  " + prepared + " " + child
	p.log.AddDebug("frame",
		padding + text)
}
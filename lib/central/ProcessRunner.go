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
		p.singleStep(process)
	}
}

// cursor must always point to a real relation!
func (p ProcessRunner) singleStep(process *goal.Process) {
	currentFrame := process.GetLastFrame()

	p.debug(currentFrame, len(process.Stack))

// todo call functions that act on multiple bindings
// since these are only simple functions, this can be done inline without much fuzz

	// execute the relation at the cursor
	p.solver.SolveSingleRelationSingleBinding(process)

	// if the relation has not pushed a new frame, then it is done processing
	if currentFrame == process.GetLastFrame() {
		process.Advance()
	}
}

func (p ProcessRunner) debug(frame *goal.StackFrame, stackDepth int) {

	padding := strings.Repeat("  ", stackDepth)
	p.log.AddDebug("frame",
		padding + frame.Relations[frame.RelationIndex].String() + "  " + frame.InBindings.Get(frame.InBindingIndex).String())
}
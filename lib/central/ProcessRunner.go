package central

import (
	"nli-go/lib/central/goal"
	"nli-go/lib/mentalese"
)

type ProcessRunner struct {
	solver *ProblemSolverAsync
	list goal.ProcessList
}

func NewProcessRunner(solver *ProblemSolverAsync) *ProcessRunner {
	return &ProcessRunner{
		solver: solver,
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

// todo call functions that act on multiple bindings

	// execute the relation at the cursor
	p.solver.SolveSingleRelationSingleBinding(process)

	// if the relation has not pushed a new frame, then it is done processing
	if currentFrame == process.GetLastFrame() {

		// quit if there are no bindings
		if currentFrame.Bindings.IsEmpty() {
			process.Clear()
		} else {
			process.Advance()
		}
	}
}
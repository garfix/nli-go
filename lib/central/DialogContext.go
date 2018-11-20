package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

const predicateOpenQuestion = "open_question"
const predicateAnswerToOpenQuestion = "answer_open_question"
const predicateOriginalInput = "original_input"

// The dialog context stores questions and answers that involve interaction with the user while solving his/her main question
// It may also be used to data relations that may be needed in the next call of the library (within the same session)
type DialogContext struct {
	factBase *knowledge.InMemoryFactBase
	solver *ProblemSolver
	values mentalese.RelationSet
}

func NewDialogContext(matcher *mentalese.RelationMatcher, solver *ProblemSolver, log *common.SystemLog) *DialogContext {

	factBase := knowledge.NewInMemoryFactBase(
		"in-memory",
		mentalese.RelationSet{},
		matcher,
		[]mentalese.RelationTransformation{},
		mentalese.DbStats{},
		log,
	)

	return &DialogContext{
		factBase: factBase,
		solver: solver,
	}
}

func (dc *DialogContext) Initialize(values mentalese.RelationSet) {
	dc.factBase.SetRelations(values)
}

func (dc *DialogContext) SetOriginalInput(originalInput string) {
	dc.factBase.AddRelation(mentalese.NewRelation(predicateOriginalInput, []mentalese.Term{
		mentalese.NewString(originalInput),
	}))
}

func (dc *DialogContext) GetOriginalInput() (string, bool) {
	results := dc.factBase.MatchRelationToDatabase(mentalese.NewRelation(predicateOriginalInput, []mentalese.Term{
		mentalese.NewVariable("A"),
	}))

	if len(results) > 0 {
		result := results[0]
		return result["A"].TermValue, true
	} else {
		return "", false
	}
}

func (dc *DialogContext) AddRelation(relation mentalese.Relation) {
	dc.factBase.AddRelation(relation)
}

func (dc *DialogContext) FindRelations(relationset mentalese.RelationSet) []mentalese.Binding {
	return dc.solver.FindFacts(dc.factBase, relationset)
}

func (dc *DialogContext) GetRelations() mentalese.RelationSet {
	return dc.factBase.GetRelations()
}

func (dc *DialogContext) SetOpenQuestion(question string) {
	dc.factBase.AddRelation(mentalese.NewRelation(predicateOpenQuestion, []mentalese.Term{
		mentalese.NewString(question),
	}))
}

func (dc *DialogContext) GetOpenQuestion() (string, bool) {
	results := dc.factBase.MatchRelationToDatabase(mentalese.NewRelation(predicateOpenQuestion, []mentalese.Term{
		mentalese.NewVariable("Q"),
	}))

	if len(results) > 0 {
		result := results[0]
		return result["Q"].TermValue, true
	} else {
		return "", false
	}
}

func (dc* DialogContext) RemoveOpenQuestion() {
	dc.factBase.RemoveRelation(mentalese.NewRelation(predicateOpenQuestion, []mentalese.Term{
		mentalese.NewAnonymousVariable(),
	}))
}

func (dc *DialogContext) SetAnswerToOpenQuestion(answer string) {
	dc.factBase.AddRelation(mentalese.NewRelation(predicateAnswerToOpenQuestion, []mentalese.Term{
		mentalese.NewString(answer),
	}))
}

func (dc *DialogContext) GetAnswerToOpenQuestion() (string, bool) {
	results := dc.factBase.MatchRelationToDatabase(mentalese.NewRelation(predicateAnswerToOpenQuestion, []mentalese.Term{
		mentalese.NewVariable("A"),
	}))

	if len(results) > 0 {
		result := results[0]
		return result["A"].TermValue, true
	} else {
		return "", false
	}
}

func (dc* DialogContext) RemoveAnswerToOpenQuestion() {
	dc.factBase.RemoveRelation(mentalese.NewRelation(predicateAnswerToOpenQuestion, []mentalese.Term{
		mentalese.NewAnonymousVariable(),
	}))
}

func (dc* DialogContext) Process(currentInput string) string {

	originalInput := ""

	_, found := dc.GetOpenQuestion()
	if found {

		// data user response in open question
		dc.SetAnswerToOpenQuestion(currentInput)

		// return to not expecting an answer
		dc.RemoveOpenQuestion()

		// continue with the user's original question
		originalInput, _ = dc.GetOriginalInput()

	} else {

		originalInput = currentInput

		// data original question
		dc.SetOriginalInput(currentInput)

	}

	return originalInput
}
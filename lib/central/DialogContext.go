package central

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
)

const predicateOption = "option"
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

	transformationString := "[" +
		"name_information(Name, Database_name, Entity_id) => name_information(Name, Database_name, Entity_id); " +
		"]"

	parser := importer.NewInternalGrammarParser()

	transformations := parser.CreateTransformations(transformationString)

	factBase := knowledge.NewInMemoryFactBase(
		"in-memory",
		mentalese.RelationSet{},
		matcher,
		transformations,
		mentalese.DbStats{},
		mentalese.Entities{},
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

func (dc* DialogContext) RemoveOriginalInput() {
	dc.factBase.RemoveRelation(mentalese.NewRelation(predicateOriginalInput, []mentalese.Term{
		mentalese.NewAnonymousVariable(),
	}))
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

func (dc *DialogContext) AddOption(option string) {
	dc.factBase.AddRelation(mentalese.NewRelation(predicateOption, []mentalese.Term{
		mentalese.NewString(option),
	}))
}

func (dc *DialogContext) GetOpenOptions() []string {
	results := dc.factBase.MatchRelationToDatabase(mentalese.NewRelation(predicateOption, []mentalese.Term{
		mentalese.NewVariable("Q"),
	}))

	var options []string

	for _, result := range results {
		options = append(options, result["Q"].TermValue)
	}

	return options
}

func (dc* DialogContext) RemoveOpenOptions() {
	dc.factBase.RemoveRelation(mentalese.NewRelation(predicateOption, []mentalese.Term{
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

	options := dc.GetOpenOptions()

	// are we expecting an answer from the user?
	if len(options) > 0 {

		// check if user response matches options
		match := false
		for _, option := range options {
			if option == currentInput {
				match = true
			}
		}

		if match {

			// data user response in open question
			dc.SetAnswerToOpenQuestion(currentInput)

			// continue with the user's original question
			originalInput, _ = dc.GetOriginalInput()

		} else {

			// the user gave a response that does not match our expectation
			// assume the user is posing a new question
			originalInput = currentInput
		}

		// stop expecting an answer
		dc.RemoveOpenOptions()

	} else {

		originalInput = currentInput

		// data original question
		dc.SetOriginalInput(currentInput)

	}

	return originalInput
}
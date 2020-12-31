package central

import "nli-go/lib/common"

const MaxSizeAnaphoraQueue = 10

// The dialog context stores questions and answers that involve interaction with the user while solving his/her main question
// It may also be used to data relations that may be needed in the next call of the library (within the same session)
type DialogContext struct {
	storage *common.FileStorage
	OriginalInput string
	AnswerToOpenQuestion string
	NameInformations []NameInformation
	Options []string
	AnaphoraQueue *AnaphoraQueue
}

func NewDialogContext(storage *common.FileStorage) *DialogContext {
	dialogContext := &DialogContext{
		storage: storage,
	}
	dialogContext.Initialize()

	if storage != nil {
		storage.Read(dialogContext)
	}

	return dialogContext
}

func (dc *DialogContext) Initialize() {
	dc.OriginalInput = ""
	dc.AnswerToOpenQuestion = ""
	dc.NameInformations = []NameInformation{}
	dc.Options = []string{}
	dc.AnaphoraQueue = &AnaphoraQueue{}
}

func (dc *DialogContext) Store() {
	if dc.storage != nil {
		dc.storage.Write(dc)
	}
}

func (dc *DialogContext) SetOriginalInput(originalInput string) {
	dc.OriginalInput = originalInput
}

func (dc *DialogContext) GetOriginalInput() (string, bool) {
	return dc.OriginalInput, dc.OriginalInput != ""
}

func (dc* DialogContext) RemoveOriginalInput() {
	dc.OriginalInput = ""
}

func (dc *DialogContext) AddNameInformations(nameInformations []NameInformation) {
	dc.NameInformations = append(dc.NameInformations, nameInformations...)
}

func (dc *DialogContext) GetNameInformations() []NameInformation {
	return dc.NameInformations
}

func (dc *DialogContext) AddOption(option string) {
	dc.Options = append(dc.Options, option)
}

func (dc *DialogContext) GetOpenOptions() []string {
	return dc.Options
}

func (dc* DialogContext) RemoveOpenOptions() {
	dc.Options = []string{}
}

func (dc *DialogContext) SetAnswerToOpenQuestion(answer string) {
	dc.AnswerToOpenQuestion = answer
}

func (dc *DialogContext) GetAnswerToOpenQuestion() (string, bool) {
	return dc.AnswerToOpenQuestion, dc.AnswerToOpenQuestion != ""
}

func (dc* DialogContext) RemoveAnswerToOpenQuestion() {
	dc.AnswerToOpenQuestion = ""
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
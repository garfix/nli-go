package central

import "nli-go/lib/common"

const MaxSizeAnaphoraQueue = 10

// The dialog context stores questions and answers that involve interaction with the user while solving his/her main question
// It may also be used to data relations that may be needed in the next call of the library (within the same session)
type DialogContext struct {
	storage *common.FileStorage
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
	dc.AnaphoraQueue = &AnaphoraQueue{}
}

func (dc *DialogContext) Store() {
	if dc.storage != nil {
		dc.storage.Write(dc)
	}
}

package central

import (
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores session data that needs not and should not be available to mentalese programs
type DialogContext struct {
	storage *common.FileStorage
	AnaphoraQueue *AnaphoraQueue
	ProcessList *goal.ProcessList
}

func NewDialogContext(storage *common.FileStorage, anaphoraQueue *AnaphoraQueue, processList *goal.ProcessList) *DialogContext {
	dialogContext := &DialogContext{
		storage:       storage,
		AnaphoraQueue: anaphoraQueue,
		ProcessList: processList,
	}
	dialogContext.Initialize()

	storage.Read(dialogContext)

	return dialogContext
}

func (dc *DialogContext) Initialize() {
	dc.AnaphoraQueue.Initialize()
	dc.ProcessList.Initialize()
}

func (dc *DialogContext) Store() {
	dc.storage.Write(dc)
}

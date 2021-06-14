package central

import (
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores data that should be available to multiple sentences in the dialog
type DialogContext struct {
	storage       *common.FileStorage
	AnaphoraQueue *AnaphoraQueue
	DeicticCenter *DeicticCenter
	Sentences     []*mentalese.ParseTreeNode
	ProcessList   *goal.ProcessList
}

func NewDialogContext(storage *common.FileStorage, anaphoraQueue *AnaphoraQueue, deicticCenter *DeicticCenter, processList *goal.ProcessList) *DialogContext {
	dialogContext := &DialogContext{
		storage:       storage,
		AnaphoraQueue: anaphoraQueue,
		DeicticCenter: deicticCenter,
		Sentences:     []*mentalese.ParseTreeNode{},
		ProcessList:   processList,
	}
	dialogContext.Initialize()

	storage.Read(dialogContext)

	return dialogContext
}

func (dc *DialogContext) Initialize() {
	dc.AnaphoraQueue.Initialize()
	dc.DeicticCenter.Initialize()
	dc.Sentences = []*mentalese.ParseTreeNode{}
	dc.ProcessList.Initialize()
}

func (dc *DialogContext) Store() {
	dc.storage.Write(dc)
}

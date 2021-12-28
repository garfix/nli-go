package central

import (
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores data that should be available to multiple sentences in the dialog
type DialogContext struct {
	storage           *common.FileStorage
	AnaphoraQueue     *AnaphoraQueue
	DeicticCenter     *DeicticCenter
	Sentences         []*mentalese.ParseTreeNode
	ProcessList       *goal.ProcessList
	VariableGenerator *mentalese.VariableGenerator
	DiscourseEntities *mentalese.Binding
}

func NewDialogContext(
	storage *common.FileStorage,
	anaphoraQueue *AnaphoraQueue,
	deicticCenter *DeicticCenter,
	processList *goal.ProcessList,
	variableGenerator *mentalese.VariableGenerator,
	discourseEntities *mentalese.Binding,
) *DialogContext {
	dialogContext := &DialogContext{
		storage:           storage,
		AnaphoraQueue:     anaphoraQueue,
		DeicticCenter:     deicticCenter,
		Sentences:         []*mentalese.ParseTreeNode{},
		ProcessList:       processList,
		VariableGenerator: variableGenerator,
		DiscourseEntities: discourseEntities,
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
	dc.VariableGenerator.Initialize()

	dc.DiscourseEntities.Clear()
}

func (dc *DialogContext) Store() {
	dc.storage.Write(dc)
}

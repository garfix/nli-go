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
	ProcessList       *goal.ProcessList
	VariableGenerator *mentalese.VariableGenerator
	DiscourseEntities *mentalese.Binding
	ClauseList        *mentalese.ClauseList
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
		ProcessList:       processList,
		VariableGenerator: variableGenerator,
		DiscourseEntities: discourseEntities,
		ClauseList:        mentalese.NewClauseList(),
	}
	dialogContext.Initialize()

	if storage != nil {
		storage.Read(dialogContext)
	}

	return dialogContext
}

func (dc *DialogContext) GetClauseList() *mentalese.ClauseList {
	return dc.ClauseList
}

func (dc *DialogContext) Initialize() {
	dc.AnaphoraQueue.Initialize()
	dc.DeicticCenter.Initialize()
	dc.ProcessList.Initialize()
	dc.VariableGenerator.Initialize()

	dc.DiscourseEntities.Clear()
	dc.ClauseList.Clear()
}

func (dc *DialogContext) Store() {
	if dc.storage != nil {
		dc.storage.Write(dc)
	}
}

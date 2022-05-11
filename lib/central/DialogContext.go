package central

import (
	"nli-go/lib/mentalese"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores data that should be available to multiple sentences in the dialog
type DialogContext struct {
	DeicticCenter     *DeicticCenter
	ProcessList       *ProcessList
	VariableGenerator *mentalese.VariableGenerator
	DiscourseEntities *mentalese.Binding
	ClauseList        *mentalese.ClauseList
	AnaphoraQueue     *mentalese.AnaphoraQueue
	TagList           *TagList
	Sorts             *mentalese.EntitySorts
}

func NewDialogContext(
	deicticCenter *DeicticCenter,
	processList *ProcessList,
	variableGenerator *mentalese.VariableGenerator,
) *DialogContext {
	discourseEntities := mentalese.NewBinding()
	dialogContext := &DialogContext{
		DeicticCenter:     deicticCenter,
		ProcessList:       processList,
		VariableGenerator: variableGenerator,
		DiscourseEntities: &discourseEntities,
		ClauseList:        mentalese.NewClauseList(),
		AnaphoraQueue:     mentalese.NewAnaphoraQueue(),
		TagList:           NewTagList(),
		Sorts:             mentalese.NewEntitySorts(),
	}
	dialogContext.Initialize()

	return dialogContext
}

func (e *DialogContext) GetAnaphoraQueue() []EntityReferenceGroup {
	ids := []EntityReferenceGroup{}
	clauses := e.AnaphoraQueue.GetClauses()

	for i := len(clauses) - 1; i >= 0; i-- {
		clause := clauses[i]
		for _, discourseVariable := range clause.GetDiscourseVariables() {
			value, found := e.DiscourseEntities.Get(discourseVariable)
			if found {
				if value.IsList() {
					group := EntityReferenceGroup{}
					for _, item := range value.TermValueList {
						_, sort := e.Sorts.GetSort(discourseVariable)
						if sort != item.TermSort {
							item.TermSort = ""
						}
						reference := EntityReference{sort, item.TermValue, discourseVariable}
						group = append(group, reference)
					}
					ids = append(ids, group)
				} else {
					_, sort := e.Sorts.GetSort(discourseVariable)
					if sort != value.TermSort {
						value.TermSort = ""
					}
					reference := EntityReference{sort, value.TermValue, discourseVariable}
					group := EntityReferenceGroup{reference}
					ids = append(ids, group)
				}
			} else {
				reference := EntityReference{"", "", discourseVariable}
				group := EntityReferenceGroup{reference}
				ids = append(ids, group)
			}
		}
	}

	return ids
}

func (e *DialogContext) GetAnaphoraQueue1() []EntityReferenceGroup {
	ids := []EntityReferenceGroup{}
	clauses := e.ClauseList.Clauses

	for i := len(clauses) - 1; i >= 0; i-- {
		clause := clauses[i]
		for _, entity := range clause.Entities {
			value, found := e.DiscourseEntities.Get(entity.DiscourseVariable)
			if found {
				if value.IsList() {
					group := EntityReferenceGroup{}
					for _, item := range value.TermValueList {
						reference := EntityReference{item.TermSort, item.TermValue, entity.DiscourseVariable}
						group = append(group, reference)
					}
					ids = append(ids, group)
				} else {
					reference := EntityReference{value.TermSort, value.TermValue, entity.DiscourseVariable}
					group := EntityReferenceGroup{reference}
					ids = append(ids, group)
				}
			} else {
				reference := EntityReference{"", "", entity.DiscourseVariable}
				group := EntityReferenceGroup{reference}
				ids = append(ids, group)
			}
		}
	}

	return ids
}

func (dc *DialogContext) Initialize() {
	dc.DeicticCenter.Initialize()
	dc.ProcessList.Initialize()
	dc.VariableGenerator.Initialize()

	dc.DiscourseEntities.Clear()
	dc.ClauseList.Clear()
}

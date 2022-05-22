package central

import (
	"nli-go/lib/mentalese"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores data that should be available to multiple sentences in the dialog
type DialogContext struct {
	VariableGenerator *mentalese.VariableGenerator
	DeicticCenter     *DeicticCenter
	ClauseList        *mentalese.ClauseList
	EntityBindings    *mentalese.Binding
	EntityTags        *TagList
	EntitySorts       *mentalese.EntitySorts
}

func NewDialogContext(
	deicticCenter *DeicticCenter,
	variableGenerator *mentalese.VariableGenerator,
) *DialogContext {
	discourseEntities := mentalese.NewBinding()
	dialogContext := &DialogContext{
		VariableGenerator: variableGenerator,
		DeicticCenter:     deicticCenter,
		ClauseList:        mentalese.NewClauseList(),
		EntityBindings:    &discourseEntities,
		EntityTags:        NewTagList(),
		EntitySorts:       mentalese.NewEntitySorts(),
	}
	dialogContext.Initialize()

	return dialogContext
}

func (e *DialogContext) ReplaceVariable(fromVariable string, toVariable string) {
	e.ClauseList.GetLastClause().ReplaceVariable(fromVariable, toVariable)
}

func (e *DialogContext) GetAnaphoraQueue() []AnaphoraQueueElement {
	ids := []AnaphoraQueueElement{}
	clauses := e.ClauseList.Clauses

	variableUsed := map[string]bool{}

	first := len(clauses) - 1 - MaxSizeAnaphoraQueue
	for i := len(clauses) - 1; i >= 0 && i >= first; i-- {
		clause := clauses[i]
		for _, discourseVariable := range clause.ResolvedEntities {

			// add each variable only once
			_, found := variableUsed[discourseVariable]
			if found {
				continue
			} else {
				variableUsed[discourseVariable] = true
			}

			value, found := e.EntityBindings.Get(discourseVariable)
			if found {
				if value.IsList() {
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{}}
					sorts := e.EntitySorts.GetSorts(discourseVariable)
					for i, item := range value.TermValueList {
						sort := sorts[i]
						reference := AnaphoraQueueElementValue{sort, item.TermValue}
						group.values = append(group.values, reference)
					}
					ids = append(ids, group)
				} else {
					sorts := e.EntitySorts.GetSorts(discourseVariable)
					sort := sorts[0]
					reference := AnaphoraQueueElementValue{sort, value.TermValue}
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
					ids = append(ids, group)
				}
			} else {
				sorts := e.EntitySorts.GetSorts(discourseVariable)
				sort := mentalese.SortEntity
				if len(sorts) > 0 {
					sort = sorts[0]
				}
				reference := AnaphoraQueueElementValue{sort, ""}
				group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
				ids = append(ids, group)
			}
		}
	}

	return ids
}

func (dc *DialogContext) Initialize() {
	dc.DeicticCenter.Initialize()
	dc.VariableGenerator.Initialize()

	dc.EntityBindings.Clear()
	dc.EntityTags.Clear()
	dc.EntitySorts.Clear()
	dc.ClauseList.Clear()
}

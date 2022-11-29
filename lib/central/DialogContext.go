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
	EntityBindings    *mentalese.EntityBindings
	EntityTags        *TagList
	EntitySorts       *mentalese.EntitySorts
	EntityLabels      *EntityLabels
	EntityDefinitions *EntityDefinitions
}

func NewDialogContext(
	variableGenerator *mentalese.VariableGenerator,
) *DialogContext {
	dialogContext := &DialogContext{
		VariableGenerator: variableGenerator,
		DeicticCenter:     NewDeicticCenter(),
		ClauseList:        mentalese.NewClauseList(),
		EntityBindings:    mentalese.NewEntityBindings(),
		EntityTags:        NewTagList(),
		EntitySorts:       mentalese.NewEntitySorts(),
		EntityLabels:      NewEntityLabels(),
		EntityDefinitions: NewEntityDefinitions(),
	}
	dialogContext.Initialize()

	return dialogContext
}

func (e *DialogContext) Fork() *DialogContext {

	return &DialogContext{
		VariableGenerator: e.VariableGenerator, // no copy!
		DeicticCenter:     e.DeicticCenter.Copy(),
		ClauseList:        e.ClauseList.Copy(),
		EntityBindings:    e.EntityBindings.Copy(),
		EntityTags:        e.EntityTags.Copy(),
		EntitySorts:       e.EntitySorts.Copy(),
		EntityLabels:      e.EntityLabels.Copy(),
		EntityDefinitions: e.EntityDefinitions.Copy(),
	}
}

func GetAnaphoraQueue(clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entitySorts *mentalese.EntitySorts) []AnaphoraQueueElement {
	ids := []AnaphoraQueueElement{}
	clauses := clauseList.Clauses

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

			value, found := entityBindings.Get(discourseVariable)
			if found {
				if value.IsList() {
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{}}
					sorts := entitySorts.GetSorts(discourseVariable)
					for i, item := range value.TermValueList {
						sort := sorts[i]
						reference := AnaphoraQueueElementValue{sort, item.TermValue}
						group.values = append(group.values, reference)
					}
					ids = append(ids, group)
				} else {
					sorts := entitySorts.GetSorts(discourseVariable)
					sort := sorts[0]
					reference := AnaphoraQueueElementValue{sort, value.TermValue}
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
					ids = append(ids, group)
				}
			} else {
				sorts := entitySorts.GetSorts(discourseVariable)
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

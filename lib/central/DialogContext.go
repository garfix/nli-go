package central

import (
	"nli-go/lib/mentalese"
)

const MaxSizeAnaphoraQueue = 10

// The dialog context stores data that should be available to multiple sentences in the dialog
type DialogContext struct {
	DeicticCenter     *DeicticCenter
	VariableGenerator *mentalese.VariableGenerator
	EntityBindings    *mentalese.Binding
	ClauseList        *mentalese.ClauseList
	AnaphoraQueue     *mentalese.AnaphoraQueue
	EntityTags        *TagList
	EntitySorts       *mentalese.EntitySorts
}

func NewDialogContext(
	deicticCenter *DeicticCenter,
	variableGenerator *mentalese.VariableGenerator,
) *DialogContext {
	discourseEntities := mentalese.NewBinding()
	dialogContext := &DialogContext{
		DeicticCenter:     deicticCenter,
		VariableGenerator: variableGenerator,
		EntityBindings:    &discourseEntities,
		ClauseList:        mentalese.NewClauseList(),
		AnaphoraQueue:     mentalese.NewAnaphoraQueue(),
		EntityTags:        NewTagList(),
		EntitySorts:       mentalese.NewEntitySorts(),
	}
	dialogContext.Initialize()

	return dialogContext
}

func (e *DialogContext) ReplaceVariable(fromVariable string, toVariable string) {
	clause := e.ClauseList.GetLastClause()
	newTree := clause.ParseTree.ReplaceVariable(fromVariable, toVariable)
	clause.ParseTree = &newTree

	if clause.Center != nil {
		clause.Center.Replacevariable(fromVariable, toVariable)
	}
	for _, e := range clause.Entities {
		e.Replacevariable(fromVariable, toVariable)
	}
}

func (e *DialogContext) GetAnaphoraQueue() []EntityReferenceGroup {
	ids := []EntityReferenceGroup{}
	clauses := e.AnaphoraQueue.GetClauses()

	for i := len(clauses) - 1; i >= 0; i-- {
		clause := clauses[i]
		for _, discourseVariable := range clause.GetDiscourseVariables() {
			value, found := e.EntityBindings.Get(discourseVariable)
			if found {
				if value.IsList() {
					group := EntityReferenceGroup{}
					sorts := e.EntitySorts.GetSorts(discourseVariable)
					for i, item := range value.TermValueList {
						sort := sorts[i]
						reference := EntityReference{sort, item.TermValue, discourseVariable}
						group = append(group, reference)
					}
					ids = append(ids, group)
				} else {
					sorts := e.EntitySorts.GetSorts(discourseVariable)
					sort := sorts[0]
					reference := EntityReference{sort, value.TermValue, discourseVariable}
					group := EntityReferenceGroup{reference}
					ids = append(ids, group)
				}
			} else {
				sorts := e.EntitySorts.GetSorts(discourseVariable)
				sort := mentalese.SortEntity
				if len(sorts) > 0 {
					sort = sorts[0]
				}
				reference := EntityReference{sort, "", discourseVariable}
				group := EntityReferenceGroup{reference}
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

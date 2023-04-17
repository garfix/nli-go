package central

import "nli-go/lib/mentalese"

type ReferentFinder struct {
}

func (finder *ReferentFinder) FindReferentCandidates() {

}

func GetAnaphoraQueue(clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entitySorts *mentalese.EntitySorts) []AnaphoraQueueElement {
	ids := []AnaphoraQueueElement{}
	clauses := clauseList.Clauses

	variableUsed := map[string]bool{}

	first := len(clauses) - 1 - MaxSizeAnaphoraQueue
	for i := len(clauses) - 1; i >= 0 && i >= first; i-- {

		clause := clauses[i]
		for _, discourseVariable := range clause.QueuedEntities {

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
					sort := entitySorts.GetSort(discourseVariable)
					for _, item := range value.TermValueList {
						reference := AnaphoraQueueElementValue{sort, item.TermValue}
						group.values = append(group.values, reference)
					}
					ids = append(ids, group)
				} else {
					sort := entitySorts.GetSort(discourseVariable)
					reference := AnaphoraQueueElementValue{sort, value.TermValue}
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
					ids = append(ids, group)
				}
			} else {
				sort := entitySorts.GetSort(discourseVariable)
				reference := AnaphoraQueueElementValue{sort, ""}
				group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
				ids = append(ids, group)
			}
		}
	}

	return ids
}

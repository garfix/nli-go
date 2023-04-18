package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type ReferentFinder struct {
	log               *common.SystemLog
	messenger         api.ProcessMessenger
	clauseList        *mentalese.ClauseList
	entityBindings    *mentalese.EntityBindings
	entityDefinitions *mentalese.EntityDefinitions
	entityTags        *mentalese.TagList
	entitySorts       *mentalese.EntitySorts
	sortFinder        SortFinder
}

func NewReferentFinder(log *common.SystemLog, meta *mentalese.Meta, messenger api.ProcessMessenger, clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entityDefinitions *mentalese.EntityDefinitions, entityTags *mentalese.TagList, entitySorts *mentalese.EntitySorts) *ReferentFinder {
	return &ReferentFinder{
		log:               log,
		messenger:         messenger,
		clauseList:        clauseList,
		entityBindings:    entityBindings,
		entityDefinitions: entityDefinitions,
		entitySorts:       entitySorts,
		entityTags:        entityTags,
		sortFinder:        NewSortFinder(meta, messenger),
	}
}

func (finder *ReferentFinder) FindAnaphoricReferent(variable string, referenceSort string, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool, requiresDefinition bool) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}

	groups := GetAnaphoraQueue(finder.clauseList, finder.entityBindings, finder.entitySorts)
	for _, group := range groups {

		// check if the referent is a group
		if len(group.values) == 1 {

			referent := group.values[0]

			found, foundVariable, foundTerm = finder.MatchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, entityDefinition, binding, collection, reflective)
			if found {
				if foundVariable != "" {
					if finder.checkDefinition(requiresDefinition, foundVariable) {
						goto end
					}
				} else {
					goto end
				}
			}

		} else {

			// try to match an element in the group
			for _, referent := range group.values {

				found, foundVariable, foundTerm = finder.MatchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, entityDefinition, binding, collection, reflective)
				if found {
					if foundVariable != "" {
						if finder.checkDefinition(requiresDefinition, foundVariable) {
							goto end
						}
					} else {
						goto end
					}
				}
			}
		}

	}

end:

	if found {
		if foundVariable != "" {
			finder.log.AddProduction("ref", "accept "+foundVariable+"\n")
		} else {
			finder.log.AddProduction("ref", "accept "+foundTerm.String()+"\n")
		}
	} else {
		finder.log.AddProduction("ref", "reject all\n")
	}

	return found, foundVariable, foundTerm
}

func (finder *ReferentFinder) checkDefinition(requiresDefinition bool, foundVariable string) bool {
	if requiresDefinition {
		definition := finder.entityDefinitions.Get(foundVariable)
		if !definition.IsEmpty() {
			finder.log.AddProduction("ref", foundVariable+" has a definition\n")
			return true
		}
		return false
	} else {
		return true
	}

}

func (finder *ReferentFinder) MatchReferenceToReferent(variable string, referenceSort string, referentVariable string, referentSort string, referentId string, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}
	agree := false
	mostSpecificFound := false

	finder.log.AddProduction("\nresolving", variable+"\n")

	// the entity itself is in the queue
	// should not be possible
	if referentVariable == variable {
		finder.log.AddProduction("ref", referentVariable+" equals "+variable+"\n")
		goto end
	}

	if referentSort == "" {
		finder.log.AddProduction("ref", referentVariable+" has no sort\n")
		goto end
	}
	_, mostSpecificFound = finder.sortFinder.getMostSpecific(referenceSort, referentSort)
	if !mostSpecificFound {
		finder.log.AddProduction("ref", referentVariable+" ("+referentSort+") does not have common sort with "+variable+" ("+referenceSort+")\n")
		goto end
	}

	agree, _, _ = NewAgreementChecker().CheckForCategoryConflictBetween(variable, referentVariable, finder.entityTags)
	if !agree {
		finder.log.AddProduction("ref", referentVariable+" does not agree with "+variable+"\n")
		goto end
	}

	if reflective {
		if !collection.IsCoArgument(variable, referentVariable) {
			finder.log.AddProduction("ref", referentVariable+" is not co-argument "+variable+"\n")
			goto end
		}
	} else {
		if collection.IsCoArgument(variable, referentVariable) {
			finder.log.AddProduction("ref", referentVariable+" is co-argument of "+variable+"\n")
			goto end
		}
	}

	// is this a definite reference?
	if len(entityDefinition) == 0 {
		// no: we're done
		found = true
		foundVariable = referentVariable
		finder.log.AddProduction("ref", referentVariable+" is OK\n")
	} else {
		// yes, it is a definite reference
		// a definite reference can only be checked against an id
		if referentId == "" {
			finder.log.AddProduction("ref", referentVariable+" has no id "+entityDefinition.String()+"\n")
			goto end
		} else {
			b := mentalese.NewBinding()
			value := mentalese.NewTermId(referentId, referentSort)
			b.Set(variable, value)

			refBinding := binding.Merge(b)
			testRangeBindings := finder.messenger.ExecuteChildStackFrame(entityDefinition, mentalese.InitBindingSet(refBinding))
			if testRangeBindings.GetLength() > 0 {
				// found: bind the reference variable to the id of the referent
				// (don't replace variable)
				found = true
				foundTerm = value
				finder.log.AddProduction("ref", referentVariable+" binding to id\n")
				goto end
			} else {
				finder.log.AddProduction("ref", referentVariable+" could not be bound\n")
				goto end
			}
		}
	}

end:

	finder.log.AddProduction("  ref end", referentVariable)

	return found, foundVariable, foundTerm
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

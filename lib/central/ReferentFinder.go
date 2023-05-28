package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"sort"
	"strconv"
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
		sortFinder:        NewSortFinder(messenger),
	}
}

type Referent struct {
	Variable string
	Term     mentalese.Term
	Score    int
}

func (finder *ReferentFinder) FindAnaphoricReferents(variable string, referenceSort string, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool, requiresDefinition bool) ([]Referent, bool) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}
	referents := []Referent{}
	ambiguous := false

	finder.log.AddProduction("\nresolving", variable+"\n")
	finder.log.AddProduction("definition", entityDefinition.String()+"\n\n")

	groups := GetAnaphoraQueue(finder.clauseList, finder.entityBindings, finder.entitySorts)
	for _, group := range groups {

		// check if the referent is a group
		if len(group.values) == 1 {

			referent := group.values[0]

			found, foundVariable, foundTerm = finder.MatchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, referent.Score, entityDefinition, binding, collection, reflective)
			if found {
				if foundVariable != "" {
					if finder.checkDefinition(requiresDefinition, foundVariable) {
						referents = append(referents, Referent{foundVariable, foundTerm, referent.Score})
					}
				} else {
					referents = append(referents, Referent{foundVariable, foundTerm, referent.Score})
				}
			}

		} else {

			// try to match an element in the group
			for _, referent := range group.values {

				found, foundVariable, foundTerm = finder.MatchReferenceToReferent(variable, referenceSort, group.Variable, referent.Sort, referent.Id, referent.Score, entityDefinition, binding, collection, reflective)
				if found {
					if foundVariable != "" {
						if finder.checkDefinition(requiresDefinition, foundVariable) {
							referents = append(referents, Referent{foundVariable, foundTerm, referent.Score})
						}
					} else {
						referents = append(referents, Referent{foundVariable, foundTerm, referent.Score})
					}
				}
			}
		}
	}

	// sort by score, desc
	sort.Slice(referents, func(i, j int) bool {
		return referents[i].Score > referents[j].Score
	})

	// check for ambiguity (2 highest scores are the same)
	if len(referents) > 1 {

		//diff := referents[0].Score - referents[1].Score
		// if diff == 1 {
		//println("* diff: " + strconv.Itoa(diff))
		// }

		if (referents[0].Score == referents[1].Score) &&
			(referents[0].Variable != referents[1].Variable) {
			println("* Ambiguity found!")
			ambiguous = true
		}
	}

	// return found, foundVariable, foundTerm
	return referents, ambiguous
}

func (finder *ReferentFinder) checkDefinition(requiresDefinition bool, foundVariable string) bool {
	if requiresDefinition {
		definition := finder.entityDefinitions.Get(foundVariable)
		if !definition.IsEmpty() {
			// finder.log.AddProduction("ref", foundVariable+" has a definition\n")
			return true
		}
		return false
	} else {
		return true
	}

}

func (finder *ReferentFinder) MatchReferenceToReferent(variable string, referenceSort string, referentVariable string, referentSort string, referentId string, referentScore int, entityDefinition mentalese.RelationSet, binding mentalese.Binding, collection *AnaphoraResolverCollection, reflective bool) (bool, string, mentalese.Term) {

	found := false
	foundVariable := ""
	foundTerm := mentalese.Term{}
	agree := false
	mostSpecificFound := false

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
			finder.log.AddProduction("ref", "[OK] "+referentVariable+" is co-argument of "+variable+"\n")
			goto end
		}
	}

	// is this a definite reference?
	if len(entityDefinition) == 0 {
		// no: we're done
		found = true
		foundVariable = referentVariable
		finder.log.AddProduction("ref", "[OK] ("+strconv.Itoa(referentScore)+") "+referentVariable+" is indefinite\n")
	} else {
		// yes, it is a definite reference
		// a definite reference can only be checked against an id
		if referentId == "" {
			finder.log.AddProduction("ref", referentVariable+" has no id, which is needed for a definite reference\n")
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
				finder.log.AddProduction("ref", "OK ("+strconv.Itoa(referentScore)+") "+referentVariable+" / "+value.String()+" matches the definition; binding to id\n")
				goto end
			} else {
				finder.log.AddProduction("ref", referentVariable+" / "+value.String()+" could not be bound to the definition\n")
				goto end
			}
		}
	}

end:

	return found, foundVariable, foundTerm
}

func GetAnaphoraQueue(clauseList *mentalese.ClauseList, entityBindings *mentalese.EntityBindings, entitySorts *mentalese.EntitySorts) []AnaphoraQueueElement {
	ids := []AnaphoraQueueElement{}
	clauses := clauseList.Clauses

	variableUsed := map[string]bool{}

	scoreBase := 0
	first := len(clauses) - 1 - MaxSizeAnaphoraQueue
	for i := len(clauses) - 1; i >= 0 && i >= first; i-- {

		clause := clauses[i]
		for i, discourseVariable := range clause.QueuedEntities {

			// add each variable only once
			_, found := variableUsed[discourseVariable]
			if found {
				continue
			} else {
				variableUsed[discourseVariable] = true
			}

			score := calculateScore(scoreBase, clause, i, discourseVariable)

			value, found := entityBindings.Get(discourseVariable)
			if found {
				if value.IsList() {
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{}}
					sort := entitySorts.GetSort(discourseVariable)
					for _, item := range value.TermValueList {
						reference := AnaphoraQueueElementValue{sort, item.TermValue, score}
						group.values = append(group.values, reference)
					}
					ids = append(ids, group)
				} else {
					sort := entitySorts.GetSort(discourseVariable)
					reference := AnaphoraQueueElementValue{sort, value.TermValue, score}
					group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
					ids = append(ids, group)
				}
			} else {
				sort := entitySorts.GetSort(discourseVariable)
				reference := AnaphoraQueueElementValue{sort, "", score}
				group := AnaphoraQueueElement{Variable: discourseVariable, values: []AnaphoraQueueElementValue{reference}}
				ids = append(ids, group)
			}
		}

		scoreBase -= 10
	}

	return ids
}

func calculateScore(scoreBase int, clause *mentalese.Clause, index int, variable string) int {
	score := scoreBase
	if index == 0 {
		score += 1
	}
	for _, function := range clause.SyntacticFunctions {
		if function.SyntacticFunction == mentalese.AtomFunctionSubject {
			score += 5
		} else if function.SyntacticFunction == mentalese.AtomFunctionObject {
			score += 3
		}
	}
	return score
}

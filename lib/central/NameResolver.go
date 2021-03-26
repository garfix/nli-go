package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"sort"
)

// uses SharedId() relations to create new sense() relations, that hold the SharedId's of the entities of different knowledge bases
type NameResolver struct {
	solverAsync   *ProblemSolverAsync
	meta          *mentalese.Meta
	log           *common.SystemLog
}

func NewNameResolver(solverAsync *ProblemSolverAsync, meta *mentalese.Meta, log *common.SystemLog) *NameResolver {
	return &NameResolver{
		solverAsync:   solverAsync,
		meta:          meta,
		log:           log,
	}
}

func (resolver *NameResolver) Choose(messenger api.ProcessMessenger, nameInformations []NameInformation) ([]NameInformation, bool) {

	resolvedInformations := []NameInformation{}

	names := mentalese.TermList{}

	for _, nameInformation := range nameInformations {
		names = append(names, mentalese.NewTermString(nameInformation.Information))
	}

	// go:wait_for(go:user_select(['A', 'B', 'C'], Selection))
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateWaitFor, []mentalese.Term{
			mentalese.NewTermRelationSet(
				mentalese.RelationSet{
					mentalese.NewRelation(false, mentalese.PredicateUserSelect, []mentalese.Term{
						mentalese.NewTermList(names),
						mentalese.NewTermVariable("Selection"),
					}),
				}),
		}),
	}

	bindings, loading := messenger.ExecuteChildStackFrameAsync(set, mentalese.InitBindingSet(mentalese.NewBinding()))
	if loading {
		return resolvedInformations, true
	}

	for _, binding := range bindings.GetAll() {
		selection := binding.MustGet("Selection")
		index, _ := selection.GetIntValue()
		resolvedInformations = append(resolvedInformations, nameInformations[index])
	}

	return resolvedInformations, false
}

func (resolver *NameResolver) ResolveName(name string, sort string) []NameInformation {

	factBaseNameInformations := []NameInformation{}

	for _, factBase := range resolver.solverAsync.factBases {
		nameInformations := resolver.resolveNameInFactBase(name, sort, factBase)
		factBaseNameInformations = append(factBaseNameInformations, nameInformations...)
	}

	return factBaseNameInformations
}

func (resolver *NameResolver) resolveNameInFactBase(name string, inducedSort string, factBase api.FactBase) []NameInformation {

	var nameInformations []NameInformation

	// go through all sorts
	for aSort, entityInfo := range resolver.meta.GetSorts() {

		if inducedSort != "entity" && inducedSort != "" && !resolver.meta.MatchesSort(aSort, inducedSort) {
			continue
		}

		b := mentalese.NewBinding()
		b.Set(mentalese.NameVar, mentalese.NewTermString(name))

		bindings := resolver.solverAsync.FindFacts(factBase, entityInfo.Name, b)

		for _, binding := range bindings.GetAll() {

			id, _ := binding.Get(mentalese.IdVar)

			// sort because the resulting strings must not be in random order
			sortedInfoTypes := []string{}
			for infoType := range entityInfo.Knownby {
				sortedInfoTypes = append(sortedInfoTypes, infoType)
			}
			sort.Strings(sortedInfoTypes)

			information := ""
			sep := ""

			for _, infoType := range sortedInfoTypes {

				relationSet := entityInfo.Knownby[infoType]

				// create a relation set for each field that gives Information about this name
				b := mentalese.NewBinding()
				b.Set(mentalese.IdVar, mentalese.NewTermId(id.TermValue, aSort))
				bindings2 := resolver.solverAsync.FindFacts(factBase, relationSet, b)

				for _, binding2 := range bindings2.GetAll() {
					value, _ := binding2.Get(mentalese.ValueVar)
					information += sep + value.TermValue
					sep = ";"
					// DBPedia sometimes returns multiple results for a date, while there should be only one
					break
				}
			}

			nameInformations = append(nameInformations, NameInformation{
				Name:         name,
				DatabaseName: factBase.GetName(),
				EntityType:   id.TermSort,
				SharedId:     id.TermValue,
				Information:  information,
			})
		}
	}

	return nameInformations
}
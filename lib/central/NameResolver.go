package central

import (
	"nli-go/lib/api"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"sort"
)

// uses SharedId() relations to create new sense() relations, that hold the SharedId's of the entities of different knowledge bases
type NameResolver struct {
	solverAsync *ProblemSolver
	meta        *mentalese.Meta
	log         *common.SystemLog
}

func NewNameResolver(solverAsync *ProblemSolver, meta *mentalese.Meta, log *common.SystemLog) *NameResolver {
	return &NameResolver{
		solverAsync: solverAsync,
		meta:        meta,
		log:         log,
	}
}

func (resolver *NameResolver) ExtractNames(rootNode mentalese.ParseTreeNode, rootVariables []string) mentalese.Binding {
	return resolver.extractSenseFromNode(rootNode, rootVariables)
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (resolver *NameResolver) extractSenseFromNode(node mentalese.ParseTreeNode, antecedentVariables []string) mentalese.Binding {

	nameBinding := resolver.extractName(node, antecedentVariables)

	// create relations for each of the children
	for i, childNode := range node.Constituents {

		consequentVariables := node.Rule.GetConsequentVariables(i)

		childNameBinding := resolver.extractSenseFromNode(*childNode, consequentVariables)
		nameBinding = nameBinding.Merge(childNameBinding)
	}

	return nameBinding
}

func (resolver *NameResolver) extractName(node mentalese.ParseTreeNode, antecedentVariables []string) mentalese.Binding {

	names := mentalese.NewBinding()

	if node.Category != mentalese.CategoryProperNounGroup {
		return names
	}

	variable := antecedentVariables[0]

	name := ""
	sep := ""
	for _, properNounNode := range node.Constituents {
		name += sep + properNounNode.Form
		sep = " "
	}
	names.Set(variable, mentalese.NewTermString(name))

	return names
}

func (resolver *NameResolver) Choose(messenger api.ProcessMessenger, nameInformations []NameInformation) ([]NameInformation, bool) {

	resolvedInformations := []NameInformation{}

	names := mentalese.TermList{}

	for _, nameInformation := range nameInformations {
		names = append(names, mentalese.NewTermString(nameInformation.Information))
	}

	// go:wait_for(go:user_select('Which one?', ['A', 'B', 'C'], Selection))
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateWaitFor, []mentalese.Term{
			mentalese.NewTermString(mentalese.MessageChoose),
			mentalese.NewTermVariable("Selection"),
			mentalese.NewTermString(common.WhichOne),
			mentalese.NewTermList(names),
		}, 0),
	}

	bindings := messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(mentalese.NewBinding()))

	for _, binding := range bindings.GetAll() {
		selection := binding.MustGet("Selection")
		index, _ := selection.GetIntValue()
		resolvedInformations = append(resolvedInformations, nameInformations[index])
	}

	return resolvedInformations, false
}

func (resolver *NameResolver) ResolveName(name string, sort string, messenger api.ProcessMessenger) []NameInformation {

	factBaseNameInformations := []NameInformation{}

	for _, factBase := range resolver.solverAsync.factBases {
		nameInformations := resolver.resolveNameInFactBase(name, sort, factBase, messenger)
		factBaseNameInformations = append(factBaseNameInformations, nameInformations...)
	}

	return factBaseNameInformations
}

func (resolver *NameResolver) resolveNameInFactBase(name string, inducedSort string, factBase api.FactBase, messenger api.ProcessMessenger) []NameInformation {

	var nameInformations []NameInformation

	// go through all sorts
	for aSort, entityInfo := range resolver.meta.GetSorts() {

		if inducedSort != "entity" && inducedSort != "" && inducedSort != aSort {

			isa := mentalese.NewRelation(false, mentalese.PredicateIsa, []mentalese.Term{
				mentalese.NewTermAtom(aSort),
				mentalese.NewTermAtom(inducedSort),
			}, 0)

			bindings := messenger.ExecuteChildStackFrame(mentalese.RelationSet{isa}, mentalese.InitBindingSet(mentalese.NewBinding()))
			if bindings.IsEmpty() {
				continue
			}
		}

		b := mentalese.NewBinding()
		b.Set(mentalese.NameVar, mentalese.NewTermString(name))

		bindings := resolver.solverAsync.FindFacts(factBase, entityInfo.Name, b)

		for _, binding := range bindings.GetAll() {

			id, _ := binding.Get(mentalese.IdVar)

			// gender
			gender := ""
			genderInBinding := mentalese.NewBinding()
			genderInBinding.Set(mentalese.IdVar, id)

			if entityInfo.Gender.Predicate != "" {
				genderOutBindings := resolver.solverAsync.FindFacts(factBase, entityInfo.Gender, genderInBinding)
				for _, genderOutBinding := range genderOutBindings.GetAll() {
					value, _ := genderOutBinding.Get(mentalese.ValueVar)
					gender = value.TermValue
				}
			}

			// number
			number := ""
			numberInBinding := mentalese.NewBinding()
			numberInBinding.Set(mentalese.IdVar, id)

			if entityInfo.Number.Predicate != "" {
				numberOutBindings := resolver.solverAsync.FindFacts(factBase, entityInfo.Number, numberInBinding)
				for _, numberOutBinding := range numberOutBindings.GetAll() {
					value, _ := numberOutBinding.Get(mentalese.ValueVar)
					number = value.TermValue
				}
			}

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
				Gender:       gender,
				Number:       number,
				DatabaseName: factBase.GetName(),
				EntityType:   id.TermSort,
				SharedId:     id.TermValue,
				Information:  information,
			})
		}
	}

	return nameInformations
}

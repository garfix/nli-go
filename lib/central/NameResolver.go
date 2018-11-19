package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"sort"
	"strconv"
)

const predicateNameInformation = "name_information"

// uses EntityId() relations to create new sense() relations, that hold the EntityId's of the entities of different knowledge bases
type NameResolver struct {
	solver 			*ProblemSolver
	matcher 		*mentalese.RelationMatcher
	log 			*common.SystemLog
	dialogContext   *DialogContext
}

func NewNameResolver(solver *ProblemSolver, matcher *mentalese.RelationMatcher, log *common.SystemLog, dialogContext *DialogContext) *NameResolver {
	return &NameResolver{
		solver:      	solver,
		matcher:	 	matcher,
		log: 			log,
		dialogContext:	dialogContext,
	}
}

// Returns a set of senses, or a human readable question to the user
func (resolver *NameResolver) Resolve(relations mentalese.RelationSet) (mentalese.RelationSet, string) {

	names := resolver.collectNames(relations)

	senses := mentalese.RelationSet{}
	userQuestion := ""

	for variable, name := range names {

		// check if the name is known in the dialog context
		dialogNameInformations := resolver.RetrieveNameInformations(name)

		if len(dialogNameInformations) == 0 {

			// look up the name in all fact bases

			multipleResultsInFactBase := false
			factBasesWithResults := 0

			var factBaseNameInformations []NameInformation

			for _, factBase := range resolver.solver.factBases {
				factBaseNameInformations = append(factBaseNameInformations, resolver.resolveName(name, factBase)...)

				if len(factBaseNameInformations) > 0 {
					factBasesWithResults++
				}

				if len(factBaseNameInformations) > 1 {
					multipleResultsInFactBase = true
				}
			}

			if factBasesWithResults > 1 || multipleResultsInFactBase {

				// check if the user has just answered this question
				answer, found := resolver.dialogContext.GetAnswerToOpenQuestion()

				if found {

					dialogNameInformations = resolver.selectNameInformationsFromAnswer(factBaseNameInformations, answer)
					resolver.SaveNameInformations(name, dialogNameInformations)
					resolver.dialogContext.RemoveAnswerToOpenQuestion()

				} else {

					// need to ask user
					userQuestion = resolver.composeUserQuestion(factBaseNameInformations)
					resolver.dialogContext.SetOpenQuestion(userQuestion)
					break

				}

			} else {

				// single meaning for name
				dialogNameInformations = factBaseNameInformations
				resolver.SaveNameInformations(name, dialogNameInformations)

			}
		}

		senses = append(senses, resolver.createNameSensesFromNameInformations(dialogNameInformations, variable)...)
	}

	return senses, userQuestion
}

func (resolver *NameResolver) selectNameInformationsFromAnswer(nameInformations []NameInformation, answer string) []NameInformation {
	answerNameInformations := []NameInformation{}

	for _, nameInformation := range nameInformations {
		if nameInformation.GetIdentifier() == answer {
			answerNameInformations = append(answerNameInformations, nameInformation)
		}
	}

	return answerNameInformations
}

func (resolver *NameResolver) composeUserQuestion(nameInformations []NameInformation) string {
	question := "Which one?"

	for _, nameInformation := range nameInformations {
		question += " [" + nameInformation.GetIdentifier() + "] " + nameInformation.Information
	}

	return question
}

func (resolver *NameResolver) createNameSensesFromNameInformations(nameInformations []NameInformation, variable string) mentalese.RelationSet {

	senses := mentalese.RelationSet{}

	for _, nameInformation := range nameInformations {

		sense := mentalese.NewRelation(mentalese.PredicateSense, []mentalese.Term{
			mentalese.NewVariable(variable),
			mentalese.NewString(nameInformation.DatabaseName),
			mentalese.NewId(nameInformation.EntityId),
		})

		senses = append(senses, sense)

	}

	return senses
}

func (resolver *NameResolver) SaveNameInformations(name string, nameInformations []NameInformation) {

	for _, nameInformation := range nameInformations {
		resolver.dialogContext.AddRelation(mentalese.NewRelation(predicateNameInformation, []mentalese.Term{
			mentalese.NewString(name),
			mentalese.NewString(nameInformation.DatabaseName),
			mentalese.NewString(nameInformation.EntityId),
		}))
	}
}

func (resolver *NameResolver) RetrieveNameInformations(name string) []NameInformation {

	bindings := resolver.dialogContext.FindRelations(mentalese.RelationSet{
		mentalese.NewRelation(predicateNameInformation, []mentalese.Term{
			mentalese.NewString(name),
			mentalese.NewVariable("databaseName"),
			mentalese.NewVariable("entityId"),
		}),
	})

	nameInformations := []NameInformation{}

	for _, binding := range bindings {
		nameInformations = append(nameInformations, NameInformation{
			Name:         name,
			DatabaseName: binding["databaseName"].TermValue,
			EntityId:     binding["entityId"].TermValue,
		})
	}

	return nameInformations
}

// in: EntityId(E1, "de", 1) EntityId(E1, "Boer", 2) EntityId(E1, "Jan", 0)
// out: E1: "Jan de Boer"
func (resolver *NameResolver) collectNames(relations mentalese.RelationSet) map[string]string {

	nameTree := map[string]map[int]string{}

	for _, relation := range relations {
		if relation.Predicate == mentalese.PredicateName {
			variable := relation.Arguments[0].TermValue
			value := relation.Arguments[1].TermValue
			indexString := relation.Arguments[2]
			index, err := strconv.Atoi(indexString.TermValue)
			if err == nil {
				_, found := nameTree[variable]
				if !found {
					nameTree[variable] = map[int]string{}
				}
				nameTree[variable][index] = value
			}
		}
	}

	names := map[string]string{}

	for variable, branch := range nameTree {

		name := ""

		for i := 1; i <= len(branch); i++ {
			value, found := branch[i]
			if found {
				if name == "" {
					name = value
				} else {
					name = name + " " + value
				}
			}
		}

		names[variable] = name
	}

	return names
}

func (resolver *NameResolver) resolveName(name string, factBase knowledge.FactBase) []NameInformation {

	var nameInformations []NameInformation

	// go through all entity types
	for entityType, entityInfo := range factBase.GetEntities() {

		// create a relation set for the name
		boundNameSet := resolver.matcher.BindRelationSetSingleBinding(entityInfo.Name, mentalese.Binding{
			mentalese.NameVar: mentalese.NewString(name),
		})

		// ask fact base for entities (ids) with this name
		bindings := resolver.solver.FindFacts(factBase, boundNameSet)

		for _, binding := range bindings {

			id, _ := binding[mentalese.IdVar]
			information := entityType

			// sort because the resulting strings must not be in random order
			sortedInfoTypes := []string{}
			for infoType := range entityInfo.Knownby {
				sortedInfoTypes = append(sortedInfoTypes, infoType)
			}
			sort.Strings(sortedInfoTypes)

			for _, infoType := range sortedInfoTypes {

				relationSet := entityInfo.Knownby[infoType]

				// create a relation set for each field that gives Information about this name
				boundRelationSet := resolver.matcher.BindRelationSetSingleBinding(relationSet, mentalese.Binding{
					mentalese.IdVar: mentalese.NewId(id.TermValue),
				})

				bindings2 := resolver.solver.FindFacts(factBase, boundRelationSet)

				for _, binding2 := range bindings2 {
					value, _ := binding2[mentalese.ValueVar]
					information += "; " + infoType + ": " + value.TermValue
					// DBPedia sometimes returns multiple results for a date, while there should be only one
					break
				}
			}

			nameInformations = append(nameInformations, NameInformation{
				Name: name,
				DatabaseName: factBase.GetName(),
				EntityId:     id.TermValue,
				Information:  information,
			})
		}
	}

	return nameInformations
}
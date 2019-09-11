package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"sort"
	"strconv"
)

const predicateNameInformation = "name_information"

const defaultEntityType = "entity"

// uses EntityId() relations to create new sense() relations, that hold the EntityId's of the entities of different knowledge bases
type NameResolver struct {
	solver 			*ProblemSolver
	matcher 		*mentalese.RelationMatcher
	predicates      mentalese.Predicates
	log 			*common.SystemLog
	dialogContext   *DialogContext
}

func NewNameResolver(solver *ProblemSolver, matcher *mentalese.RelationMatcher, predicates mentalese.Predicates, log *common.SystemLog, dialogContext *DialogContext) *NameResolver {
	return &NameResolver{
		solver:      	solver,
		matcher:	 	matcher,
		predicates:		predicates,
		log: 			log,
		dialogContext:	dialogContext,
	}
}

// Returns a set of senses, or a human readable question to the user
func (resolver *NameResolver) Resolve(relations mentalese.RelationSet) (*mentalese.KeyCabinet, mentalese.RelationSet) {

	keyCabinet := mentalese.NewKeyCabinet()
	namelessRelations := mentalese.RelationSet{}
	userResponse := ""
	options := common.NewOptions()

	namesAndTypes := resolver.collectNamesAndTypes(relations.UnScope())

	for variable, nameAndType := range namesAndTypes {

		name := nameAndType.name
		entityType := nameAndType.entityType

		// check if the nameAndType is known in the dialog context
		dialogNameInformations := resolver.RetrieveNameInDialogContext(name)

		if len(dialogNameInformations) == 0 {

			// look up the nameAndType in all fact bases

			multipleResultsInFactBase := false
			factBasesWithResults := 0

			var factBaseNameInformations []NameInformation

			for _, factBase := range resolver.solver.factBases {
				factBaseNameInformations = append(factBaseNameInformations, resolver.resolveName(name, entityType, factBase)...)

				if len(factBaseNameInformations) > 0 {
					factBasesWithResults++
				}

				if len(factBaseNameInformations) > 1 {
					multipleResultsInFactBase = true
				}
			}

			if factBasesWithResults == 0 {

				userResponse = "Name not found in any knowledge base: " + name

			} else if factBasesWithResults > 1 || multipleResultsInFactBase {

				// check if the user has just answered this question
				answer, found := resolver.dialogContext.GetAnswerToOpenQuestion()

				if found {

					dialogNameInformations = resolver.selectNameInformationsFromAnswer(factBaseNameInformations, answer)
					resolver.SaveNameInformations(name, dialogNameInformations)
					resolver.dialogContext.RemoveAnswerToOpenQuestion()

				} else {

					// need to ask user
					userResponse = "Which one?"
					options = resolver.composeOptions(factBaseNameInformations)

					// store options
					resolver.storeOptions(factBaseNameInformations)

					break

				}

			} else {

				// single meaning for nameAndType
				dialogNameInformations = factBaseNameInformations
				resolver.SaveNameInformations(name, dialogNameInformations)

			}
		}

		for _, info := range dialogNameInformations {
			keyCabinet.AddName(variable, info.DatabaseName, info.EntityId)
		}
	}

	namelessRelations = relations.RemoveMatchingPredicate(mentalese.PredicateName)

	if userResponse != "" {
		resolver.log.SetClarificationRequest(userResponse, options)
	}

	return keyCabinet, namelessRelations
}

func (resolver *NameResolver) storeOptions(nameInformations []NameInformation) {

	for _, nameInformation := range nameInformations {
		resolver.dialogContext.AddOption(nameInformation.GetIdentifier())
	}
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

func (resolver *NameResolver) composeOptions(nameInformations []NameInformation) *common.Options {

	options := &common.Options{}

	for _, nameInformation := range nameInformations {
		options.AddOption(nameInformation.GetIdentifier(), nameInformation.Information)
	}

	return options
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

func (resolver *NameResolver) RetrieveNameInDialogContext(name string) []NameInformation {

	bindings := resolver.solver.FindFacts(resolver.dialogContext.GetFactBase(), mentalese.RelationSet{
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

type nameInfo struct {
	name string
	entityType string
}

// in: EntityId(E1, "de", 1) EntityId(E1, "Boer", 2) EntityId(E1, "Jan", 0)
// out: E1: "Jan de Boer"
func (resolver *NameResolver) collectNamesAndTypes(relations mentalese.RelationSet) map[string]nameInfo {

	nameTree := map[string]map[int]string{}

	for _, relation := range relations {
		if relation.Predicate == mentalese.PredicateName {
			variable := relation.Arguments[0].TermValue
			name := relation.Arguments[1].TermValue
			indexString := relation.Arguments[2]
			index, err := strconv.Atoi(indexString.TermValue)
			if err == nil {
				_, found := nameTree[variable]
				if !found {
					nameTree[variable] = map[int]string{}
				}
				nameTree[variable][index] = name
			}
		}
	}

	names := map[string]nameInfo{}

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

		entityType := resolver.getEntityTypeFromRelations(variable, relations)

		names[variable] = nameInfo{ name: name, entityType: entityType }
	}

	return names
}

func (resolver *NameResolver) getEntityTypeFromRelations(variable string, relations mentalese.RelationSet) string {

	entityType := defaultEntityType

	for _, relation := range relations {
		predicate := relation.Predicate

		for i, argument := range relation.Arguments {
			if argument.IsVariable() && argument.TermValue == variable {

				entityTypes, found := resolver.predicates[predicate]

				if found {

					if entityType != defaultEntityType && entityTypes.EntityTypes[i] != entityType {
						panic("Conflict in entity types!")
					}

					entityType = entityTypes.EntityTypes[i]
				}
			}
		}
	}

	return entityType
}

func (resolver *NameResolver) resolveName(name string, inducedEntityType string, factBase knowledge.FactBase) []NameInformation {

	var nameInformations []NameInformation

	// go through all entity types
	for entityType, entityInfo := range factBase.GetEntities() {

		if inducedEntityType != defaultEntityType && entityType != inducedEntityType {
			continue
		}

		bindings := resolver.solver.SolveRelationSet(entityInfo.Name, nil, mentalese.Bindings{{
			mentalese.NameVar: mentalese.NewString(name),
		}})

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
				boundRelationSet := relationSet.BindSingle(mentalese.Binding{
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
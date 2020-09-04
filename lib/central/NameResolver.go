package central

import (
	"nli-go/lib/common"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"sort"
)

const defaultEntityType = "entity"

// uses SharedId() relations to create new sense() relations, that hold the SharedId's of the entities of different knowledge bases
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

func (resolver *NameResolver) Resolve(nameInformations []NameInformation) []NameInformation {

	resolvedInformations := []NameInformation{}
	userResponse := ""
	options := common.NewOptions()

	if len(nameInformations) > 0 {

		name := nameInformations[0].Name

		dialogNameInformations := resolver.RetrieveNameInDialogContext(name)

		if len(dialogNameInformations) == 0 {

			multipleResultsInFactBase, factBasesWithResults := resolver.collectMetaData(nameInformations)

			if factBasesWithResults > 1 || multipleResultsInFactBase {

				// check if the user has just answered this question
				answer, found := resolver.dialogContext.GetAnswerToOpenQuestion()

				if found {

					dialogNameInformations = resolver.selectNameInformationsFromAnswer(nameInformations, answer)
					resolver.SaveNameInformations(name, dialogNameInformations)
					resolver.dialogContext.RemoveAnswerToOpenQuestion()

				} else {

					// need to ask user
					userResponse = "Which one?"
					options = resolver.composeOptions(nameInformations)

					// store options
					resolver.storeOptions(nameInformations)

				}
			} else {

				// single meaning for nameAndType
				dialogNameInformations = nameInformations
				resolver.SaveNameInformations(name, dialogNameInformations)

			}
		}

		resolvedInformations = dialogNameInformations
	}

	if userResponse != "" {
		resolver.log.SetClarificationRequest(userResponse, options)
	}

	return resolvedInformations
}

func (resolver *NameResolver) collectMetaData(nameInformations []NameInformation) (bool, int) {

	factBases := map[string]bool{}

	multipleResultsInFactBase := false

	for _, nameInformation := range nameInformations {

		_, found := factBases[nameInformation.DatabaseName]
		if found {
			multipleResultsInFactBase = true
		} else {
			factBases[nameInformation.DatabaseName] = true
		}
	}

	return multipleResultsInFactBase, len(factBases)
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

func (resolver *NameResolver) SaveNameInformations(name string, nameInformations []NameInformation) {

	resolver.dialogContext.AddNameInformations(nameInformations)
}

func (resolver *NameResolver) RetrieveNameInDialogContext(name string) []NameInformation {

	nameInformations := []NameInformation{}

	for _, nameInformation := range resolver.dialogContext.GetNameInformations() {
		if nameInformation.Name == name {
			nameInformations = append(nameInformations, nameInformation)
		}
	}

	return nameInformations
}

type nameInfo struct {
	name string
	entityType string
}

func (resolver *NameResolver) ResolveName(name string, entityType string) []NameInformation {

	factBaseNameInformations := []NameInformation{}

	for _, factBase := range resolver.solver.factBases {
		factBaseNameInformations = append(factBaseNameInformations, resolver.resolveNameInFactBase(name, entityType, factBase)...)
	}

	return factBaseNameInformations
}

func (resolver *NameResolver) resolveNameInFactBase(name string, inducedEntityType string, factBase knowledge.FactBase) []NameInformation {

	var nameInformations []NameInformation

	// go through all entity types
	for entityType, entityInfo := range factBase.GetEntities() {

		if inducedEntityType != defaultEntityType && entityType != inducedEntityType {
			continue
		}

		bindings := resolver.solver.FindFacts(factBase, entityInfo.Name, mentalese.Binding{
			mentalese.NameVar: mentalese.NewTermString(name),
		})

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
				bindings2 := resolver.solver.FindFacts(factBase, relationSet, mentalese.Binding{
					mentalese.IdVar: mentalese.NewTermId(id.TermValue, entityType),
				})

				for _, binding2 := range bindings2 {
					value, _ := binding2[mentalese.ValueVar]
					information += "; " + infoType + ": " + value.TermValue
					// DBPedia sometimes returns multiple results for a date, while there should be only one
					break
				}
			}

			nameInformations = append(nameInformations, NameInformation{
				Name:         name,
				DatabaseName: factBase.GetName(),
				EntityType:   id.TermEntityType,
				SharedId:     id.TermValue,
				Information:  information,
			})
		}
	}

	return nameInformations
}
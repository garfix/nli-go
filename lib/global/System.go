package global

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"strconv"
)

type System struct {
	log                   *common.SystemLog
	dialogContext         *central.DialogContext
	dialogContextStorage  *DialogContextFileStorage
	internalGrammarParser *importer.InternalGrammarParser
	nameResolver          *central.NameResolver
	grammars              []parse.Grammar
	parser                *earley.Parser
	meta                  *mentalese.Meta
	relationizer          *earley.Relationizer
	matcher               *central.RelationMatcher
	solver                *central.ProblemSolver
	answerer              *central.Answerer
	generator             *generate.Generator
	surfacer              *generate.SurfaceRepresentation
}

func (system *System) PopulateDialogContext(sessionId string, clearWhenCorrupt bool) {
	system.dialogContextStorage.Read(sessionId, system.dialogContext, clearWhenCorrupt)
}

func (system *System) ClearDialogContext() {
	system.dialogContext.Initialize()
}

func (system *System) StoreDialogContext(sessionId string) {
	system.dialogContextStorage.Write(sessionId, system.dialogContext)
}

func (system *System) RemoveDialogContext(sessionId string) {
	system.dialogContextStorage.Remove(sessionId)
}

// Low-level function to inspect the internal state of the system
func (system *System) Query(relations string) mentalese.BindingSet {
	set := system.internalGrammarParser.CreateRelationSet(relations)
	return system.solver.SolveRelationSet(set, mentalese.InitBindingSet( mentalese.NewBinding()))
}

func (system *System) Answer(input string) (string, *common.Options) {

	// process possible user responses and start with the original question
	originalInput := system.dialogContext.Process(input)

	// process it (again)
	answer, options := system.Process(originalInput)

	// does the System ask the user a question?
	if !options.HasOptions() {
		// the original question has been answered
		system.dialogContext.RemoveOriginalInput()
	}

	return answer, options
}

func (system *System) Process(originalInput string) (string, *common.Options) {

	options := common.NewOptions()
	sortFinder := central.NewSortFinder(system.meta)
	namesProcessed := false
	nameNotFound := ""
	answer := ""
	tokens := []string{}
	parseTrees := []earley.ParseTreeNode{}
	requestRelations := mentalese.RelationSet{}
	answerRelations := mentalese.RelationSet{}
	answerWords := []string{}
	names := mentalese.NewBinding()
	entityIds := mentalese.NewBinding()

	system.log.AddProduction("Anaphora queue before", system.dialogContext.AnaphoraQueue.String())

	for _, grammar := range system.grammars {

		if !system.log.IsDone() {
			tokens = grammar.GetTokenizer().Process(originalInput)
			system.log.AddProduction("TokenExpression", fmt.Sprintf("%v", tokens))
		}

		if !system.log.IsDone() {
			parseTrees = system.parser.Parse(grammar.GetReadRules(), tokens)
			if len(parseTrees) == 0 {
				system.log.AddError("Parser returned no parse trees")
			} else {
				system.log.AddProduction("Parse trees found", strconv.Itoa(len(parseTrees)))
			}
		}

		if !system.log.IsDone() {
			for _, aTree := range parseTrees {

				system.log.AddProduction("Parser", aTree.String())

				requestRelations, names = system.relationizer.Relationize(aTree)
				system.log.AddProduction("Relationizer", requestRelations.String())

				// extract sorts: variable => sort
				sorts, sortFound := sortFinder.FindSorts(requestRelations)
				if !sortFound {
					// conflicting sorts
					system.log.AddProduction("Break", "Breaking due to conflicting sorts: " + sorts.String())
					goto next
				}

				// look up entity ids by name
				entityIds = mentalese.NewBinding()
				for variable, name := range names.GetAll() {

					// find sort
					sort, found := sorts[variable]
					if !found {
						system.log.AddProduction("Info",
							"The name '" + name.TermValue + "' could not be looked up because no sort could be derived from the relations.")
						if nameNotFound == "" {
							nameNotFound = name.TermValue
						}
						goto next
					}

					// find name information
					nameInformations := system.nameResolver.ResolveName(name.TermValue, sort)
					if len(nameInformations) == 0 {
						system.log.AddProduction("Info",
							"Database lookup for name '" + name.TermValue + "'  with sort '" + sort + "' did not give any results")
						nameNotFound = name.TermValue
						goto next
					}

					// make the user choose one entity from multiple with the same name
					if len(nameInformations) > 1 {
						nameInformations = system.nameResolver.Resolve(nameInformations)
					}

					// link variable to ID
					for _, nameInformation := range nameInformations {
						entityIds.Set(variable, mentalese.NewTermId(nameInformation.SharedId, nameInformation.EntityType))
					}
				}

				// names found and linked to id
				for _, value := range entityIds.GetAll() {
					system.dialogContext.AnaphoraQueue.AddReferenceGroup(
						central.EntityReferenceGroup{ central.CreateEntityReference(value.TermValue, value.TermSort) })
				}
				system.log.AddProduction("Named entities", entityIds.String())

				// success!
				namesProcessed = true
				break

				next:
			}
		}

		if !namesProcessed && nameNotFound != "" {
			answer = common.NameNotFound + ": " + nameNotFound
			system.log.AddError(answer)
		}

		if !system.log.IsDone() {
			answerRelations = system.answerer.Answer(requestRelations, mentalese.InitBindingSet(entityIds))
			system.log.AddProduction("Answer", answerRelations.String())
			system.log.AddProduction("Anaphora queue after", system.dialogContext.AnaphoraQueue.String())
		}

		if !system.log.IsDone() {
			answerWords = system.generator.Generate(grammar.GetWriteRules(), answerRelations)
			system.log.AddProduction("Answer Words", fmt.Sprintf("%v", answerWords))
		}

		if !system.log.IsDone() {
			answer = system.surfacer.Create(answerWords)
			system.log.AddProduction("Answer", fmt.Sprintf("%v", answer))
		}

		// for now, just use the first grammar
		break
	}

	if system.log.GetClarificationQuestion() != "" {
		answer = system.log.GetClarificationQuestion()
		options = system.log.GetClarificationOptions()
		system.log.SetClarificationRequest("", &common.Options{})
	}

	return answer, options
}

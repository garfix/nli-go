package global

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"strconv"
)

type System struct {
	log                  *common.SystemLog
	dialogContext        *central.DialogContext
	dialogContextStorage *DialogContextFileStorage
	nameResolver         *central.NameResolver
	grammars             []parse.Grammar
	parser               *earley.Parser
	meta                 *mentalese.Meta
	relationizer         *earley.Relationizer
	matcher              *mentalese.RelationMatcher
	solver               *central.ProblemSolver
	answerer             *central.Answerer
	generator            *generate.Generator
	surfacer             *generate.SurfaceRepresentation
}

func (system *System) PopulateDialogContext(sessionDataPath string, clearWhenCorrupt bool) {
	system.dialogContextStorage.Read(sessionDataPath, system.dialogContext, clearWhenCorrupt)
}

func (system *System) ClearDialogContext() {
	system.dialogContext.Initialize()
}

func (system *System) StoreDialogContext(sessionDataPath string) {
	system.dialogContextStorage.Write(sessionDataPath, system.dialogContext)
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
	answer := ""
	tokens := []string{}
	parseTree := earley.ParseTreeNode{}
	requestRelations := mentalese.RelationSet{}
	answerRelations := mentalese.RelationSet{}
	answerWords := []string{}
	nameBinding := mentalese.Binding{}

	system.log.AddProduction("Anaphora queue", system.dialogContext.AnaphoraQueue.String())

	for _, grammar := range system.grammars {

		if !system.log.IsDone() {
			tokens = grammar.GetTokenizer().Process(originalInput)
			system.log.AddProduction("TokenExpression", fmt.Sprintf("%v", tokens))
		}

		if !system.log.IsDone() {
			parseTrees := system.parser.Parse(grammar.GetReadRules(), tokens)
			if len(parseTrees) == 0 {
				system.log.AddError("Parser returned no parse trees")
			} else {
				parseTree = parseTrees[0]
				system.log.AddProduction("Parse trees found: ", strconv.Itoa(len(parseTrees)))
				system.log.AddProduction("Parser", parseTree.String())
			}
		}

		if !system.log.IsDone() {
			requestRelations, nameBinding = system.relationizer.Relationize(parseTree, system.nameResolver)
			system.storeNamedEntities(nameBinding)
			system.log.AddProduction("Relationizer", requestRelations.String())
			system.log.AddProduction("Named entities", nameBinding.String())
		}

		if !system.log.IsDone() {
			answerRelations = system.answerer.Answer(requestRelations, []mentalese.Binding{nameBinding})
			system.log.AddProduction("Answer", answerRelations.String())
			system.log.AddProduction("Anaphora queue", system.dialogContext.AnaphoraQueue.String())
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

func (system System) storeNamedEntities(binding mentalese.Binding) {
	 for _, value := range binding {
		 system.dialogContext.AnaphoraQueue.AddReferenceGroup(central.EntityReferenceGroup{ central.CreateEntityReference(value.TermValue, value.TermEntityType) })
	 }
}
package global

import (
	"encoding/json"
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"os"
	"path/filepath"
	"strconv"
)

type system struct {
	log               *common.SystemLog
	dialogContext     *central.DialogContext
	dialogContextStorage *DialogContextFileStorage
	nameResolver      *central.NameResolver
	lexicon           *parse.Lexicon
	grammar           *parse.Grammar
	generationLexicon *generate.GenerationLexicon
	generationGrammar *generate.GenerationGrammar
	tokenizer         *parse.Tokenizer
	parser            *earley.Parser
	relationizer      *earley.Relationizer
	answerer          *central.Answerer
	generator         *generate.Generator
	surfacer          *generate.SurfaceRepresentation
}

func NewSystem(configPath string, log *common.SystemLog) *system {

	system := &system{ log: log }
	config := system.ReadConfig(configPath, log)

	if log.IsOk() {
		builder := NewSystemBuilder(filepath.Dir(configPath), log)
		builder.BuildFromConfig(system, config)
	}

	return system
}

func (system *system) ReadConfig(configPath string, log *common.SystemLog) (systemConfig) {

	config := systemConfig{}

	configJson, err := common.ReadFile(configPath)
	if err != nil {
		log.AddError("Error reading JSON file " + configPath + " (" + err.Error() + ")")
	}

	if log.IsOk() {
		err := json.Unmarshal([]byte(configJson), &config)
		if err != nil {
			log.AddError("Error parsing JSON file " + configPath + " (" + err.Error() + ")")
		}
	}

	if config.ParentConfig != "" {
		parentConfigPath := config.ParentConfig
		if len(parentConfigPath) > 0 && parentConfigPath[0] != os.PathSeparator {
			parentConfigPath = filepath.Dir(configPath) + string(os.PathSeparator) + parentConfigPath
		}
		parentConfig := system.ReadConfig(parentConfigPath, log)

		config = parentConfig.Merge(config)
		config.ParentConfig = ""
	}

	return config
}

func (system *system) PopulateDialogContext(sessionDataPath string) {
	system.dialogContextStorage.Read(sessionDataPath, system.dialogContext)
}

func (system *system) ClearDialogContext() {
	system.dialogContext.Initialize()
}

func (system *system) StoreDialogContext(sessionDataPath string) {
	system.dialogContextStorage.Write(sessionDataPath, system.dialogContext)
}

func (system *system) Answer(input string) (string, *common.Options) {

	// process possible user responses and start with the original question
	originalInput := system.dialogContext.Process(input)

	// process it (again)
	answer, options := system.Process(originalInput)

	// does the system ask the user a question?
	if !options.HasOptions() {
		// the original question has been answered
		system.dialogContext.RemoveOriginalInput()
	}

	return answer, options
}

func (system *system) Process(originalInput string) (string, *common.Options) {

	options := common.NewOptions()
	answer := ""
	tokens := []string{}
	parseTree := earley.ParseTreeNode{}
	requestRelations := mentalese.RelationSet{}
	answerRelations := mentalese.RelationSet{}
	answerWords := []string{}
	nameBinding := mentalese.Binding{}

	system.log.AddProduction("Anaphora queue", system.dialogContext.AnaphoraQueue.String())

	if !system.log.IsDone() {
		tokens = system.tokenizer.Process(originalInput)
		system.log.AddProduction("Tokenizer", fmt.Sprintf("%v", tokens))
	}

	if !system.log.IsDone() {
		parseTrees := system.parser.Parse(tokens)
		system.log.AddProduction("Parse trees found: ", strconv.Itoa(len(parseTrees)))
		if len(parseTrees) == 0 {
			system.log.AddError("Parser returned no parse trees")
		} else {
			parseTree = parseTrees[0]
			system.log.AddProduction("Parser", parseTree.String())
		}
	}

	if !system.log.IsDone() {
		requestRelations, nameBinding = system.relationizer.Relationize(parseTree, system.nameResolver)
		system.log.AddProduction("Relationizer", requestRelations.String())
	}

	if !system.log.IsDone() {
		answerRelations = system.answerer.Answer(requestRelations, []mentalese.Binding{ nameBinding })
		system.log.AddProduction("Answer", answerRelations.String())
	}

	if !system.log.IsDone() {
		for _, id := range answerRelations.GetIds() {
			system.dialogContext.AnaphoraQueue.AddEntityReference(central.CreateEntityReference(id.TermValue, id.TermEntityType))
		}
		system.log.AddProduction("Anaphora queue", system.dialogContext.AnaphoraQueue.String())
	}

	if !system.log.IsDone() {
		answerWords = system.generator.Generate(answerRelations)
		system.log.AddProduction("Answer Words", fmt.Sprintf("%v", answerWords))
	}

	if !system.log.IsDone() {
		answer = system.surfacer.Create(answerWords)
		system.log.AddProduction("Answer", fmt.Sprintf("%v", answer))
	}

	if system.log.GetClarificationQuestion() != "" {
		answer = system.log.GetClarificationQuestion()
		options = system.log.GetClarificationOptions()
	}

	return answer, options
}

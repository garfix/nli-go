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
	system.dialogContext.Initialize([]mentalese.Relation{})
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
	keyCabinet := mentalese.NewKeyCabinet()

	if !system.log.IsDone() {
		system.log.AddProduction("Dialog Context", system.dialogContext.GetRelations().String())
	}

	if !system.log.IsDone() {
		tokens = system.tokenizer.Process(originalInput)
		system.log.AddProduction("Tokenizer", fmt.Sprintf("%v", tokens))
	}

	if !system.log.IsDone() {
		parseTree = system.parser.Parse(tokens)
		system.log.AddProduction("Parser", parseTree.String())
	}

	if !system.log.IsDone() {
		requestRelations = system.relationizer.Relationize(parseTree, keyCabinet, system.nameResolver)
		system.log.AddProduction("Relationizer", requestRelations.String())
		system.log.AddProduction("Key cabinet", keyCabinet.String())
	}

	if !system.log.IsDone() {
		answerRelations = system.answerer.Answer(requestRelations, keyCabinet)
		system.log.AddProduction("DS Answer", answerRelations.String())
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

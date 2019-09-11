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
	quantifierScoper  mentalese.QuantifierScoper
	relationizer      earley.Relationizer
	transformer       *mentalese.RelationTransformer
	answerer          *central.Answerer
	generator         *generate.Generator
	surfacer          *generate.SurfaceRepresentation
	generic2ds        []mentalese.RelationTransformation
	ds2generic        []mentalese.RelationTransformation
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
	syntacticRelations := mentalese.RelationSet{}
	dsRelations := mentalese.RelationSet{}
	namelessDsRelations := mentalese.RelationSet{}
	dsAnswer := mentalese.RelationSet{}
	genericAnswer := mentalese.RelationSet{}
	answerWords := []string{}

	var keyCabinet *mentalese.KeyCabinet

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
		syntacticRelations = system.relationizer.Relationize(parseTree)
		system.log.AddProduction("Relationizer", syntacticRelations.String())
	}

	//if !system.log.IsDone() {
	//	dsRelations = system.transformer.Replace(system.generic2ds, syntacticRelations)
	//	system.log.AddProduction("Generic 2 DS", dsRelations.String())
	//}
	dsRelations = syntacticRelations

	if !system.log.IsDone() {
		keyCabinet, namelessDsRelations = system.nameResolver.Resolve(dsRelations)
		system.log.AddProduction("Nameless", namelessDsRelations.String())
		system.log.AddProduction("Key cabinet", keyCabinet.String())
	}

	if !system.log.IsDone() {
		dsAnswer = system.answerer.Answer(namelessDsRelations, keyCabinet)
		system.log.AddProduction("DS Answer", dsAnswer.String())
	}

	if !system.log.IsDone() {
		genericAnswer = system.transformer.Replace(system.ds2generic, dsAnswer)
		system.log.AddProduction("Generic Answer", genericAnswer.String())
	}

	if !system.log.IsDone() {
		answerWords = system.generator.Generate(genericAnswer)
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

func (system *system) Suggest(input string) []string {

	system.log.Clear()

	tokens := system.tokenizer.Process(input)

	if system.log.IsOk() {
		system.log.AddProduction("Tokenizer", fmt.Sprintf("%v", tokens))
	} else {
		return []string{}
	}

	suggests := system.parser.Suggest(tokens)

	if system.log.IsOk() {
		system.log.AddProduction("Answer Words", fmt.Sprintf("%v", suggests))
	} else {
		return []string{}
	}

	return suggests
}

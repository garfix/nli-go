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

func (system *system) Answer(input string) string {

	originalInput := system.dialogContext.Process(input)

	if system.log.IsOk() {
		system.log.AddProduction("Dialog Context", system.dialogContext.GetRelations().String())
	} else {
		return ""
	}

	tokens := system.tokenizer.Process(originalInput)

	if system.log.IsOk() {
		system.log.AddProduction("Tokenizer", fmt.Sprintf("%v", tokens))
	} else {
		return ""
	}

	parseTree := system.parser.Parse(tokens)

	if system.log.IsOk() {
		system.log.AddProduction("Parser", parseTree.String())
	} else {
		return ""
	}

	genericRelations := system.relationizer.Relationize(parseTree)

	if system.log.IsOk() {
		system.log.AddProduction("Relationizer", genericRelations.String())
	} else {
		return ""
	}

	// name(E5, "John") => name(E5, "John") reference(E5, 'dbpedia', <http://dbpedia.org/resource/John>)
	// each access to a data store, replace E5 with its ID
	nameStore, namelessRelations, userQuestion := system.nameResolver.Resolve(genericRelations)

	if userQuestion == "" {
		system.log.AddProduction("NameResolver", namelessRelations.String())
		system.log.AddProduction("NameResolver", nameStore.String())
	} else {
		return userQuestion
	}

	dsRelations := system.transformer.Replace(system.generic2ds, namelessRelations)

	if system.log.IsOk() {
		system.log.AddProduction("Generic 2 DS", dsRelations.String())
	} else {
		return ""
	}

	dsAnswer := system.answerer.Answer(dsRelations, nameStore)

	if system.log.IsOk() {
		system.log.AddProduction("DS Answer", dsAnswer.String())
	} else {
		return ""
	}

	genericAnswer := system.transformer.Replace(system.ds2generic, dsAnswer)

	if system.log.IsOk() {
		system.log.AddProduction("Generic Answer", genericAnswer.String())
	} else {
		return ""
	}

	answerWords := system.generator.Generate(genericAnswer)

	if system.log.IsOk() {
		system.log.AddProduction("Answer Words", fmt.Sprintf("%v", answerWords))
	} else {
		return ""
	}

	answer := system.surfacer.Create(answerWords)

	if system.log.IsOk() {
		system.log.AddProduction("Answer", fmt.Sprintf("%v", answer))
	} else {
		return ""
	}

	return answer
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

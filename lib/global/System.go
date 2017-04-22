package global

import (
    "nli-go/lib/parse"
    "nli-go/lib/central"
    "nli-go/lib/mentalese"
    "nli-go/lib/parse/earley"
    "nli-go/lib/generate"
    "nli-go/lib/common"
    "encoding/json"
    "path/filepath"
)

type system struct {
    log *systemLog
    lexicon *parse.Lexicon
    grammar *parse.Grammar
    generationLexicon *generate.GenerationLexicon
    generationGrammar *generate.GenerationGrammar
    tokenizer *parse.Tokenizer
    parser *earley.Parser
    quantifierScoper mentalese.QuantifierScoper
    relationizer earley.Relationizer
    transformer *mentalese.RelationTransformer
    answerer *central.Answerer
    generator *generate.Generator
    surfacer *generate.SurfaceRepresentation
    generic2ds []mentalese.RelationTransformation
    ds2generic []mentalese.RelationTransformation
}

func NewSystem(configPath string, log *systemLog) *system {

    system := &system{ log: log }
    logBlock := NewLogBlock("Build system")
    config := systemConfig{}

    configJson, err := common.ReadFile(configPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading JSON file " + configPath)
        logBlock.AddLine(err.Error())
    }

    if logBlock.IsOk() {
        err := json.Unmarshal([]byte(configJson), &config)
        if err != nil {
            logBlock.Fail()
            logBlock.AddLine("Error parsing config file " + configPath)
            logBlock.AddLine(err.Error())
        }
    }

    if logBlock.IsOk() {
        builder := newSystemBuilder(filepath.Dir(configPath))
        builder.buildFromConfig(system, config, logBlock)
    }

    log.AddBlock(logBlock)

    return system
}

func (system *system) Process(input string) (string, bool) {

    tokens := system.tokenizer.Process(input)
    parseTree, _ := system.parser.Parse(tokens)
    rawRelations := system.relationizer.Relationize(parseTree)
    genericRelations := system.transformer.Replace(system.generic2ds, rawRelations)
    domainSpecificSense := system.quantifierScoper.Scope(genericRelations)
    dsAnswer := system.answerer.Answer(domainSpecificSense)
    genericAnswer := system.transformer.Replace(system.ds2generic, dsAnswer)
    answerWords := system.generator.Generate(genericAnswer)
    answer := system.surfacer.Create(answerWords)

    return answer, true
}

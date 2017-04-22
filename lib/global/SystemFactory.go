package global

import (
    "encoding/json"
    "nli-go/lib/common"
    "path/filepath"
    "nli-go/lib/parse/earley"
    "nli-go/lib/mentalese"
    "nli-go/lib/knowledge"
    "nli-go/lib/central"
    "nli-go/lib/generate"
    "nli-go/lib/parse"
)

type systemFactory struct {
}

func NewSystemFactory() systemFactory {
    return systemFactory{}
}

func (systemFactory systemFactory) NewSystem(configPath string, log *systemLog) (*system, bool) {

    system := &system{ log: log }
    logBlock := NewLogBlock("Build system")
    config := systemConfig{}
    ok := true

    system.lexicon = parse.NewLexicon()
    system.grammar = parse.NewGrammar()
    system.generationLexicon = generate.NewGenerationLexicon()
    system.generationGrammar = generate.NewGenerationGrammar()
    system.tokenizer = parse.NewTokenizer()
    system.parser = earley.NewParser(system.grammar, system.lexicon)
    system.quantifierScoper = mentalese.NewQuantifierScoper()
    system.relationizer = earley.NewRelationizer(system.lexicon)
    system.generic2ds = []mentalese.RelationTransformation{}
    system.ds2generic = []mentalese.RelationTransformation{}

    systemFunctionBase := knowledge.NewSystemFunctionBase()
    matcher := mentalese.NewRelationMatcher()
    matcher.AddFunctionBase(systemFunctionBase)
    system.transformer = mentalese.NewRelationTransformer(matcher)

    systemPredicateBase := knowledge.NewSystemPredicateBase()
    system.answerer = central.NewAnswerer(matcher)
    system.answerer.AddMultipleBindingsBase(systemPredicateBase)

    system.generator = generate.NewGenerator(system.generationGrammar, system.generationLexicon)
    system.surfacer = generate.NewSurfaceRepresentation()

    configJson, err := common.ReadFile(configPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading JSON file " + configPath)
        logBlock.AddLine(err.Error())
        ok = false
    }

    if ok {
        err := json.Unmarshal([]byte(configJson), &config)
        if err != nil {
            logBlock.Fail()
            logBlock.AddLine("Error parsing config file " + configPath)
            logBlock.AddLine(err.Error())
            ok = false
        }
    }

    if ok {
        builder := newSystemBuilder(filepath.Dir(configPath))
        builder.buildFromConfig(system, config, logBlock)
    }

    log.AddBlock(logBlock)

    return system, ok
}

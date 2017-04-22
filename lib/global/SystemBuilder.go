package global

import (
    "nli-go/lib/common"
    "nli-go/lib/importer"
    "nli-go/lib/knowledge"
)

type systemBuilder struct {
    baseDir string
    parser *importer.InternalGrammarParser
}

func newSystemBuilder(baseDir string) systemBuilder {

    parser := importer.NewInternalGrammarParser()
    parser.SetPanicOnParseFail(false)

    return systemBuilder{
        baseDir: baseDir,
        parser: parser,
    }
}

func (builder systemBuilder) buildFromConfig(system *system, config systemConfig, logBlock *LogBlock) {

    systemFacts := builder.parser.CreateRelationSet(`[
		act(question, _)
		focus(_)
		every(_)
		isa(_, do)
	]`)

    ds2system := builder.parser.CreateDbMappings(`[
		act(question, X) ->> act(question, X);
		focus(A) ->> focus(A);
		every(A) ->> every(A);
		isa(_, B) ->> isa(_, B);
	]`)

    system.answerer.AddFactBase(knowledge.NewInMemoryFactBase(systemFacts, ds2system))

    for _, lexiconPath := range config.Lexicons {
        path := common.AbsolutePath(builder.baseDir, lexiconPath)
        builder.ImportLexiconFromPath(system, path, logBlock)
    }
    for _, grammarPath := range config.Grammars {
        path := common.AbsolutePath(builder.baseDir, grammarPath)
        builder.ImportGrammarFromPath(system, path, logBlock)
    }
    for _, lexiconPath := range config.Generationlexicons {
        path := common.AbsolutePath(builder.baseDir, lexiconPath)
        builder.ImportGenerationLexiconFromPath(system, path, logBlock)
    }
    for _, grammarPath := range config.Generationgrammars {
        path := common.AbsolutePath(builder.baseDir, grammarPath)
        builder.ImportGenerationGrammarFromPath(system, path, logBlock)
    }
    for _, ruleBasePath := range config.Rulebases {
        path := common.AbsolutePath(builder.baseDir, ruleBasePath)
        builder.ImportRuleBaseFromPath(system, path, logBlock)
    }
    for _, factBase := range config.Factbases.Relation {
        builder.ImportRelationSetFactBase(system, factBase, logBlock)
    }
    for _, factBase := range config.Factbases.Mysql {
        builder.ImportMySqlDatabase(system, factBase, logBlock)
    }
    for _, solutionBasePath := range config.Solutions {
        builder.ImportSolutionBaseFromPath(system, solutionBasePath, logBlock)
    }
    for _, transformationsPath := range config.Ds2generic {
        path := common.AbsolutePath(builder.baseDir, transformationsPath)
        builder.ImportDs2GenericTransformations(system, path, logBlock)
    }
    for _, transformationsPath := range config.Generic2ds {
        path := common.AbsolutePath(builder.baseDir, transformationsPath)
        builder.ImportGeneric2DsTransformations(system, path, logBlock)
    }
}

func (builder systemBuilder) ImportLexiconFromPath(system *system, lexiconPath string, logBlock *LogBlock) {

    lexiconString, err := common.ReadFile(lexiconPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine(err.Error())
        return
    }

    lexicon := builder.parser.CreateLexicon(lexiconString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing lexicon file " + lexiconPath)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.ImportLexicon(lexicon)
}

func (builder systemBuilder) ImportGrammarFromPath(system *system, grammarPath string, logBlock *LogBlock) {

    grammarString, err := common.ReadFile(grammarPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine(err.Error())
        return
    }

    grammar := builder.parser.CreateGrammar(grammarString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing grammar file " + grammarPath)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.ImportGrammar(grammar)
}

func (builder systemBuilder) ImportGenerationLexiconFromPath(system *system, lexiconPath string, logBlock *LogBlock) {

    lexiconString, err := common.ReadFile(lexiconPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine(err.Error())
        return
    }

    lexicon := builder.parser.CreateGenerationLexicon(lexiconString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing lexicon file " + lexiconPath)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.ImportGenerationLexicon(lexicon)
}

func (builder systemBuilder) ImportGenerationGrammarFromPath(system *system, grammarPath string, logBlock *LogBlock) {

    grammarString, err := common.ReadFile(grammarPath)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine(err.Error())
        return
    }

    grammar := builder.parser.CreateGenerationGrammar(grammarString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing grammar file " + grammarPath)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.ImportGenerationGrammar(grammar)
}

func (builder systemBuilder) ImportRuleBaseFromPath(system *system, ruleBasePath string, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, ruleBasePath)
    ruleBaseString, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading rules " + path)
        logBlock.AddLine(err.Error())
        return
    }

    rules := builder.parser.CreateRules(ruleBaseString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing rules file " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.answerer.AddRuleBase(knowledge.NewRuleBase(rules))
}

func (builder systemBuilder) ImportRelationSetFactBase(system *system, factBase relationSetFactBase, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, factBase.Facts)
    factString, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading facts " + path)
        logBlock.AddLine(err.Error())
        return
    }

    facts := builder.parser.CreateRelationSet(factString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing facts " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    path = common.AbsolutePath(builder.baseDir, factBase.Map)
    mapString, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading map " + path)
        logBlock.AddLine(err.Error())
        return
    }

    dbMap := builder.parser.CreateDbMappings(mapString)
    lastResult = builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing map " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.answerer.AddFactBase(knowledge.NewInMemoryFactBase(facts, dbMap))
}

func (builder systemBuilder) ImportMySqlDatabase(system *system, factBase mysqlFactBase, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, factBase.Map)
    mapString, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading map " + path)
        logBlock.AddLine(err.Error())
        return
    }

    dbMap := builder.parser.CreateDbMappings(mapString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing map " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    database := knowledge.NewMySqlFactBase(factBase.Domain, factBase.Username, factBase.Password, factBase.Database, dbMap)

    for _, table := range factBase.Tables  {
        columns := []string{}
        for _, column := range table.Columns {
            columns = append(columns, column.Name)
        }
        database.AddTableDescription(table.Name, columns)
    }

    if factBase.Enabled {
        system.answerer.AddFactBase(database)
    }
}

func (builder systemBuilder) ImportSolutionBaseFromPath(system *system, solutionBasePath string, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, solutionBasePath)
    solutionString, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading solutions " + path)
        logBlock.AddLine(err.Error())
        return
    }

    solutions := builder.parser.CreateSolutions(solutionString)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing solutions " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    system.answerer.AddSolutions(solutions)
}


func (builder systemBuilder) ImportGeneric2DsTransformations(system *system, transformationsPath string, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, transformationsPath)
    transformationstring, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading transformations " + path)
        logBlock.AddLine(err.Error())
        return
    }

    transformations := builder.parser.CreateTransformations(transformationstring)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing transformations " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    for _, transformation := range transformations {
        system.generic2ds = append(system.generic2ds, transformation)
    }
}

func (builder systemBuilder) ImportDs2GenericTransformations(system *system, transformationsPath string, logBlock *LogBlock) {

    path := common.AbsolutePath(builder.baseDir, transformationsPath)
    transformationstring, err := common.ReadFile(path)
    if err != nil {
        logBlock.Fail()
        logBlock.AddLine("Error reading transformations " + path)
        logBlock.AddLine(err.Error())
        return
    }

    transformations := builder.parser.CreateTransformations(transformationstring)
    lastResult := builder.parser.GetLastParseResult()
    if !lastResult.Ok {
        logBlock.Fail()
        logBlock.AddLine("Error parsing transformations " + path)
        logBlock.AddLine(lastResult.String())
        return
    }

    for _, transformation := range transformations {
        system.ds2generic = append(system.ds2generic, transformation)
    }
}

package global

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
)

type systemBuilder struct {
	log     *common.SystemLog
	baseDir string
	parser  *importer.InternalGrammarParser
}

func newSystemBuilder(baseDir string, log *common.SystemLog) systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return systemBuilder{
		baseDir: baseDir,
		parser:  parser,
		log:     log,
	}
}

func (builder systemBuilder) buildFromConfig(system *system, config systemConfig) {

	system.lexicon = parse.NewLexicon()
	system.grammar = parse.NewGrammar()
	system.generationLexicon = generate.NewGenerationLexicon(builder.log)
	system.generationGrammar = generate.NewGenerationGrammar()
	system.tokenizer = parse.NewTokenizer(builder.log)
	system.parser = earley.NewParser(system.grammar, system.lexicon, builder.log)
	system.quantifierScoper = mentalese.NewQuantifierScoper(builder.log)
	system.relationizer = earley.NewRelationizer(system.lexicon, builder.log)
	system.generic2ds = []mentalese.RelationTransformation{}
	system.ds2generic = []mentalese.RelationTransformation{}

	systemFunctionBase := knowledge.NewSystemFunctionBase()
	matcher := mentalese.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)
	system.transformer = mentalese.NewRelationTransformer(matcher, builder.log)

	systemPredicateBase := knowledge.NewSystemPredicateBase(builder.log)
	system.answerer = central.NewAnswerer(matcher, builder.log)
	system.answerer.AddMultipleBindingsBase(systemPredicateBase)

	system.generator = generate.NewGenerator(system.generationGrammar, system.generationLexicon, builder.log)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)

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

	system.answerer.AddFactBase(knowledge.NewInMemoryFactBase(systemFacts, ds2system, builder.log))

	for _, lexiconPath := range config.Lexicons {
		path := common.AbsolutePath(builder.baseDir, lexiconPath)
		builder.ImportLexiconFromPath(system, path)
	}
	for _, grammarPath := range config.Grammars {
		path := common.AbsolutePath(builder.baseDir, grammarPath)
		builder.ImportGrammarFromPath(system, path)
	}
	for _, lexiconPath := range config.Generationlexicons {
		path := common.AbsolutePath(builder.baseDir, lexiconPath)
		builder.ImportGenerationLexiconFromPath(system, path)
	}
	for _, grammarPath := range config.Generationgrammars {
		path := common.AbsolutePath(builder.baseDir, grammarPath)
		builder.ImportGenerationGrammarFromPath(system, path)
	}
	for _, ruleBasePath := range config.Rulebases {
		path := common.AbsolutePath(builder.baseDir, ruleBasePath)
		builder.ImportRuleBaseFromPath(system, path)
	}
	for _, factBase := range config.Factbases.Relation {
		builder.ImportRelationSetFactBase(system, factBase)
	}
	for _, factBase := range config.Factbases.Mysql {
		builder.ImportMySqlDatabase(system, factBase)
	}
	for _, solutionBasePath := range config.Solutions {
		builder.ImportSolutionBaseFromPath(system, solutionBasePath)
	}
	for _, transformationsPath := range config.Ds2generic {
		path := common.AbsolutePath(builder.baseDir, transformationsPath)
		builder.ImportDs2GenericTransformations(system, path)
	}
	for _, transformationsPath := range config.Generic2ds {
		path := common.AbsolutePath(builder.baseDir, transformationsPath)
		builder.ImportGeneric2DsTransformations(system, path)
	}
}

func (builder systemBuilder) ImportLexiconFromPath(system *system, lexiconPath string) {

	lexiconString, err := common.ReadFile(lexiconPath)
	if err != nil {
		builder.log.Fail(err.Error())
		return
	}

	lexicon := builder.parser.CreateLexicon(lexiconString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing lexicon file " + lexiconPath + " (" + lastResult.String() + ")")
		return
	}

	system.lexicon.ImportFrom(lexicon)
}

func (builder systemBuilder) ImportGrammarFromPath(system *system, grammarPath string) {

	grammarString, err := common.ReadFile(grammarPath)
	if err != nil {
		builder.log.Fail(err.Error())
		return
	}

	grammar := builder.parser.CreateGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing grammar file " + grammarPath + " (" + lastResult.String() + ")")
		return
	}

	system.grammar.ImportFrom(grammar)
}

func (builder systemBuilder) ImportGenerationLexiconFromPath(system *system, lexiconPath string) {

	lexiconString, err := common.ReadFile(lexiconPath)
	if err != nil {
		builder.log.Fail(err.Error())
		return
	}

	lexicon := builder.parser.CreateGenerationLexicon(lexiconString, builder.log)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing lexicon file " + lexiconPath + " (" + lastResult.String() + ")")
		return
	}

	system.generationLexicon.ImportFrom(lexicon)
}

func (builder systemBuilder) ImportGenerationGrammarFromPath(system *system, grammarPath string) {

	grammarString, err := common.ReadFile(grammarPath)
	if err != nil {
		builder.log.Fail(err.Error())
		return
	}

	grammar := builder.parser.CreateGenerationGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing grammar file " + grammarPath + " (" + lastResult.String() + ")")
		return
	}

	system.generationGrammar.ImportFrom(grammar)
}

func (builder systemBuilder) ImportRuleBaseFromPath(system *system, ruleBasePath string) {

	path := common.AbsolutePath(builder.baseDir, ruleBasePath)
	ruleBaseString, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading rules " + path + " (" + err.Error() + ")")
		return
	}

	rules := builder.parser.CreateRules(ruleBaseString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing rules file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.answerer.AddRuleBase(knowledge.NewRuleBase(rules, builder.log))
}

func (builder systemBuilder) ImportRelationSetFactBase(system *system, factBase relationSetFactBase) {

	path := common.AbsolutePath(builder.baseDir, factBase.Facts)
	factString, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading facts " + path + " (" + err.Error() + ")")
		return
	}

	facts := builder.parser.CreateRelationSet(factString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing facts file " + path + " (" + lastResult.String() + ")")
		return
	}

	path = common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading map file " + path + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateDbMappings(mapString)
	lastResult = builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing map file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.answerer.AddFactBase(knowledge.NewInMemoryFactBase(facts, dbMap, builder.log))
}

func (builder systemBuilder) ImportMySqlDatabase(system *system, factBase mysqlFactBase) {

	path := common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading map file " + path + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateDbMappings(mapString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing map file " + path + " (" + lastResult.String() + ")")
		return
	}

	database := knowledge.NewMySqlFactBase(factBase.Domain, factBase.Username, factBase.Password, factBase.Database, dbMap, builder.log)

	for _, table := range factBase.Tables {
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

func (builder systemBuilder) ImportSolutionBaseFromPath(system *system, solutionBasePath string) {

	path := common.AbsolutePath(builder.baseDir, solutionBasePath)
	solutionString, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading solutions file " + path + " (" + err.Error() + ")")
		return
	}

	solutions := builder.parser.CreateSolutions(solutionString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing solutions file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.answerer.AddSolutions(solutions)
}

func (builder systemBuilder) ImportGeneric2DsTransformations(system *system, transformationsPath string) {

	path := common.AbsolutePath(builder.baseDir, transformationsPath)
	transformationstring, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading transformations file " + path + " (" + err.Error() + ")")
		return
	}

	transformations := builder.parser.CreateTransformations(transformationstring)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing transformations file " + path + " (" + lastResult.String() + ")")
		return
	}

	for _, transformation := range transformations {
		system.generic2ds = append(system.generic2ds, transformation)
	}
}

func (builder systemBuilder) ImportDs2GenericTransformations(system *system, transformationsPath string) {

	path := common.AbsolutePath(builder.baseDir, transformationsPath)
	transformationstring, err := common.ReadFile(path)
	if err != nil {
		builder.log.Fail("Error reading transformations file " + path + " (" + err.Error() + ")")
		return
	}

	transformations := builder.parser.CreateTransformations(transformationstring)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.Fail("Error parsing transformations file " + path + " (" + lastResult.String() + ")")
		return
	}

	for _, transformation := range transformations {
		system.ds2generic = append(system.ds2generic, transformation)
	}
}

package global

import (
	"encoding/json"
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

func NewSystemBuilder(baseDir string, log *common.SystemLog) systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return systemBuilder{
		baseDir: baseDir,
		parser:  parser,
		log:     log,
	}
}

func (builder systemBuilder) BuildFromConfig(system *system, config systemConfig) {

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function")
	matcher := mentalese.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)

	system.lexicon = parse.NewLexicon()
	system.grammar = parse.NewGrammar()
	system.generationLexicon = generate.NewGenerationLexicon(builder.log, matcher)
	system.generationGrammar = generate.NewGenerationGrammar()
	system.tokenizer = parse.NewTokenizer(builder.log)
	system.quantifierScoper = mentalese.NewQuantifierScoper(builder.log)
	system.relationizer = earley.NewRelationizer(system.lexicon, builder.log)
	system.generic2ds = []mentalese.RelationTransformation{}
	system.ds2generic = []mentalese.RelationTransformation{}
	system.transformer = mentalese.NewRelationTransformer(matcher, builder.log)
	system.dialogContext = central.NewDialogContext(matcher, builder.log)

	solver := central.NewProblemSolver(matcher, system.dialogContext, builder.log)

	solver.AddFunctionBase(systemFunctionBase)

	systemAggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", builder.log)
	solver.AddMultipleBindingsBase(systemAggregateBase)

	nestedStructureBase := knowledge.NewSystemNestedStructureBase(builder.log)
	solver.AddNestedStructureBase(nestedStructureBase)

	predicates, _ := builder.CreatePredicates(config.Predicates)

	system.dialogContextStorage = NewDialogContextFileStorage(builder.log)
	system.nameResolver = central.NewNameResolver(solver, matcher, predicates, builder.log, system.dialogContext)
	system.parser = earley.NewParser(system.grammar, system.lexicon, system.nameResolver, builder.log)
	system.answerer = central.NewAnswerer(matcher, solver, builder.log)
	system.generator = generate.NewGenerator(system.generationGrammar, system.generationLexicon, builder.log, matcher)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)

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
		builder.ImportRuleBaseFromPath(solver, path)
	}
	for _, factBase := range config.Factbases.Relation {
		builder.ImportInMemoryFactBase(factBase.Name, solver, factBase, matcher)
	}
	for _, factBase := range config.Factbases.Mysql {
		builder.ImportMySqlDatabase(factBase.Database, solver, system.nameResolver, factBase, matcher)
	}
	for _, factBase := range config.Factbases.Sparql {
		builder.ImportSparqlDatabase(factBase.Name, solver, system.nameResolver, factBase, matcher)
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

func (builder systemBuilder) CreatePredicates(path string) (mentalese.Predicates, bool) {

	predicates := mentalese.Predicates{}

	if path != "" {

		absolutePath := common.AbsolutePath(builder.baseDir, path)

		content, err := common.ReadFile(absolutePath)
		if err != nil {
			builder.log.AddError("Error reading predicates file " + absolutePath + " (" + err.Error() + ")")
			return predicates, false
		}

		err = json.Unmarshal([]byte(content), &predicates)
		if err != nil {
			builder.log.AddError("Error parsing predicates file " + absolutePath + " (" + err.Error() + ")")
			return predicates, false
		}
	}

	return predicates, true
}

func (builder systemBuilder) ImportLexiconFromPath(system *system, lexiconPath string) {

	lexiconString, err := common.ReadFile(lexiconPath)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	lexicon := builder.parser.CreateLexicon(lexiconString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing lexicon file " + lexiconPath + " (" + lastResult.String() + ")")
		return
	}

	system.lexicon.ImportFrom(lexicon)
}

func (builder systemBuilder) ImportGrammarFromPath(system *system, grammarPath string) {

	grammarString, err := common.ReadFile(grammarPath)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	grammar := builder.parser.CreateGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + grammarPath + " (" + lastResult.String() + ")")
		return
	}

	system.grammar.ImportFrom(grammar)
}

func (builder systemBuilder) ImportGenerationLexiconFromPath(system *system, lexiconPath string) {

	lexiconString, err := common.ReadFile(lexiconPath)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	lexicon := builder.parser.CreateGenerationLexicon(lexiconString, builder.log)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing lexicon file " + lexiconPath + " (" + lastResult.String() + ")")
		return
	}

	system.generationLexicon.ImportFrom(lexicon)
}

func (builder systemBuilder) ImportGenerationGrammarFromPath(system *system, grammarPath string) {

	grammarString, err := common.ReadFile(grammarPath)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	grammar := builder.parser.CreateGenerationGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + grammarPath + " (" + lastResult.String() + ")")
		return
	}

	system.generationGrammar.ImportFrom(grammar)
}

func (builder systemBuilder) ImportRuleBaseFromPath(solver *central.ProblemSolver, ruleBasePath string) {

	path := common.AbsolutePath(builder.baseDir, ruleBasePath)
	ruleBaseString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading rules " + path + " (" + err.Error() + ")")
		return
	}

	rules := builder.parser.CreateRules(ruleBaseString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing rules file " + path + " (" + lastResult.String() + ")")
		return
	}

	solver.AddRuleBase(knowledge.NewInMemoryRuleBase("rules", rules, builder.log))
}

func (builder systemBuilder) ImportInMemoryFactBase(name string, solver *central.ProblemSolver, factBase relationSetFactBase, matcher *mentalese.RelationMatcher) {

	path := common.AbsolutePath(builder.baseDir, factBase.Facts)
	factString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading facts " + path + " (" + err.Error() + ")")
		return
	}

	facts := builder.parser.CreateRelationSet(factString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing facts file " + path + " (" + lastResult.String() + ")")
		return
	}

	path = common.AbsolutePath(builder.baseDir, factBase.ReadMap)
	readMapString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading read map file " + path + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateTransformations(readMapString)
	lastResult = builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing read map file " + path + " (" + lastResult.String() + ")")
		return
	}

	dbMapWrite := []mentalese.RelationTransformation{}
	path = common.AbsolutePath(builder.baseDir, factBase.WriteMap)
	if path != "" {
		writeMapString, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading write map file " + path + " (" + err.Error() + ")")
			return
		}

		dbMapWrite = builder.parser.CreateTransformations(writeMapString)
		lastResult = builder.parser.GetLastParseResult()
		if !lastResult.Ok {
			builder.log.AddError("Error parsing write map file " + path + " (" + lastResult.String() + ")")
			return
		}
	}

	stats, _ := builder.CreateDbStats(factBase.Stats)

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	solver.AddFactBase(knowledge.NewInMemoryFactBase(name, facts, matcher, dbMap, dbMapWrite, stats, entities, builder.log))
}

func (builder systemBuilder) ImportMySqlDatabase(name string, solver *central.ProblemSolver, nameResolver *central.NameResolver, factBase mysqlFactBase, matcher *mentalese.RelationMatcher) {

	path := common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading map file " + path + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateTransformations(mapString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing map file " + path + " (" + lastResult.String() + ")")
		return
	}

	stats, ok := builder.CreateDbStats(factBase.Stats)
	if !ok {
		return
	}

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	database := knowledge.NewMySqlFactBase(name, factBase.Domain, factBase.Username, factBase.Password, factBase.Database, matcher, dbMap, stats, entities, builder.log)

	for _, table := range factBase.Tables {
		columns := []string{}
		for _, column := range table.Columns {
			columns = append(columns, column.Name)
		}
		database.AddTableDescription(table.Name, columns)
	}

	if factBase.Enabled {
		solver.AddFactBase(database)
	}
}

func (builder systemBuilder) ImportSparqlDatabase(name string, solver *central.ProblemSolver, nameResolver *central.NameResolver, factBase sparqlFactBase, matcher *mentalese.RelationMatcher) {

	mapPath := common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(mapPath)
	if err != nil {
		builder.log.AddError("Error reading map file " + mapPath + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateTransformations(mapString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing map file " + mapPath + " (" + lastResult.String() + ")")
		return
	}

	names, ok := builder.CreateConfigMap(factBase.Names)
	if !ok {
		return
	}

	stats, ok := builder.CreateDbStats(factBase.Stats)
	if !ok {
		return
	}

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	doCache := factBase.DoCache

	database := knowledge.NewSparqlFactBase(name, factBase.Baseurl, factBase.Defaultgraphuri, matcher, dbMap, names, stats, entities, doCache, builder.log)

	solver.AddFactBase(database)
}

func (builder systemBuilder) CreateConfigMap(path string) (mentalese.ConfigMap, bool) {

	configMap := mentalese.ConfigMap{}
	absolutePath := common.AbsolutePath(builder.baseDir, path)

	content, err := common.ReadFile(absolutePath)
	if err != nil {
		builder.log.AddError("Error reading config map file " + absolutePath + " (" + err.Error() + ")")
		return configMap, false
	}

	err = json.Unmarshal([]byte(content), &configMap)
	if err != nil {
		builder.log.AddError("Error parsing config map file " + absolutePath + " (" + err.Error() + ")")
		return configMap, false
	}

	return configMap, true
}

func (builder systemBuilder) CreateDbStats(path string) (mentalese.DbStats, bool) {

	stats := mentalese.DbStats{}

	if path != "" {

		absolutePath := common.AbsolutePath(builder.baseDir, path)

		content, err := common.ReadFile(absolutePath)
		if err != nil {
			builder.log.AddError("Error reading db stats file " + absolutePath + " (" + err.Error() + ")")
			return stats, false
		}

		err = json.Unmarshal([]byte(content), &stats)
		if err != nil {
			builder.log.AddError("Error parsing db stats file " + absolutePath + " (" + err.Error() + ")")
			return stats, false
		}
	}

	return stats, true
}

func (builder systemBuilder) CreateEntities(path string) (mentalese.Entities, bool) {

	entities := mentalese.Entities{}

	if path != "" {

		absolutePath := common.AbsolutePath(builder.baseDir, path)

		content, err := common.ReadFile(absolutePath)
		if err != nil {
			builder.log.AddError("Error reading entities file " + absolutePath + " (" + err.Error() + ")")
			return entities, false
		}

		entityStructure := Entities{}

		err = json.Unmarshal([]byte(content), &entityStructure)
		if err != nil {
			builder.log.AddError("Error parsing entities file " + absolutePath + " (" + err.Error() + ")")
			return entities, false
		}

		for key, entityInfo := range entityStructure {

			nameRelationSet := builder.parser.CreateRelationSet(entityInfo.Name)

			parseResult := builder.parser.GetLastParseResult()
			if !parseResult.Ok {
				builder.log.AddError("Error parsing " + path + " (" + parseResult.String() + ")")
				return entities, false
			}

			knownBy := map[string]mentalese.RelationSet{}
			for knownByKey, knownByValue := range entityInfo.Knownby {
				knownBy[knownByKey] = builder.parser.CreateRelationSet(knownByValue)

				parseResult := builder.parser.GetLastParseResult()
				if !parseResult.Ok {
					builder.log.AddError("Error parsing " + path + " (" + parseResult.String() + ")")
					return entities, false
				}
			}

			entities[key] = mentalese.EntityInfo{
				Name: nameRelationSet,
				Knownby: knownBy,
			}
		}
	}

	return entities, true
}

func (builder systemBuilder) ImportSolutionBaseFromPath(system *system, solutionBasePath string) {

	path := common.AbsolutePath(builder.baseDir, solutionBasePath)
	solutionString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading solutions file " + path + " (" + err.Error() + ")")
		return
	}

	solutions := builder.parser.CreateSolutions(solutionString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing solutions file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.answerer.AddSolutions(solutions)
}

func (builder systemBuilder) ImportGeneric2DsTransformations(system *system, transformationsPath string) {

	path := common.AbsolutePath(builder.baseDir, transformationsPath)
	transformationstring, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading transformations file " + path + " (" + err.Error() + ")")
		return
	}

	transformations := builder.parser.CreateTransformations(transformationstring)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing transformations file " + path + " (" + lastResult.String() + ")")
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
		builder.log.AddError("Error reading transformations file " + path + " (" + err.Error() + ")")
		return
	}

	transformations := builder.parser.CreateTransformations(transformationstring)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing transformations file " + path + " (" + lastResult.String() + ")")
		return
	}

	for _, transformation := range transformations {
		system.ds2generic = append(system.ds2generic, transformation)
	}
}

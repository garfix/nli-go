package global

import (
	"encoding/json"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/nested"
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

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", builder.log)
	matcher := mentalese.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)

	system.grammar = parse.NewGrammar()
	system.generationGrammar = parse.NewGrammar()
	system.tokenizer = parse.NewTokenizer(builder.log)
	system.relationizer = earley.NewRelationizer(builder.log)
	system.dialogContext = central.NewDialogContext()

	predicates, _ := builder.CreatePredicates(config.Predicates)

	solver := central.NewProblemSolver(matcher, predicates, system.dialogContext, builder.log)

	solver.AddFunctionBase(systemFunctionBase)

	systemAggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", builder.log)
	solver.AddMultipleBindingsBase(systemAggregateBase)

	nestedStructureBase := nested.NewSystemNestedStructureBase(solver, system.dialogContext, predicates, builder.log)
	solver.AddNestedStructureBase(nestedStructureBase)

	shellBase := knowledge.NewShellBase("shell", builder.log)
	solver.AddFunctionBase(shellBase)

	system.dialogContextStorage = NewDialogContextFileStorage(builder.log)
	system.nameResolver = central.NewNameResolver(solver, matcher, predicates, builder.log, system.dialogContext)
	system.parser = earley.NewParser(system.grammar, system.nameResolver, predicates, builder.log)
	system.answerer = central.NewAnswerer(matcher, solver, builder.log)
	system.generator = generate.NewGenerator(system.generationGrammar, builder.log, matcher)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)

	for _, grammarPath := range config.Grammars {
		builder.ImportGrammarFromPath(system, grammarPath)
	}
	for _, grammarPath := range config.Generationgrammars {
		builder.ImportGenerationGrammarFromPath(system, grammarPath)
	}
	for _, ruleBasePath := range config.Rulebases {
		builder.ImportRuleBaseFromPath(solver, ruleBasePath)
	}
	for _, factBase := range config.Factbases.Relation {
		builder.ImportInMemoryFactBase(factBase.Name, solver, factBase, matcher)
	}
	for _, factBase := range config.Factbases.Mysql {
		builder.ImportMySqlDatabase(factBase.Database, solver, system.nameResolver, factBase, matcher)
	}
	for _, factBase := range config.Factbases.Sparql {
		builder.ImportSparqlDatabase(factBase.Name, solver, predicates, factBase, matcher)
	}
	for _, solutionBasePath := range config.Solutions {
		builder.ImportSolutionBaseFromPath(system, solutionBasePath)
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

func (builder systemBuilder) ImportGrammarFromPath(system *system, grammarPath string) {

	path := common.AbsolutePath(builder.baseDir, grammarPath)
	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	grammar := builder.parser.CreateGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.grammar.ImportFrom(grammar)
}

func (builder systemBuilder) ImportGenerationGrammarFromPath(system *system, grammarPath string) {

	path := common.AbsolutePath(builder.baseDir, grammarPath)
	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	grammar := builder.parser.CreateGenerationGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + path + " (" + lastResult.String() + ")")
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

	dbMap := builder.parser.CreateRules(readMapString)
	lastResult = builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing read map file " + path + " (" + lastResult.String() + ")")
		return
	}

	dbMapWrite := []mentalese.Rule{}
	path = common.AbsolutePath(builder.baseDir, factBase.WriteMap)
	if path != "" {
		writeMapString, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading write map file " + path + " (" + err.Error() + ")")
			return
		}

		dbMapWrite = builder.parser.CreateRules(writeMapString)
		lastResult = builder.parser.GetLastParseResult()
		if !lastResult.Ok {
			builder.log.AddError("Error parsing write map file " + path + " (" + lastResult.String() + ")")
			return
		}
	}

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	database := knowledge.NewInMemoryFactBase(name, facts, matcher, dbMap, dbMapWrite, entities, builder.log)

	if factBase.SharedIds != "" {
		sharedIds, ok := builder.LoadSharedIds(factBase.SharedIds)
		if ok {
			database.SetSharedIds(sharedIds)
		}
	}

	solver.AddFactBase(database)
}

func (builder systemBuilder) LoadSharedIds(path string) (knowledge.SharedIds, bool) {

	sharedIds := knowledge.SharedIds{}

	if path != "" {

		absolutePath := common.AbsolutePath(builder.baseDir, path)

		content, err := common.ReadFile(absolutePath)
		if err != nil {
			builder.log.AddError("Error reading shared ids file " + absolutePath + " (" + err.Error() + ")")
			return sharedIds, false
		}

		err = json.Unmarshal([]byte(content), &sharedIds)
		if err != nil {
			builder.log.AddError("Error parsing shared ids file " + absolutePath + " (" + err.Error() + ")")
			return sharedIds, false
		}
	}

	return sharedIds, true
}

func (builder systemBuilder) ImportMySqlDatabase(name string, solver *central.ProblemSolver, nameResolver *central.NameResolver, factBase mysqlFactBase, matcher *mentalese.RelationMatcher) {

	path := common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading map file " + path + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateRules(mapString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing map file " + path + " (" + lastResult.String() + ")")
		return
	}

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	database := knowledge.NewMySqlFactBase(name, factBase.Domain, factBase.Username, factBase.Password, factBase.Database, matcher, dbMap, entities, builder.log)

	for _, table := range factBase.Tables {
		columns := []string{}
		for _, column := range table.Columns {
			columns = append(columns, column.Name)
		}
		database.AddTableDescription(table.Name, columns)
	}

	if factBase.SharedIds != "" {
		sharedIds, ok := builder.LoadSharedIds(factBase.SharedIds)
		if ok {
			database.SetSharedIds(sharedIds)
		}
	}

	if factBase.Enabled {
		solver.AddFactBase(database)
	}
}

func (builder systemBuilder) ImportSparqlDatabase(name string, solver *central.ProblemSolver, predicates mentalese.Predicates, factBase sparqlFactBase, matcher *mentalese.RelationMatcher) {

	mapPath := common.AbsolutePath(builder.baseDir, factBase.Map)
	mapString, err := common.ReadFile(mapPath)
	if err != nil {
		builder.log.AddError("Error reading map file " + mapPath + " (" + err.Error() + ")")
		return
	}

	dbMap := builder.parser.CreateRules(mapString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing map file " + mapPath + " (" + lastResult.String() + ")")
		return
	}

	names, ok := builder.CreateConfigMap(factBase.Names)
	if !ok {
		return
	}

	entities, ok := builder.CreateEntities(factBase.Entities)
	if !ok {
		return
	}

	doCache := factBase.DoCache

	database := knowledge.NewSparqlFactBase(name, factBase.Baseurl, factBase.Defaultgraphuri, matcher, dbMap, names, entities, predicates, doCache, builder.log)

	if factBase.SharedIds != "" {
		sharedIds, ok := builder.LoadSharedIds(factBase.SharedIds)
		if ok {
			database.SetSharedIds(sharedIds)
		}
	}

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

			nameRelationSet := builder.parser.CreateRelation(entityInfo.Name)

			parseResult := builder.parser.GetLastParseResult()
			if !parseResult.Ok {
				builder.log.AddError("Error parsing " + path + " (" + parseResult.String() + ")")
				return entities, false
			}

			knownBy := map[string]mentalese.Relation{}
			for knownByKey, knownByValue := range entityInfo.Knownby {
				knownBy[knownByKey] = builder.parser.CreateRelation(knownByValue)

				parseResult := builder.parser.GetLastParseResult()
				if !parseResult.Ok {
					builder.log.AddError("Error parsing " + path + " (" + parseResult.String() + ")")
					return entities, false
				}
			}

			entities[key] = mentalese.EntityInfo{
				Name:    nameRelationSet,
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

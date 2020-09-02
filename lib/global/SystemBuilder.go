package global

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/nested"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/earley"
	"path/filepath"
	"strings"
)

type systemBuilder struct {
	log                *common.SystemLog
	baseDir            string
	parser             *importer.InternalGrammarParser
	loadedModules      []string
	applicationAliases map[string]string
}

func NewSystem(systemPath string, log *common.SystemLog) *system {

	system := &system{ log: log }

	absolutePath, err := filepath.Abs(systemPath)
	if err != nil {
		log.AddError(err.Error())
		return system
	}

	builder := newSystemBuilder(absolutePath, log)
	builder.build(system)

	return system
}

func newSystemBuilder(baseDir string, log *common.SystemLog) *systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return &systemBuilder {
		baseDir: baseDir,
		parser:  parser,
		log:     log,
	}
}

func (builder *systemBuilder) build(system *system) {

	indexes, ok := builder.readIndexes()
	if !ok {
		return
	}

	config, ok := builder.readConfig()
	if !ok {
		return
	}

	builder.buildBasic(config, system)

	builder.applicationAliases = map[string]string{}
	for alias, moduleSpec := range config.Uses {
		parts := strings.Split(moduleSpec, ":")
		moduleName := parts[0]
		builder.applicationAliases[moduleName] = alias
	}

	builder.loadedModules = []string{ "go" }
	for alias, moduleSpec := range config.Uses {
		builder.loadModule(moduleSpec, alias, &indexes, system)
	}
}

func (builder *systemBuilder) buildBasic(config config, system *system) {

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", builder.log)
	matcher := mentalese.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)
	system.matcher = matcher

	system.grammar = parse.NewGrammar()
	system.generationGrammar = parse.NewGrammar()
	system.tokenizer = builder.createTokenizer(config.Tokenizer)
	system.relationizer = earley.NewRelationizer(builder.log)
	system.dialogContext = central.NewDialogContext()

	predicates, _ := builder.CreatePredicates(config.Predicates)

	system.predicates = &predicates

	solver := central.NewProblemSolver(matcher, predicates, system.dialogContext, builder.log)

	solver.AddFunctionBase(systemFunctionBase)

	systemAggregateBase := knowledge.NewSystemAggregateBase("system-aggregate", builder.log)
	solver.AddMultipleBindingsBase(systemAggregateBase)

	nestedStructureBase := nested.NewSystemNestedStructureBase(solver, system.dialogContext, predicates, builder.log)
	solver.AddNestedStructureBase(nestedStructureBase)

	shellBase := knowledge.NewShellBase("shell", builder.log)
	solver.AddFunctionBase(shellBase)

	system.solver = solver
	system.dialogContextStorage = NewDialogContextFileStorage(builder.log)
	system.nameResolver = central.NewNameResolver(solver, matcher, predicates, builder.log, system.dialogContext)
	system.parser = earley.NewParser(system.grammar, system.nameResolver, predicates, builder.log)
	system.answerer = central.NewAnswerer(matcher, solver, builder.log)
	system.generator = generate.NewGenerator(system.generationGrammar, builder.log, matcher)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)
}

func (builder *systemBuilder) CreatePredicates(path string) (mentalese.Predicates, bool) {

	predicates := mentalese.Predicates{}

	if path != "" {

		absolutePath := common.AbsolutePath(builder.baseDir, path)

		content, err := common.ReadFile(absolutePath)
		if err != nil {
			builder.log.AddError("Error reading predicates file " + absolutePath + " (" + err.Error() + ")")
			return predicates, false
		}

		rawPredicates := mentalese.Predicates{}
		err = json.Unmarshal([]byte(content), &rawPredicates)
		if err != nil {
			builder.log.AddError("Error parsing predicates file " + absolutePath + " (" + err.Error() + ")")
			return predicates, false
		}

		// hack: fix it by replacing the entire JSON file with something better
		for predicate, entityTypes := range rawPredicates {
			predicate = strings.Replace(predicate, ":", "_", 1)
			predicates[predicate] = entityTypes
		}
	}

	return predicates, true
}

func (builder *systemBuilder) createTokenizer(configExpression string) *parse.Tokenizer {

	expression := parse.DefaultTokenizerExpression

	if configExpression != "" {
		expression = configExpression
	}

	return parse.NewTokenizer(expression, builder.log)
}

func (builder *systemBuilder) loadModule(moduleSpec string, alias string, indexes *map[string]index, system *system) {

	parts := strings.Split(moduleSpec, ":")
	if len(parts) != 2 {
		builder.log.AddError("A module specification must have a module name and a version: module-name:1.2.3")
		return
	}

	moduleName := parts[0]
	version := parts[1]

	// check if the module has been loaded already
	for _, aModuleName := range builder.loadedModules {
		if aModuleName == moduleName {
			// no need to load again. also: avoid circular dependencies
			return
		}
	}

	builder.loadedModules = append(builder.loadedModules, moduleName)

	index, found := (*indexes)[moduleName]
	if !found {
		builder.log.AddError("Module not found: " + moduleName)
		return
	}

	if !builder.checkVersion(moduleName, version, index.Version) {
		return
	}

	aliasMap := builder.createAliasMap(index, moduleName)

	moduleBaseDir := builder.baseDir + "/" + moduleName
	applicationAlias := builder.applicationAliases[moduleName]
	builder.processIndex(index, system, applicationAlias, moduleBaseDir, aliasMap)

	builder.loadDependentModules(index, indexes, system)
}

func (builder *systemBuilder) checkVersion(moduleName string, expectedVersion string, actualVersion string) bool {
	// elementary version check
	if expectedVersion != actualVersion {
		builder.log.AddError("Module " + moduleName + " has version " + actualVersion + ", but version " + expectedVersion + " is required")
		return false
	}

	return true
}

func (builder *systemBuilder) createAliasMap(index index, moduleName string) map[string]string {

	aliasMap := map[string]string{
		"": builder.GetApplicationAlias(moduleName),
		"go": "go",
	}

	for moduleAlias, moduleSpec := range index.Uses {
		parts := strings.Split(moduleSpec, ":")
		applicationAlias := builder.GetApplicationAlias(parts[0])
		aliasMap[moduleAlias] = applicationAlias
	}

	return aliasMap
}

func (builder *systemBuilder) loadDependentModules(index index, indexes *map[string]index, system *system) {
	for _, moduleSpec := range index.Uses {
		builder.loadModule(moduleSpec, "", indexes, system)
	}
}

func (builder *systemBuilder) readConfig() (config, bool) {

	config := config{}
	configPath := builder.baseDir + "/config.yml"
	configYml, err := common.ReadFile(configPath)
	if err != nil {
		builder.log.AddError(err.Error())
		return config, false
	}

	err = yaml.Unmarshal([]byte(configYml), &config)
	if err != nil {
		builder.log.AddError("Error parsing YAML file " + configPath + " (" + err.Error() + ")")
		return config, false
	}

	return config, true
}

func (builder *systemBuilder) readIndexes() (map[string]index, bool) {

	indexes := map[string]index{}
	ok := true

	files, err := ioutil.ReadDir(builder.baseDir)
	if err != nil {
		builder.log.AddError(err.Error())
		ok = false
		goto end
	}

	for _, fileInfo := range files {
		if !fileInfo.IsDir() { continue }

		dirName := fileInfo.Name()
		indexPath := builder.baseDir + "/" + dirName + "/index.yml"

		indexYml, err := common.ReadFile(indexPath)
		if err != nil {
			builder.log.AddError(err.Error())
			ok = false
			goto end
		}

		index := index{ }
		err = yaml.Unmarshal([]byte(indexYml), &index)
		if err != nil {
			builder.log.AddError("Error parsing YAML file " + indexPath + " (" + err.Error() + ")")
			ok = false
			goto end
		}

		if index.Type == "" {
			builder.log.AddError("'type' is required; index.yml from: " + dirName)
			goto end
		}

		if index.Version == "" {
			builder.log.AddError("'version' is required; index.yml from: " + dirName)
			goto end
		}

		indexes[dirName] = index
	}

	end:

	return indexes, ok
}

func (builder *systemBuilder) GetApplicationAlias(module string) string {

	alias, found := builder.applicationAliases[module]

	if !found {
		builder.log.AddError("Module not found: " + module)
		alias = ""
	}

	return alias
}

func (builder *systemBuilder) processIndex(index index, system *system, applicationAlias string, moduleBaseDir string, aliasMap map[string]string) bool {

	ok := true

	builder.parser.SetAliasMap(aliasMap)

	switch index.Type {
	case "domain":
		builder.buildDomain(index, system, moduleBaseDir)
	case "grammar":
		builder.buildGrammar(index, system, moduleBaseDir)
	case "solution":
		builder.buildSolution(index, system, moduleBaseDir)
	case "db/internal":
		builder.buildInternalDatabase(index, system, moduleBaseDir, applicationAlias)
	case "db/sparql":
		builder.buildSparqlDatabase(index, system, moduleBaseDir, applicationAlias)
	case "db/mysql":
		builder.buildMySqlDatabase(index, system, moduleBaseDir, applicationAlias)
	default:
		builder.log.AddError("Unknown type: " + index.Type)
		ok = false
	}

	return ok
}

func (builder *systemBuilder) buildDomain(index index, system *system, moduleBaseDir string) {
	for _, rule := range index.Rules {
		builder.importRuleBaseFromPath(system, moduleBaseDir + "/" + rule)
	}
}

func (builder *systemBuilder) buildGrammar(index index, system *system, moduleBaseDir string) {

	for _, read := range index.Read {
		builder.importGrammarFromPath(system, moduleBaseDir + "/" + read)
	}
	for _, write := range index.Write {
		builder.importGenerationGrammarFromPath(system, moduleBaseDir + "/" + write)
	}
}

func (builder *systemBuilder) buildSolution(index index, system *system, moduleBaseDir string) {

	for _, solution := range index.Solution {
		builder.importSolutionBaseFromPath(system, moduleBaseDir + "/" + solution)
	}
}

func (builder *systemBuilder) buildInternalDatabase(index index, system *system, baseDir string, applicationAlias string) {

	facts := mentalese.RelationSet{}

	for _, file := range index.Facts {
		path := common.AbsolutePath(baseDir, file)
		factString, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading facts " + path + " (" + err.Error() + ")")
			return
		}

		facts = builder.parser.CreateRelationSet(factString)
		lastResult := builder.parser.GetLastParseResult()
		if !lastResult.Ok {
			builder.log.AddError("Error parsing facts file " + path + " (" + lastResult.String() + ")")
			return
		}
	}

	readMap := builder.buildReadMap(index, baseDir)
	writeMap := builder.buildWriteMap(index, baseDir)
	entities := builder.buildEntities(index, baseDir)

	database := knowledge.NewInMemoryFactBase(applicationAlias, facts, system.matcher, readMap, writeMap, entities, builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solver.AddFactBase(database)
}

func (builder *systemBuilder) buildSparqlDatabase(index index, system *system, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	entities := builder.buildEntities(index, baseDir)
	names, ok := builder.buildNames(index, baseDir, applicationAlias)
	if !ok {
		return
	}

	database := knowledge.NewSparqlFactBase(applicationAlias, index.BaseUrl, index.DefaultGraphUri, system.matcher, readMap, names, entities, *system.predicates, index.Cache, builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solver.AddFactBase(database)
}

func (builder *systemBuilder) buildMySqlDatabase(index index, system *system, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	entities := builder.buildEntities(index, baseDir)

	prefix := ""

	if applicationAlias != "" {
		prefix = applicationAlias + "_"
	}

	database := knowledge.NewMySqlFactBase(applicationAlias, index.Username, index.Password, index.Database, system.matcher, readMap, entities, builder.log)

	for _, table := range index.Tables {
		columns := []string{}
		for _, column := range table.Columns {
			columns = append(columns, column.Name)
		}
		database.AddTableDescription(prefix + table.Name, table.Name, columns)
	}

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solver.AddFactBase(database)
}

func (builder *systemBuilder) buildReadMap(index index, baseDir string) []mentalese.Rule {

	readMap := []mentalese.Rule{}

	for _, file := range index.Read {
		path := common.AbsolutePath(baseDir, file)
		readMapString, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading read map file " + path + " (" + err.Error() + ")")
			return readMap
		}

		readMap = append(readMap, builder.parser.CreateRules(readMapString)...)
		lastResult := builder.parser.GetLastParseResult()
		if !lastResult.Ok {
			builder.log.AddError("Error parsing read map file " + path + " (" + lastResult.String() + ")")
			return readMap
		}
	}

	return readMap
}

func (builder *systemBuilder) buildWriteMap(index index, baseDir string) []mentalese.Rule {

	writeMap := []mentalese.Rule{}

	for _, file := range index.Write {
		path := common.AbsolutePath(baseDir, file)
		if path != "" {
			writeMapString, err := common.ReadFile(path)
			if err != nil {
				builder.log.AddError("Error reading write map file " + path + " (" + err.Error() + ")")
				return writeMap
			}

			writeMap = append(writeMap, builder.parser.CreateRules(writeMapString)...)
			lastResult := builder.parser.GetLastParseResult()
			if !lastResult.Ok {
				builder.log.AddError("Error parsing write map file " + path + " (" + lastResult.String() + ")")
				return writeMap
			}
		}
	}

	return writeMap
}

func (builder *systemBuilder) buildEntities(index index, baseDir string) mentalese.Entities {

	entities := mentalese.Entities{}

	for _, file := range index.Entities {
		path := common.AbsolutePath(baseDir, file)
		newEntities, ok := builder.CreateEntities(path)
		if !ok {
			return entities
		}
		for key, value := range newEntities {
			entities[key] = value
		}
	}

	return entities
}

func (builder *systemBuilder) buildSharedIds(index index, baseDir string) (knowledge.SharedIds, bool) {

	sharedIds := knowledge.SharedIds{}
	ok := true

	for _, file := range index.Shared {
		path := common.AbsolutePath(baseDir, file)
		sharedIds, ok = builder.LoadSharedIds(path)
		if !ok {
			break
		}
	}

	return sharedIds, ok
}

func (builder systemBuilder) buildNames(index index, baseDir string, applicationAlias string) (mentalese.ConfigMap, bool) {

	configMap := mentalese.ConfigMap{}
	names := mentalese.ConfigMap{}

	path := common.AbsolutePath(baseDir, index.Names)
	content, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading config map file " + path + " (" + err.Error() + ")")
		return configMap, false
	}

	err = json.Unmarshal([]byte(content), &configMap)
	if err != nil {
		builder.log.AddError("Error parsing config map file " + path + " (" + err.Error() + ")")
		return configMap, false
	}

	for key, value := range configMap {
		names[applicationAlias + "_" + key] = value
	}

	return names, true
}

func (builder *systemBuilder) importGrammarFromPath(system *system, path string) {

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

func (builder *systemBuilder) importGenerationGrammarFromPath(system *system, path string) {

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

func (builder *systemBuilder) importSolutionBaseFromPath(system *system, path string) {

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

func (builder systemBuilder) importRuleBaseFromPath(system *system, path string) {

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

	system.solver.AddRuleBase(knowledge.NewInMemoryRuleBase("rules", rules, builder.log))
}

func (builder systemBuilder) CreateEntities(path string) (mentalese.Entities, bool) {

	entities := mentalese.Entities{}

	if path != "" {

		content, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading entities file " + path + " (" + err.Error() + ")")
			return entities, false
		}

		entityStructure := Entities{}

		err = json.Unmarshal([]byte(content), &entityStructure)
		if err != nil {
			builder.log.AddError("Error parsing entities file " + path + " (" + err.Error() + ")")
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

func (builder systemBuilder) LoadSharedIds(path string) (knowledge.SharedIds, bool) {

	sharedIds := knowledge.SharedIds{}

	if path != "" {

		content, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading shared ids file " + path + " (" + err.Error() + ")")
			return sharedIds, false
		}

		err = json.Unmarshal([]byte(content), &sharedIds)
		if err != nil {
			builder.log.AddError("Error parsing shared ids file " + path + " (" + err.Error() + ")")
			return sharedIds, false
		}
	}

	return sharedIds, true
}
package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"nli-go/lib/central"
	"nli-go/lib/central/goal"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/function"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/parse/morphology"
	"path/filepath"
	"regexp"
	"strings"
)

type systemBuilder struct {
	log                *common.SystemLog
	baseDir            string
	sessionId		   string
	varDir			   string
	parser             *importer.InternalGrammarParser
	loadedModules      []string
	applicationAliases map[string]string
}

// systemDir: absolute base dir of the interaction system
// outputDir: absolute base dir of all output files
// log: the log file that will store progress and errors
func NewSystem(systemDir string, sessionId string, varDir string, log *common.SystemLog) *System {

	system := &System{ log: log }

	absolutePath, err := filepath.Abs(systemDir)
	if err != nil {
		log.AddError(err.Error())
		return system
	}

	builder := newSystemBuilder(absolutePath, sessionId, varDir, log)
	builder.build(system)

	return system
}

func newSystemBuilder(baseDir string, sessionId string, varDir string, log *common.SystemLog) *systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return &systemBuilder {
		baseDir: baseDir,
		sessionId: sessionId,
		varDir:  varDir,
		parser:  parser,
		log:     log,
	}
}

func (builder *systemBuilder) getSystemName() string {
	return filepath.Base(builder.baseDir)
}

func (builder *systemBuilder) build(system *System) {

	indexes, ok := builder.readIndexes()
	if !ok {
		return
	}

	config, ok := builder.readConfig()
	if !ok {
		return
	}

	builder.buildBasic(system)

	builder.applicationAliases = map[string]string{}
	for alias, moduleSpec := range config.Uses {
		parts := strings.Split(moduleSpec, ":")
		moduleName := parts[0]
		builder.applicationAliases[moduleName] = alias
	}

	builder.loadedModules = []string{ "go" }
	for _, moduleSpec := range config.Uses {
		builder.loadModule(moduleSpec, &indexes, system)
	}

	languageBase := knowledge.NewLanguageBase("language", system.grammars, system.meta, system.dialogContext, system.nameResolver, system.answerer, system.generator, builder.log)
	system.solverAsync.AddSolverFunctionBase(languageBase)

	system.solverAsync.Reindex()
}

func (builder *systemBuilder) buildBasic(system *System) {

	system.meta = mentalese.NewMeta()

	systemFunctionBase := knowledge.NewSystemFunctionBase("System-function", system.meta, builder.log)
	matcher := central.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)
	system.matcher = matcher

	path := builder.varDir + "/" + builder.getSystemName() + "/session"
	storage := common.NewFileStorage(path, builder.sessionId, common.StorageSession, "session", system.log)

	system.grammars = []parse.Grammar{}
	system.relationizer = parse.NewRelationizer(builder.log)
	system.internalGrammarParser = builder.parser
	system.processList = goal.NewProcessList()

	solverAsync := central.NewProblemSolverAsync(matcher, builder.log)
	solverAsync.AddFunctionBase(systemFunctionBase)

	systemMultiBindingBase := knowledge.NewSystemMultiBindingBase("System-aggregate", builder.log)
	solverAsync.AddMultipleBindingBase(systemMultiBindingBase)

	anaphoraQueue := central.NewAnaphoraQueue()
	deicticCenter := central.NewDeicticCenter()
	system.dialogContext = central.NewDialogContext(storage, anaphoraQueue, deicticCenter, system.processList)
	nestedStructureBase := function.NewSystemSolverFunctionBase(anaphoraQueue, deicticCenter, system.meta, builder.log)
	solverAsync.AddSolverFunctionBase(nestedStructureBase)

	system.solverAsync = solverAsync
	system.processRunner = central.NewProcessRunner(solverAsync, builder.log)

	system.nameResolver = central.NewNameResolver(solverAsync, system.meta, builder.log)
	system.answerer = central.NewAnswerer(matcher, builder.log)
	system.generator = generate.NewGenerator(builder.log, matcher)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)

	domainIndex, ok := builder.buildIndex(common.Dir() + "/../base/domain")
	if ok {
		builder.buildDomain(domainIndex, system, common.Dir() + "/../base/domain", "go")
	}
	dbIndex, ok := builder.buildIndex(common.Dir() + "/../base/db")
	if ok {
		builder.buildInternalDatabase(dbIndex, system, common.Dir() + "/../base/db", "nligo-db")
	}
}

func (builder *systemBuilder) AddPredicates(path string, system *System) bool {

	if path != "" {

		content, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading relation types file " + path + " (" + err.Error() + ")")
			return false
		}

		relationTypes := builder.parser.CreateRelationSet(content)

		for _, relationType := range relationTypes {
			name := strings.Replace(relationType.Predicate, ":", "_", 1)
			sorts := []string{}
			for _, argument := range relationType.Arguments {
				sorts = append(sorts, argument.TermValue)
			}

			system.meta.AddPredicate(name, sorts)
		}
	}

	return true
}

func (builder *systemBuilder) AddSubSorts(path string, system *System) bool {

	if path != "" {

		content, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading relation types file " + path + " (" + err.Error() + ")")
			return false
		}

		sortRelations := builder.parser.CreateSortRelations(content)

		for _, sortRelation := range sortRelations {
			system.meta.AddSubSort(sortRelation.GetSuperSort(), sortRelation.GetSubSort())
		}
	}

	return true
}

func (builder *systemBuilder) loadModule(moduleSpec string, indexes *map[string]index, system *System) {

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

func (builder *systemBuilder) loadDependentModules(index index, indexes *map[string]index, system *System) {
	for _, moduleSpec := range index.Uses {
		builder.loadModule(moduleSpec, indexes, system)
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
		anIndex := index{}
		anIndex, ok = builder.buildIndex(builder.baseDir + "/" + dirName)
		if ! ok {
			goto end
		}

		indexes[dirName] = anIndex
	}

	end:

	return indexes, ok
}

func (builder *systemBuilder) buildIndex(dirName string) (index, bool) {

	index := index{ }
	indexPath := dirName + "/index.yml"

	ok := true

	indexYml, err := common.ReadFile(indexPath)
	if err != nil {
		builder.log.AddError(err.Error())
		ok = false
		goto end
	}

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

	end:
		return index, ok
}

func (builder *systemBuilder) GetApplicationAlias(module string) string {

	alias, found := builder.applicationAliases[module]

	if !found {
		builder.log.AddError("Module not found: " + module)
		alias = ""
	}

	return alias
}

func (builder *systemBuilder) processIndex(index index, system *System, applicationAlias string, moduleBaseDir string, aliasMap map[string]string) bool {

	ok := true

	builder.parser.SetAliasMap(aliasMap)

	switch index.Type {
	case "domain":
		builder.buildDomain(index, system, moduleBaseDir, applicationAlias)
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

func (builder *systemBuilder) buildDomain(index index, system *System, moduleBaseDir string, applicationAlias string) {

	ok := true
	path := ""

	builder.importRuleBaseFromPath(index, system, moduleBaseDir, applicationAlias)

	path = common.AbsolutePath(moduleBaseDir, index.Sorts)
	ok = builder.AddSorts(path, system)
	if !ok {
		return
	}

	path = common.AbsolutePath(moduleBaseDir, index.Predicates)
	ok = builder.AddPredicates(path, system)
	if !ok {
		return
	}

	path = common.AbsolutePath(moduleBaseDir, index.Subsorts)
	ok = builder.AddSubSorts(path, system)
	if !ok {
		return
	}
}

func (builder *systemBuilder) buildGrammar(index index, system *System, moduleBaseDir string) {

	grammar := parse.NewGrammar(index.Locale)

	for _, read := range index.Read {
		builder.importGrammarFromPath(&grammar, moduleBaseDir + "/" + read)
	}
	for _, write := range index.Write {
		builder.importGenerationGrammarFromPath(&grammar, moduleBaseDir + "/" + write)
	}

	if index.TokenExpression != "" {
		grammar.SetTokenizer(parse.NewTokenizer(index.TokenExpression))
	}

	if index.Text != "" {
		grammar.SetTexts(builder.importTexts(moduleBaseDir + "/" + index.Text))
	}

	if index.Morphology != nil {
		grammar.SetMorphologicalAnalyzer(builder.importMorphologicalAnalyzer(index.Morphology, system, moduleBaseDir))
	}

	system.grammars = append(system.grammars, grammar)
}

func (builder *systemBuilder) importTexts(textFile string) map[string]string {
	texts := map[string]string{}
	csvString, err := common.ReadFile(textFile)
	if err != nil {
		builder.log.AddError(err.Error())
		return texts
	}

	expression, _ := regexp.Compile("((?:\\\\,|[^,])+),((?:\\\\,|[^,])+)(?:\n|$)")

	lines := expression.FindAllStringSubmatch(csvString, -1)

	for _, parts := range lines {
		source := parts[1]
		translation := parts[2]

		source = strings.ReplaceAll(source, "\\", "")
		translation = strings.ReplaceAll(translation, "\\", "")

		source = strings.Trim(source, " \t")
		translation = strings.Trim(translation, " \t")

		texts[source] = translation
	}

	return texts
}

func (builder *systemBuilder) importMorphologicalAnalyzer(parts map[string]string, system *System, moduleBaseDir string) *parse.MorphologicalAnalyzer {

	parsingRules := mentalese.NewGrammarRules()
	segmentationRules := morphology.NewSegmentationRules()

	segmenterPath, found := parts["segmentation"]
	if found {
		segmentationRules = builder.readSegmentationRulesFromPath(moduleBaseDir + "/" + segmenterPath)
	}

	parsingPath, found := parts["parsing"]
	if found {
		parsingRules = builder.readGrammarFromPath(moduleBaseDir + "/" + parsingPath)
	}

	segmenter := morphology.NewSegmenter(segmentationRules)

	return parse.NewMorphologicalAnalyzer(
		parsingRules,
		segmenter,
		parse.NewParser(parsingRules, system.log),
		parse.NewRelationizer(system.log),
		system.log)
}

func (builder *systemBuilder) readSegmentationRulesFromPath(path string) *morphology.SegmentationRules {

	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return morphology.NewSegmentationRules()
	}

	rules := builder.parser.CreateSegmentationRules(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing segmentation rule file " + path + " (" + lastResult.String() + ")")
		return morphology.NewSegmentationRules()
	}

	return rules
}

func (builder *systemBuilder) readGrammarFromPath(path string) *mentalese.GrammarRules {

	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return mentalese.NewGrammarRules()
	}

	rules := builder.parser.CreateGrammarRules(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + path + " (" + lastResult.String() + ")")
		return mentalese.NewGrammarRules()
	}

	return rules
}

func (builder *systemBuilder) importGrammarFromPath(grammar *parse.Grammar, path string) {

	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	rules := builder.parser.CreateGrammarRules(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + path + " (" + lastResult.String() + ")")
		return
	}

	grammar.GetReadRules().ImportFrom(rules)
}

func (builder *systemBuilder) importGenerationGrammarFromPath(grammar *parse.Grammar, path string) {

	grammarString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError(err.Error())
		return
	}

	rules := builder.parser.CreateGenerationGrammar(grammarString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing grammar file " + path + " (" + lastResult.String() + ")")
		return
	}

	grammar.GetWriteRules().ImportFrom(rules)
}

func (builder *systemBuilder) buildSolution(index index, system *System, moduleBaseDir string) {

	for _, solution := range index.Solution {
		builder.importSolutionBaseFromPath(system, moduleBaseDir + "/" + solution)
	}
}

func (builder *systemBuilder) buildInternalDatabase(index index, system *System, baseDir string, applicationAlias string) {

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

	storageType := index.StorageType

	if storageType == "" {
		storageType = common.StorageNone
	}

	path := builder.varDir + "/" + builder.getSystemName() + "/database"
	storage := common.NewFileStorage(path, builder.sessionId, storageType, applicationAlias, system.log)
	database := knowledge.NewInMemoryFactBase(applicationAlias, facts, system.matcher, readMap, writeMap, storage, builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solverAsync.AddFactBase(database)
}

func (builder *systemBuilder) buildSparqlDatabase(index index, system *System, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	names, ok := builder.buildNames(index, baseDir, applicationAlias)
	if !ok {
		return
	}

	database := knowledge.NewSparqlFactBase(applicationAlias, index.BaseUrl, index.DefaultGraphUri, system.matcher, readMap, names, index.Cache, builder.varDir + "/sparql-cache", builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solverAsync.AddFactBase(database)
}

func (builder *systemBuilder) buildMySqlDatabase(index index, system *System, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	writeMap := builder.buildWriteMap(index, baseDir)

	prefix := ""

	if applicationAlias != "" {
		prefix = applicationAlias + "_"
	}

	database := knowledge.NewMySqlFactBase(applicationAlias, index.Username, index.Password, index.Database, system.matcher, readMap, writeMap,  builder.log)

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

	system.solverAsync.AddFactBase(database)
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

func (builder *systemBuilder) AddSorts(path string, system *System) bool {

	sorts, ok := builder.CreateSorts(path)
	if !ok {
		return false
	}
	for name, info := range sorts {
		system.meta.AddSortInfo(name, info)
	}

	return true
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

	err = yaml.Unmarshal([]byte(content), &configMap)
	if err != nil {
		builder.log.AddError("Error parsing config map file " + path + " (" + err.Error() + ")")
		return configMap, false
	}

	for key, value := range configMap {
		names[applicationAlias + "_" + key] = value
	}

	return names, true
}

func (builder *systemBuilder) importSolutionBaseFromPath(system *System, path string) {

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

func (builder systemBuilder) importRuleBaseFromPath(index index, system *System, baseDir string, applicationAlias string) {

	rules := []mentalese.Rule{}

	for _, rule := range index.Rules {
		path := baseDir + "/" + rule
		ruleBaseString, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading rules " + path + " (" + err.Error() + ")")
			return
		}

		rules = append(rules, builder.parser.CreateRules(ruleBaseString)...)
		lastResult := builder.parser.GetLastParseResult()
		if !lastResult.Ok {
			builder.log.AddError("Error parsing rules file " + path + " (" + lastResult.String() + ")")
			return
		}
	}

	writeList, ok := builder.readWritelist(index, baseDir, applicationAlias)
	if !ok {
		return
	}

	storageType := index.StorageType

	if storageType == "" {
		storageType = common.StorageNone
	}

	path := builder.varDir + "/" + builder.getSystemName() + "/rules"
	storage := common.NewFileStorage(path, builder.sessionId, storageType, applicationAlias, system.log)

	system.solverAsync.AddRuleBase(knowledge.NewInMemoryRuleBase("rules", rules, writeList, storage, builder.log))
}

func (builder systemBuilder) readWritelist(index index, baseDir string, applicationAlias string) ([]string, bool) {

	writelist := []string{}

	for _, write := range index.Write {
		path := baseDir + "/" + write
		predicatesList, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading rules " + path + " (" + err.Error() + ")")
			return writelist, false
		}

		aWritelist := []string{}
		err = yaml.Unmarshal([]byte(predicatesList), &aWritelist)
		if err != nil {
			return writelist, false
		}
		for _, predicate := range aWritelist {
			aliasedPredicate := applicationAlias + "_" + predicate
			writelist = append(writelist, aliasedPredicate)
		}
	}

	return writelist, true
}

func (builder systemBuilder) CreateSorts(path string) (mentalese.Entities, bool) {

	entities := mentalese.Entities{}

	if path != "" {

		content, err := common.ReadFile(path)
		if err != nil {
			builder.log.AddError("Error reading entities file " + path + " (" + err.Error() + ")")
			return entities, false
		}

		entityStructure := Entities{}

		err = yaml.Unmarshal([]byte(content), &entityStructure)
		if err != nil {
			builder.log.AddError("Error parsing entities file " + path + " (" + err.Error() + ")")
			return entities, false
		}

		for key, entityInfo := range entityStructure {

			nameRelation := builder.parser.CreateRelation(entityInfo.Name)
			EntityRelationSet := mentalese.RelationSet{}
			if entityInfo.Entity != "" {
				EntityRelationSet = builder.parser.CreateRelationSet(entityInfo.Entity)
			}

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

			entities[key] = mentalese.SortInfo{
				Name:    nameRelation,
				Knownby: knownBy,
				Entity:  EntityRelationSet,
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

		err = yaml.Unmarshal([]byte(content), &sharedIds)
		if err != nil {
			builder.log.AddError("Error parsing shared ids file " + path + " (" + err.Error() + ")")
			return sharedIds, false
		}
	}

	return sharedIds, true
}
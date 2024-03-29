package global

import (
	"io/ioutil"
	"nli-go/lib/central"
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

	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v2"
)

type systemBuilder struct {
	log                *common.SystemLog
	appDir             string
	workDir            string
	sessionId          string
	parser             *importer.InternalGrammarParser
	loadedModules      []string
	applicationAliases map[string]string
	conn               *websocket.Conn
}

func NewSystem(appDir string, workDir string, sessionId string, log *common.SystemLog, conn *websocket.Conn) *System {

	system := &System{log: log}

	builder := newSystemBuilder(appDir, workDir, sessionId, log, conn)
	builder.build(system)

	return system
}

func newSystemBuilder(appDir string, workDir string, sessionId string, log *common.SystemLog, conn *websocket.Conn) *systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	logListener := func(production common.LogMessage) {

		response := mentalese.Response{
			Resource:    central.NO_RESOURCE,
			MessageType: mentalese.MessageLog,
			Message:     production,
		}

		websocket.JSON.Send(conn, response)
	}

	log.AddListener(logListener)

	return &systemBuilder{
		appDir:    appDir,
		workDir:   workDir,
		sessionId: sessionId,
		parser:    parser,
		log:       log,
		conn:      conn,
	}
}

func (builder *systemBuilder) getSystemName() string {
	return filepath.Base(builder.appDir)
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

	builder.loadedModules = []string{"go"}
	for _, moduleSpec := range config.Uses {
		builder.loadModule(moduleSpec, &indexes, system)
	}

	languageBase := knowledge.NewLanguageBase(
		"language",
		system.grammars,
		system.relationizer,
		system.meta,
		system.dialogContext,
		system.nameResolver,
		system.answerer,
		system.generator,
		system.clientConnector,
		builder.log)
	system.solver.AddSolverFunctionBase(languageBase)

	system.solver.Reindex()
}

func (builder *systemBuilder) buildBasic(system *System) {

	system.meta = mentalese.NewMeta()
	system.clientConnector = system.CreatClientConnector(builder.conn)

	systemFunctionBase := knowledge.NewSystemFunctionBase("System-function", system.meta, builder.log)
	matcher := central.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)
	system.matcher = matcher

	variableGenerator := mentalese.NewVariableGenerator()
	system.variableGenerator = variableGenerator

	system.grammars = []parse.Grammar{}
	system.relationizer = parse.NewRelationizer(variableGenerator, builder.log)
	system.internalGrammarParser = builder.parser
	system.processList = central.NewProcessList()

	modifier := central.NewFactBaseModifier(builder.log, variableGenerator)

	solverAsync := central.NewProblemSolver(matcher, variableGenerator, builder.log)
	solverAsync.AddFunctionBase(systemFunctionBase)
	solverAsync.SetModifier(modifier)

	systemMultiBindingBase := knowledge.NewSystemMultiBindingBase("System-aggregate", builder.log)
	solverAsync.AddMultipleBindingBase(systemMultiBindingBase)

	system.dialogContext = central.NewDialogContext(variableGenerator)
	nestedStructureBase := function.NewSystemSolverFunctionBase(system.dialogContext, system.meta, builder.log, system.clientConnector)
	solverAsync.AddSolverFunctionBase(nestedStructureBase)
	system.solver = solverAsync

	system.processRunner = central.NewProcessRunner(system.processList, solverAsync, builder.log)

	callback := func() {
		if system.processList.IsEmpty() {
			system.clientConnector.SendToClient(central.NO_RESOURCE, "processlist_clear", nil)
		}
	}

	system.processList.AddListener(callback)

	system.nameResolver = central.NewNameResolver(solverAsync, system.meta, builder.log)
	system.answerer = central.NewAnswerer(matcher, builder.log)

	generationState := mentalese.NewGenerationState()
	generationMatcher := central.NewRelationMatcher(builder.log)
	generationFunctionBase := knowledge.NewGenerationFunctionBase("generation", generationState, builder.log)
	generationMatcher.AddFunctionBase(systemFunctionBase)
	generationMatcher.AddFunctionBase(generationFunctionBase)

	system.generator = generate.NewGenerator(builder.log, generationMatcher, generationState)
	system.surfacer = generate.NewSurfaceRepresentation(builder.log)

	domainIndex, ok := builder.buildIndex(common.Dir() + "/../base/domain")
	if ok {
		builder.buildDomain(domainIndex, system, common.Dir()+"/../base/domain", "go")
	}
	dbIndex, ok := builder.buildIndex(common.Dir() + "/../base/db")
	if ok {
		builder.buildInternalDatabase(dbIndex, system, common.Dir()+"/../base/db", "nligo-db")
	}
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

	moduleBaseDir := builder.appDir + "/" + moduleName
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
		"":   builder.GetApplicationAlias(moduleName),
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
	configPath := builder.appDir + "/config.yml"
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

	files, err := ioutil.ReadDir(builder.appDir)
	if err != nil {
		builder.log.AddError(err.Error())
		ok = false
		goto end
	}

	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			continue
		}

		dirName := fileInfo.Name()
		if dirName == "test" {
			continue
		}

		anIndex := index{}
		anIndex, ok = builder.buildIndex(builder.appDir + "/" + dirName)
		if !ok {
			goto end
		}

		indexes[dirName] = anIndex
	}

end:

	return indexes, ok
}

func (builder *systemBuilder) buildIndex(dirName string) (index, bool) {

	index := index{}
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
	case "intent":
		builder.buildIntent(index, system, moduleBaseDir)
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
}

func (builder *systemBuilder) buildGrammar(index index, system *System, moduleBaseDir string) {

	grammar := parse.NewGrammar(index.Locale)

	for _, read := range index.Read {
		builder.importGrammarFromPath(&grammar, moduleBaseDir+"/"+read)
	}
	for _, write := range index.Write {
		builder.importGenerationGrammarFromPath(&grammar, moduleBaseDir+"/"+write)
	}

	if index.TokenExpression != "" {
		grammar.SetTokenizer(parse.NewTokenizer(index.TokenExpression))
	}

	if index.Text != "" {
		grammar.SetTexts(builder.importTexts(moduleBaseDir + "/" + index.Text))
	}

	if index.Morphology != nil {
		grammar.SetMorphologicalAnalyzer(builder.importMorphologicalAnalyzer(index.Morphology, system, moduleBaseDir, grammar.GetReadRules()))
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

func (builder *systemBuilder) importMorphologicalAnalyzer(parts map[string]string, system *System, moduleBaseDir string, readRules *mentalese.GrammarRules) *parse.MorphologicalAnalyzer {

	segmentationRules := morphology.NewSegmentationRules()

	segmenterPath, found := parts["segmentation"]
	if found {
		segmentationRules = builder.readSegmentationRulesFromPath(moduleBaseDir + "/" + segmenterPath)
	}

	parsingPath, found := parts["parsing"]
	if found {
		parsingRules := builder.readGrammarFromPath(moduleBaseDir + "/" + parsingPath)
		readRules.Merge(parsingRules)
	}

	segmenter := morphology.NewSegmenter(segmentationRules, readRules)

	return parse.NewMorphologicalAnalyzer(
		readRules,
		segmenter,
		parse.NewParser(readRules, system.log),
		parse.NewDialogizer(system.variableGenerator),
		parse.NewRelationizer(system.variableGenerator, system.log),
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

func (builder *systemBuilder) buildIntent(index index, system *System, moduleBaseDir string) {

	for _, intent := range index.Intent {
		builder.importIntentBaseFromPath(system, moduleBaseDir+"/"+intent)
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

	database := knowledge.NewInMemoryFactBase(applicationAlias, facts, system.matcher, readMap, writeMap, builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solver.AddFactBase(database)
}

func (builder *systemBuilder) buildSparqlDatabase(index index, system *System, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	names, ok := builder.buildNames(index, baseDir, applicationAlias)
	if !ok {
		return
	}

	database := knowledge.NewSparqlFactBase(applicationAlias, index.BaseUrl, index.DefaultGraphUri, system.matcher, readMap, names, index.Cache, builder.workDir+"/sparql-cache", builder.log)

	sharedIds, ok := builder.buildSharedIds(index, baseDir)
	if ok {
		database.SetSharedIds(sharedIds)
	}

	system.solver.AddFactBase(database)
}

func (builder *systemBuilder) buildMySqlDatabase(index index, system *System, baseDir string, applicationAlias string) {

	readMap := builder.buildReadMap(index, baseDir)
	writeMap := builder.buildWriteMap(index, baseDir)

	prefix := ""

	if applicationAlias != "" {
		prefix = applicationAlias + "_"
	}

	database := knowledge.NewMySqlFactBase(applicationAlias, index.Username, index.Password, index.Database, system.matcher, readMap, writeMap, builder.log)

	for _, table := range index.Tables {
		columns := []string{}
		for _, column := range table.Columns {
			columns = append(columns, column.Name)
		}
		database.AddTableDescription(prefix+table.Name, table.Name, columns)
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
		names[applicationAlias+"_"+key] = value
	}

	return names, true
}

func (builder *systemBuilder) importIntentBaseFromPath(system *System, path string) {

	intentString, err := common.ReadFile(path)
	if err != nil {
		builder.log.AddError("Error reading intents file " + path + " (" + err.Error() + ")")
		return
	}

	intents := builder.parser.CreateIntent(intentString)
	lastResult := builder.parser.GetLastParseResult()
	if !lastResult.Ok {
		builder.log.AddError("Error parsing intents file " + path + " (" + lastResult.String() + ")")
		return
	}

	system.answerer.AddIntents(intents)
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

	system.solver.AddRuleBase(knowledge.NewInMemoryRuleBase("rules", rules, writeList, builder.log))
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

func (builder systemBuilder) CreateSorts(path string) (mentalese.SortProperties, bool) {

	entities := mentalese.SortProperties{}

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
			GenderRelation := mentalese.Relation{}
			if entityInfo.Gender != "" {
				GenderRelation = builder.parser.CreateRelation(entityInfo.Gender)
			}
			NumberRelation := mentalese.Relation{}
			if entityInfo.Number != "" {
				NumberRelation = builder.parser.CreateRelation(entityInfo.Number)
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

			entities[key] = mentalese.SortProperty{
				Name:    nameRelation,
				Gender:  GenderRelation,
				Number:  NumberRelation,
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

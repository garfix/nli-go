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
	log     *common.SystemLog
	baseDir string
	parser  *importer.InternalGrammarParser
	usesStack []string
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

func newSystemBuilder(baseDir string, log *common.SystemLog) systemBuilder {

	parser := importer.NewInternalGrammarParser()
	parser.SetPanicOnParseFail(false)

	return systemBuilder {
		baseDir: baseDir,
		parser:  parser,
		log:     log,
	}
}

func (builder systemBuilder) build(system *system) {

	indexes, ok := builder.readIndexes()
	if !ok {
		return
	}

	config, ok := builder.readConfig()
	if !ok {
		return
	}

	builder.buildBasic(config, system)

	for _, moduleSpec := range config.Modules {
		builder.use(moduleSpec, &indexes, system)
	}
}

func (builder systemBuilder) buildBasic(config config, system *system) {

	systemFunctionBase := knowledge.NewSystemFunctionBase("system-function", builder.log)
	matcher := mentalese.NewRelationMatcher(builder.log)
	matcher.AddFunctionBase(systemFunctionBase)

	system.grammar = parse.NewGrammar()
	system.generationGrammar = parse.NewGrammar()
	system.tokenizer = builder.createTokenizer(config.Tokenizer)
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

func (builder systemBuilder) createTokenizer(configExpression string) *parse.Tokenizer {

	expression := parse.DefaultTokenizerExpression

	if configExpression != "" {
		expression = configExpression
	}

	return parse.NewTokenizer(expression, builder.log)
}

func (builder systemBuilder) use(moduleSpec string, indexes *map[string]index, system *system) {

	parts := strings.Split(moduleSpec, ":")
	if len(parts) != 2 {
		builder.log.AddError("A uses specification must have a module internalName and a version: module-internalName:1.2.3")
		return
	}

	moduleName := parts[0]
	version := parts[1]

	index, found := (*indexes)[moduleName]
	if !found {
		builder.log.AddError("Module not found: " + moduleName)
		return
	}

	if !builder.checkVersion(moduleName, version, index.Version) {
		return
	}

	builder.processDependendentIndexes(index, indexes, moduleName, system)

	moduleBaseDir := builder.baseDir + "/" + moduleName
	if !builder.processIndex(index, system, moduleBaseDir) {
		return
	}
}

func (builder systemBuilder) checkVersion(moduleName string, expectedVersion string, actualVersion string) bool {
	// elementary version check
	if expectedVersion != actualVersion {
		builder.log.AddError("Module " + moduleName + " has version " + actualVersion + ", but version " + expectedVersion + " is required")
		return false
	}

	return true
}

func (builder systemBuilder) processDependendentIndexes(index index, indexes *map[string]index, moduleName string, system *system) {

	// check for recursion
	for _, aModuleName := range builder.usesStack {
		if aModuleName == moduleName {
			// circular dependencies are allowed
			return
		}
	}

	builder.usesStack = append(builder.usesStack, moduleName)

	// load dependencies
	for _, moduleSpec := range index.Uses {
		builder.use(moduleSpec, indexes, system)
	}

	builder.usesStack = builder.usesStack[0: len(builder.usesStack) - 1]
}

func (builder systemBuilder) readConfig() (config, bool) {

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

func (builder systemBuilder) readIndexes() (map[string]index, bool) {

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

		//if index.Name != "" && index.Name != dirName {
		//	builder.log.AddError("The internalName in an index.yml file is optional; but if it is given, it must be the same as the directory internalName: " + dirName)
		//	goto end
		//}

		if index.Type == "" {
			builder.log.AddError("'type' is required; index.yml from: " + dirName)
			goto end
		}

		indexes[dirName] = index
	}

	end:

	return indexes, ok
}

func (builder systemBuilder) processIndex(index index, system *system, moduleBaseDir string) bool {

	ok := true

	switch index.Type {
	case "domain":
		builder.buildDomain(index, system, moduleBaseDir)
	case "grammar":
		builder.buildGrammar(index, system, moduleBaseDir)
	case "solution":
		builder.buildSolution(index, system, moduleBaseDir)
	default:
		builder.log.AddError("Unknown type: " + index.Type)
		ok = false
	}

	return ok
}

func (builder systemBuilder) buildDomain(index index, system *system, moduleBaseDir string) {

}

func (builder systemBuilder) buildGrammar(index index, system *system, moduleBaseDir string) {

	for _, read := range index.Read {
		builder.importGrammarFromPath(system, moduleBaseDir + "/" + read)
	}
	for _, write := range index.Write {
		builder.importGenerationGrammarFromPath(system, moduleBaseDir + "/" + write)
	}
}

func (builder systemBuilder) buildSolution(index index, system *system, moduleBaseDir string) {

	for _, solution := range index.Solution {
		builder.ImportSolutionBaseFromPath(system, moduleBaseDir + "/" + solution)
	}
}

func (builder systemBuilder) importGrammarFromPath(system *system, path string) {

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

func (builder systemBuilder) importGenerationGrammarFromPath(system *system, path string) {

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

func (builder systemBuilder) ImportSolutionBaseFromPath(system *system, path string) {

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
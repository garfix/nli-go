package knowledge

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
)

type LanguageBase struct {
	KnowledgeBaseCore
	matcher 			  *central.RelationMatcher
	grammars              []parse.Grammar
	meta                  *mentalese.Meta
	dialogContext         *central.DialogContext
	nameResolver          *central.NameResolver
	answerer 			  *central.Answerer
	log 			      *common.SystemLog
}

func NewLanguageBase(
	name string,
	grammars []parse.Grammar,
	meta *mentalese.Meta,
	dialogContext *central.DialogContext,
	nameResolver *central.NameResolver,
	answerer *central.Answerer,
	log *common.SystemLog) *LanguageBase {
	return &LanguageBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ name },
		matcher: central.NewRelationMatcher(log),
		grammars: grammars,
		meta: meta,
		dialogContext: dialogContext,
		nameResolver: nameResolver,
		answerer: answerer,
		log: log,
	}
}

func (base *LanguageBase) GetFunctions() map[string]api.SolverFunction {
	return map[string]api.SolverFunction{
		mentalese.PredicateLocale: base.locale,
		mentalese.PredicateTokenize: base.tokenize,
		mentalese.PredicateParse: base.parse,
		mentalese.PredicateRelationize: base.relationize,
		mentalese.PredicateGenerate:     base.generate,
		mentalese.PredicateSurface:      base.surface,
		mentalese.PredicateFindSolution: base.findSolution,
		mentalese.PredicateSolve: base.solve,
		mentalese.PredicateFindResponse: base.findResponse,
		mentalese.PredicateCreateAnswer: base.createAnswer,
	}
}

func (base *LanguageBase) locale(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	localeVar := input.Arguments[0].TermValue

	newBindings := mentalese.NewBindingSet()

	for _, grammar := range base.grammars {
		newBinding := binding.Copy()
		newBinding.Set(localeVar, mentalese.NewTermString(grammar.GetLocale()))
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *LanguageBase) tokenize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return mentalese.NewBindingSet()
	}

	locale := bound.Arguments[0].TermValue
	rawInput := bound.Arguments[1].TermValue
	tokenVar := input.Arguments[2].TermValue

	grammar, found := base.getGrammar(locale)
	if !found {
		return mentalese.NewBindingSet()
	}

	tokens := grammar.GetTokenizer().Process(rawInput)

	terms := []mentalese.Term{}
	for _, token := range tokens {
		terms = append(terms, mentalese.NewTermString(token))
	}

	newBinding := binding.Copy()
	newBinding.Set(tokenVar, mentalese.NewTermList(terms))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) getGrammar(locale string) (parse.Grammar, bool) {

	grammar := parse.Grammar{}

	found := false
	for _, aGrammar := range base.grammars {
		if aGrammar.GetLocale() == locale {
			grammar = aGrammar
			found = true
		}
	}

	return grammar, found
}

func (base *LanguageBase) parse(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "slv", base.log) {
		return mentalese.NewBindingSet()
	}

	locale := bound.Arguments[0].TermValue
	tokenList := bound.Arguments[1].TermValueList
	sentenceVar := input.Arguments[2].TermValue

	tokens := []string{}
	for _, token := range tokenList {
		tokens = append(tokens, token.TermValue)
	}

	grammar, found := base.getGrammar(locale)
	if !found {
		return mentalese.NewBindingSet()
	}

	parser := parse.NewParser(grammar.GetReadRules(), base.log)
	parser.SetMorphologicalAnalyzer(grammar.GetMorphologicalAnalyzer())
	parseTrees := parser.Parse(tokens, "s", []string{"S"})

	newBindings := mentalese.NewBindingSet()
	for _, parseTree := range parseTrees {
		newBinding := binding.Copy()
		newBinding.Set(sentenceVar, mentalese.NewTermJson(parseTree))
		newBindings.Add(newBinding)
	}

	base.log.AddProduction("Parse trees found", strconv.Itoa(len(parseTrees)))

	return newBindings
}

func (base *LanguageBase) relationize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jvv", base.log) {
		return mentalese.NewBindingSet()
	}

	senseVar := input.Arguments[1].TermValue
	requestBindingVar := input.Arguments[2].TermValue
	var parseTree parse.ParseTreeNode

	bound.Arguments[0].GetJsonValue(&parseTree)
	relationizer := parse.NewRelationizer(base.log)
	sortFinder := central.NewSortFinder(base.meta)

	requestRelations, names := relationizer.Relationize(parseTree, []string{ "S"})

	// extract sorts: variable => sort
	sorts, sortFound := sortFinder.FindSorts(requestRelations)
	if !sortFound {
		// conflicting sorts
		base.log.AddProduction("Break", "Breaking due to conflicting sorts: " + sorts.String())
		return mentalese.NewBindingSet()
	}

	entityIds, nameNotFound := base.findNames(names, sorts)


// todo
nameNotFound = nameNotFound


	// names found and linked to id
	for _, value := range entityIds.GetAll() {
		base.dialogContext.AnaphoraQueue.AddReferenceGroup(
			central.EntityReferenceGroup{ central.CreateEntityReference(value.TermValue, value.TermSort) })
	}
	base.log.AddProduction("Named entities", entityIds.String())

	newBinding := binding.Copy()

	newBinding.Set(senseVar, mentalese.NewTermRelationSet(requestRelations))
	newBinding.Set(requestBindingVar, mentalese.NewTermJson(entityIds))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) findNames(names mentalese.Binding, sorts mentalese.Sorts) (mentalese.Binding, string) {

	entityIds := mentalese.NewBinding()
	nameNotFound := ""

	// look up entity ids by name
	entityIds = mentalese.NewBinding()
	for variable, name := range names.GetAll() {

		// find sort
		sort, found := sorts[variable]
		if !found {
			base.log.AddProduction("Info",
				"The name '" + name.TermValue + "' could not be looked up because no sort could be derived from the relations.")
			if nameNotFound == "" {
				nameNotFound = name.TermValue
			}
			goto next
		}

		// find name information
		nameInformations := base.nameResolver.ResolveName(name.TermValue, sort)
		if len(nameInformations) == 0 {
			base.log.AddProduction("Info",
				"Database lookup for name '" + name.TermValue + "'  with sort '" + sort + "' did not give any results")
			nameNotFound = name.TermValue
			goto next
		}

		// make the user choose one entity from multiple with the same name
		if len(nameInformations) > 1 {
			nameInformations = base.nameResolver.Resolve(nameInformations)
		}

		// link variable to ID
		for _, nameInformation := range nameInformations {
			entityIds.Set(variable, mentalese.NewTermId(nameInformation.SharedId, nameInformation.EntityType))
		}
	}

	next:

	return entityIds, nameNotFound
}

func (base *LanguageBase) findSolution(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "rv", base.log) {
		return mentalese.NewBindingSet()
	}

	request := bound.Arguments[0].TermValueRelationSet
	solutionVar := input.Arguments[1].TermValue

	solutions := base.answerer.FindSolutions(request)

	newBindings := mentalese.NewBindingSet()

	for _, solution := range solutions {
		newBinding := mentalese.NewBinding()
		newBinding.Set(solutionVar, mentalese.NewTermJson(solution))
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *LanguageBase) solve(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	transformer := central.NewRelationTransformer(base.matcher, base.log)

	//Request, RequestBinding, Solution, ResultBindings
	if !Validate(bound, "rjjv", base.log) {
		return mentalese.NewBindingSet()
	}

	requestBinding := mentalese.Binding{}
	solution := mentalese.Solution{}

	request := bound.Arguments[0].TermValueRelationSet
	bound.Arguments[1].GetJsonValue(&requestBinding)
	bound.Arguments[2].GetJsonValue(&solution)
	resultBindingsVar := input.Arguments[3].TermValue

	child := messenger.GetCursor().GetState("child", 0)
	if child == 0 {

		messenger.GetCursor().SetState("child", 1)

		// apply transformation, if available
		transformedRequest := transformer.Replace(solution.Transformations, request)

		messenger.CreateChildStackFrame(transformedRequest, mentalese.InitBindingSet(requestBinding))
		return mentalese.NewBindingSet()

	} else {

		resultBindings := messenger.GetCursor().GetChildFrameResultBindings()

		newBinding := mentalese.NewBinding()
		newBinding.Set(resultBindingsVar, mentalese.NewTermJson(resultBindings.ToRaw()))

		// queue ids
		group := central.EntityReferenceGroup{}
		for _, id := range resultBindings.GetIds(solution.Result.TermValue) {
			group = append(group, central.CreateEntityReference(id.TermValue, id.TermSort))
		}
		base.dialogContext.AnaphoraQueue.AddReferenceGroup(group)

		return mentalese.InitBindingSet(newBinding)

	}
}

func (base *LanguageBase) findResponse(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jjvv", base.log) {
		return mentalese.NewBindingSet()
	}

	solution := mentalese.Solution{}
	resultBindings := mentalese.NewBindingSet()

	bound.Arguments[0].GetJsonValue(&solution)

	resultBindingsRaw := []map[string]mentalese.Term{}
	bound.Arguments[1].GetJsonValue(&resultBindingsRaw)
	resultBindings.FromRaw(resultBindingsRaw)

	conditionBindingsVar := input.Arguments[2].TermValue
	responseIndexVar := input.Arguments[3].TermValue

	index := messenger.GetCursor().GetState("index", 0)

	// process child results
	if index > 0 {
		conditionBindings := messenger.GetCursor().GetChildFrameResultBindings()
		if !conditionBindings.IsEmpty() {
			newBinding := mentalese.NewBinding()
			newBinding.Set(conditionBindingsVar, mentalese.NewTermJson(conditionBindings.ToRaw()))
			newBinding.Set(responseIndexVar, mentalese.NewTermString(strconv.Itoa(index - 1)))
			return mentalese.InitBindingSet(newBinding)
		}
	}

	if index < len(solution.Responses) {
		response := solution.Responses[index]
		if response.Condition.IsEmpty() {
			newBinding := mentalese.NewBinding()
			newBinding.Set(conditionBindingsVar, mentalese.NewTermJson(resultBindings))
			return mentalese.InitBindingSet(newBinding)
		} else {
			messenger.CreateChildStackFrame(response.Condition, resultBindings)
		}
	}

	messenger.GetCursor().SetState("index", index + 1)

	return mentalese.NewBindingSet()
}

func (base *LanguageBase) createAnswer(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jjiv", base.log) {
		return mentalese.NewBindingSet()
	}

	solution := mentalese.Solution{}
	resultBindings := mentalese.NewBindingSet()

	bound.Arguments[0].GetJsonValue(&solution)

	responseBindingsRaw := []map[string]mentalese.Term{}
	bound.Arguments[1].GetJsonValue(&responseBindingsRaw)
	resultBindings.FromRaw(responseBindingsRaw)

	responseIndex, _ := bound.Arguments[2].GetIntValue()
	answerVar := input.Arguments[3].TermValue

	index := messenger.GetCursor().GetState("index", 0)

	solutionBindings := resultBindings
	resultHandler := solution.Responses[responseIndex]

	if index == 0 {

		messenger.GetCursor().SetState("index", 1)

		if !resultHandler.Preparation.IsEmpty() {
			messenger.CreateChildStackFrame(resultHandler.Preparation, resultBindings)
			return mentalese.NewBindingSet()
		}

	} else {

		solutionBindings = messenger.GetCursor().GetChildFrameResultBindings()

	}

	// create answer relation sets by binding 'answer' to solutionBindings
	answer := base.answerer.Build(resultHandler.Answer, solutionBindings)

	newBinding := mentalese.NewBinding()
	newBinding.Set(answerVar, mentalese.NewTermRelationSet(answer))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) generate(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "srv", base.log) {
		return mentalese.NewBindingSet()
	}

	locale := bound.Arguments[0].TermValue
	answerRelations := bound.Arguments[1].TermValueRelationSet
	tokenVar := input.Arguments[2].TermValue

	grammar, found := base.getGrammar(locale)
	if !found {
		return mentalese.NewBindingSet()
	}

	generator := generate.NewGenerator(base.log, base.matcher)
	tokens := generator.Generate(grammar.GetWriteRules(), answerRelations)

	tokenTerms := []mentalese.Term{}
	for _, token := range tokens {
		tokenTerms = append(tokenTerms, mentalese.NewTermString(token))
	}

	newBinding := binding.Copy()
	newBinding.Set(tokenVar, mentalese.NewTermList(tokenTerms))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) surface(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "lv", base.log) {
		return mentalese.NewBindingSet()
	}

	tokenList := bound.Arguments[0].TermValueList
	surfaceVar := input.Arguments[1].TermValue

	tokens := []string{}
	for _, token := range tokenList {
		tokens = append(tokens, token.TermValue)
	}

	surfacer := generate.NewSurfaceRepresentation(base.log)
	surface := surfacer.Create(tokens)

	newBinding := binding.Copy()
	newBinding.Set(surfaceVar, mentalese.NewTermString(surface))

	return mentalese.InitBindingSet(newBinding)
}
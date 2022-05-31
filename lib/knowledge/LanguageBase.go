package knowledge

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
	"strings"
)

type LanguageBase struct {
	KnowledgeBaseCore
	matcher       *central.RelationMatcher
	grammars      []parse.Grammar
	relationizer  *parse.Relationizer
	meta          *mentalese.Meta
	dialogContext *central.DialogContext
	nameResolver  *central.NameResolver
	answerer      *central.Answerer
	generator     *generate.Generator
	log           *common.SystemLog
}

func NewLanguageBase(
	name string,
	grammars []parse.Grammar,
	relationizer *parse.Relationizer,
	meta *mentalese.Meta,
	dialogContext *central.DialogContext,
	nameResolver *central.NameResolver,
	answerer *central.Answerer,
	generator *generate.Generator,
	log *common.SystemLog) *LanguageBase {
	return &LanguageBase{
		KnowledgeBaseCore: KnowledgeBaseCore{name},
		matcher:           central.NewRelationMatcher(log),
		grammars:          grammars,
		relationizer:      relationizer,
		meta:              meta,
		dialogContext:     dialogContext,
		nameResolver:      nameResolver,
		answerer:          answerer,
		generator:         generator,
		log:               log,
	}
}

func (base *LanguageBase) GetFunctions() map[string]api.SolverFunction {
	return map[string]api.SolverFunction{
		mentalese.PredicateFindLocale:          base.findLocale,
		mentalese.PredicateTokenize:            base.tokenize,
		mentalese.PredicateParse:               base.parse,
		mentalese.PredicateDialogize:           base.dialogize,
		mentalese.PredicateEllipsize:           base.ellipsize,
		mentalese.PredicateRelationize:         base.relationize,
		mentalese.PredicateSortalFiltering:     base.sortalFiltering,
		mentalese.PredicateResolveNames:        base.resolveNames,
		mentalese.PredicateResolveAnaphora:     base.resolveAnaphora,
		mentalese.PredicateExtractRootClauses:  base.extractRootClauses,
		mentalese.PredicateDialogAddRootClause: base.dialogAddRootClause,
		mentalese.PredicateDialogUpdateCenter:  base.dialogUpdateCenter,
		mentalese.PredicateDialogGetCenter:     base.dialogGetCenter,
		mentalese.PredicateGenerate:            base.generate,
		mentalese.PredicateSurface:             base.surface,
		mentalese.PredicateFindSolution:        base.findSolution,
		mentalese.PredicateSolve:               base.solve,
		mentalese.PredicateFindResponse:        base.findResponse,
		mentalese.PredicateCreateAnswer:        base.createAnswer,
		mentalese.PredicateCreateCanned:        base.createCanned,
		mentalese.PredicateTranslate:           base.translate,
	}
}

func (base *LanguageBase) findLocale(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

	base.log.AddProduction("Tokens", strings.Join(tokens, " "))

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

	for _, parseTree := range parseTrees {
		base.log.AddProduction("Parse tree", parseTree.IndentedString(""))
	}

	return newBindings
}

func (base *LanguageBase) dialogize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jv", base.log) {
		return mentalese.NewBindingSet()
	}

	resultVar := input.Arguments[1].TermValue
	var parseTree mentalese.ParseTreeNode

	bound.Arguments[0].GetJsonValue(&parseTree)

	dialogizer := parse.NewDialogizer(base.dialogContext.VariableGenerator)
	newParseTree := dialogizer.Dialogize(&parseTree)

	newBinding := mentalese.NewBinding()
	newBinding.Set(resultVar, mentalese.NewTermJson(newParseTree))

	base.log.AddProduction("Dialogized parse tree", newParseTree.IndentedString(""))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) ellipsize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jv", base.log) {
		return mentalese.NewBindingSet()
	}

	ellipsisVar := input.Arguments[1].TermValue
	var parseTree mentalese.ParseTreeNode

	bound.Arguments[0].GetJsonValue(&parseTree)

	clauses := base.dialogContext.ClauseList.GetRootNodes()
	ellipsizer := parse.NewEllipsizer(clauses, base.log)
	newParseTree, ok := ellipsizer.Ellipsize(parseTree)
	if !ok {
		return mentalese.NewBindingSet()
	}

	newBinding := mentalese.NewBinding()
	newBinding.Set(ellipsisVar, mentalese.NewTermJson(newParseTree))

	base.log.AddProduction("Ellipsized parse tree", newParseTree.IndentedString(""))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) extractRootClauses(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jv", base.log) {
		return mentalese.NewBindingSet()
	}

	rootClauseVar := input.Arguments[1].TermValue

	var parseTree mentalese.ParseTreeNode
	bound.Arguments[0].GetJsonValue(&parseTree)

	rootClauseExtracter := parse.NewRootClauseExtracter()
	rootClauses := rootClauseExtracter.Extract(&parseTree)

	newBindings := mentalese.NewBindingSet()

	for _, rootClause := range rootClauses {
		newBinding := mentalese.NewBinding()
		newBinding.Set(rootClauseVar, mentalese.NewTermJson(rootClause))
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *LanguageBase) dialogAddRootClause(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ja", base.log) {
		return mentalese.NewBindingSet()
	}

	authorIsSystem := input.Arguments[1].TermValue

	var parseTree mentalese.ParseTreeNode
	bound.Arguments[0].GetJsonValue(&parseTree)

	clauseList := base.dialogContext.ClauseList
	entities := mentalese.ExtractEntities(&parseTree)
	clause := mentalese.NewClause(&parseTree, authorIsSystem == "true", entities)
	clauseList.AddClause(clause)

	return mentalese.InitBindingSet(binding)
}

func (base *LanguageBase) dialogUpdateCenter(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	var previousCenter = base.dialogContext.DeicticCenter.GetCenter()
	var center = ""
	var priority = 0

	priorities := map[string]int{
		"previousCenter":              100,
		mentalese.AtomFunctionSubject: 10,
		mentalese.AtomFunctionObject:  5,
	}

	c := base.dialogContext.ClauseList.GetLastClause()

	// new clause has no entities? keep existing center
	if len(c.Functions) == 0 {
		center = previousCenter
	}

	for _, entity := range c.Functions {
		if previousCenter != "" {
			a := getValue(entity.DiscourseVariable, binding)
			b := getValue(previousCenter, binding)
			if a == b {
				priority = priorities["previousCenter"]
				center = entity.DiscourseVariable
				continue
			}
		}
		prio, found := priorities[entity.SyntacticFunction]
		if found {
			if prio > priority {
				priority = prio
				center = entity.DiscourseVariable
			}
		}
	}

	base.dialogContext.DeicticCenter.SetCenter(center)

	return mentalese.InitBindingSet(binding)
}

func getValue(variable string, binding mentalese.Binding) string {
	v, found := binding.Get(variable)
	if found {
		return v.TermValue
	} else {
		return ""
	}
}

func (base *LanguageBase) dialogGetCenter(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	if !Validate(input, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	centerVar := input.Arguments[0].TermValue

	center := mentalese.NewTermAtom("none")
	//clause := base.dialogContext.ClauseList.GetLastClause()
	centerVariable := base.dialogContext.DeicticCenter.GetCenter()
	if centerVariable != "" {
		value, found := base.dialogContext.EntityBindings.Get(centerVariable)
		if found {
			center = value
		}
	}

	newBindings := mentalese.NewBindingSet()
	if center.IsList() {
		for _, item := range center.TermValueList {
			newBinding := mentalese.NewBinding()
			newBinding.Set(centerVar, item)
			newBindings.Add(newBinding)
		}
	} else {
		newBinding := mentalese.NewBinding()
		newBinding.Set(centerVar, center)
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *LanguageBase) relationize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jv", base.log) {
		return mentalese.NewBindingSet()
	}

	var parseTree mentalese.ParseTreeNode
	bound.Arguments[0].GetJsonValue(&parseTree)

	senseVar := input.Arguments[1].TermValue

	requestRelations := base.relationizer.Relationize(parseTree, []string{"S"})

	base.log.AddProduction("Relations", requestRelations.IndentedString(""))

	tags := base.relationizer.ExtractTags(parseTree)
	base.dialogContext.EntityTags.AddTags(tags)

	newBinding := binding.Copy()

	newBinding.Set(senseVar, mentalese.NewTermRelationSet(requestRelations))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) sortalFiltering(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	if !Validate(bound, "r", base.log) {
		return mentalese.NewBindingSet()
	}

	requestRelations := bound.Arguments[0].TermValueRelationSet

	// extract sorts: variable => sort
	sortFinder := central.NewSortFinder(base.meta)
	sorts, sortFound := sortFinder.FindSorts(requestRelations)
	if !sortFound {
		// conflicting sorts
		base.log.AddProduction("Break", "Breaking due to conflicting sorts: "+sorts.String())
		return mentalese.NewBindingSet()
	}

	for variable, sort := range sorts {
		base.dialogContext.EntitySorts.SetSorts(variable, []string{sort})
	}

	return mentalese.InitBindingSet(binding)
}

func (base *LanguageBase) resolveNames(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	if !Validate(bound, "jjvv", base.log) {
		return mentalese.NewBindingSet()
	}

	requestBindingVar := input.Arguments[2].TermValue
	unboundNameVar := input.Arguments[3].TermValue
	var parseTree mentalese.ParseTreeNode

	bound.Arguments[0].GetJsonValue(&parseTree)

	dialogBinding := mentalese.NewBinding()
	dialogBindingsRaw := map[string]mentalese.Term{}
	bound.Arguments[1].GetJsonValue(&dialogBindingsRaw)
	dialogBinding.FromRaw(dialogBindingsRaw)

	names := base.nameResolver.ExtractNames(parseTree, []string{"S"})

	sorts := base.dialogContext.EntitySorts

	entityIds, nameNotFound, genderTags := base.findNames(messenger, names, *sorts)
	base.dialogContext.EntityTags.AddTags(genderTags)

	requestBinding := dialogBinding.Merge(entityIds)

	base.log.AddProduction("Named entities", entityIds.String())

	newBinding := binding.Copy()

	newBinding.Set(requestBindingVar, mentalese.NewTermJson(requestBinding.ToRaw()))
	newBinding.Set(unboundNameVar, mentalese.NewTermString(nameNotFound))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) resolveAnaphora(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	if !Validate(bound, "rjvvv", base.log) {
		return mentalese.NewBindingSet()
	}

	request := bound.Arguments[0].TermValueRelationSet

	dialogBinding := mentalese.NewBinding()
	dialogBindingsRaw := map[string]mentalese.Term{}
	bound.Arguments[1].GetJsonValue(&dialogBindingsRaw)
	dialogBinding.FromRaw(dialogBindingsRaw)

	resolvedRequestVar := input.Arguments[2].TermValue
	resolvedBindingsVar := input.Arguments[3].TermValue
	outputVar := input.Arguments[4].TermValue

	resolver := central.NewAnaphoraResolver(base.dialogContext, base.meta, messenger)
	resolvedRequest, resolvedBindings, output := resolver.Resolve(request, dialogBinding)

	newBinding := mentalese.NewBinding()
	newBinding.Set(resolvedRequestVar, mentalese.NewTermRelationSet(resolvedRequest))
	newBinding.Set(resolvedBindingsVar, mentalese.NewTermJson(resolvedBindings.ToRaw()))
	newBinding.Set(outputVar, mentalese.NewTermString(output))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) findNames(messenger api.ProcessMessenger, names mentalese.Binding, sorts mentalese.EntitySorts) (mentalese.Binding, string, mentalese.RelationSet) {

	entityIds := mentalese.NewBinding()
	nameNotFound := ""
	genderTags := mentalese.RelationSet{}

	// look up entity ids by name
	entityIds = mentalese.NewBinding()
	for variable, name := range names.GetAll() {

		// find sort
		sort, found := sorts[variable]
		if !found {
			base.log.AddProduction("Info",
				"The name '"+name.TermValue+"' could not be looked up because no sort could be derived from the relations.")
			if nameNotFound == "" {
				nameNotFound = name.TermValue
			}
			goto next
		}

		// find name information
		nameInformations := base.nameResolver.ResolveName(name.TermValue, sort[0])
		if len(nameInformations) == 0 {
			base.log.AddProduction("Info",
				"Database lookup for name '"+name.TermValue+"'  with sort '"+sort[0]+"' did not give any results")
			nameNotFound = name.TermValue
			goto next
		}

		// make the user choose one entity from multiple with the same name
		if len(nameInformations) > 1 {
			nameInformations, _ = base.nameResolver.Choose(messenger, nameInformations)
		}

		// link variable to ID
		for _, nameInformation := range nameInformations {
			entityIds.Set(variable, mentalese.NewTermId(nameInformation.SharedId, nameInformation.EntityType))
			if nameInformation.Gender != "" {
				genderTags = append(genderTags, mentalese.NewRelation(false, mentalese.TagAgree, []mentalese.Term{
					mentalese.NewTermVariable(variable),
					mentalese.NewTermAtom(mentalese.AtomGender),
					mentalese.NewTermAtom(nameInformation.Gender),
				}))
			}
		}
	}

next:

	return entityIds, nameNotFound, genderTags
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
	if !Validate(bound, "rjjvvv", base.log) {
		return mentalese.NewBindingSet()
	}

	solution := mentalese.Solution{}

	request := bound.Arguments[0].TermValueRelationSet
	bound.Arguments[2].GetJsonValue(&solution)
	resultBindingsVar := input.Arguments[3].TermValue
	resultCountVar := input.Arguments[4].TermValue
	essentialVar := input.Arguments[5].TermValue

	requestBindings := mentalese.NewBindingSet()
	requestBindingsRaw := []map[string]mentalese.Term{}
	bound.Arguments[1].GetJsonValue(&requestBindingsRaw)
	requestBindings.FromRaw(requestBindingsRaw)

	base.log.AddProduction("Solution", solution.Condition.IndentedString(""))

	// apply transformation, if available
	transformedRequest := transformer.Replace(solution.Transformations, request)

	resultBindings := messenger.ExecuteChildStackFrame(transformedRequest, requestBindings)

	newBinding := mentalese.NewBinding()
	newBinding.Set(resultBindingsVar, mentalese.NewTermJson(resultBindings.ToRaw()))
	newBinding.Set(resultCountVar, mentalese.NewTermString(strconv.Itoa(resultBindings.GetLength())))

	// queue ids
	variable := solution.Result.TermValue

	essential := mentalese.NewBindingSet()
	for _, id := range resultBindings.GetIds(variable) {
		b := mentalese.NewBinding()
		b.Set(variable, id)
		essential.Add(b)
	}

	newBinding.Set(essentialVar, mentalese.NewTermJson(essential.ToRaw()))

	return mentalese.InitBindingSet(newBinding)
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

	responseBindingsVar := input.Arguments[2].TermValue
	responseIndexVar := input.Arguments[3].TermValue

	for index := 0; index < len(solution.Responses); index++ {
		response := solution.Responses[index]
		if response.Condition.IsEmpty() {
			newBinding := mentalese.NewBinding()
			newBinding.Set(responseBindingsVar, mentalese.NewTermJson(resultBindings))
			newBinding.Set(responseIndexVar, mentalese.NewTermString(strconv.Itoa(index)))
			return mentalese.InitBindingSet(newBinding)
		} else {
			responseBindings := messenger.ExecuteChildStackFrame(response.Condition, resultBindings)
			if !responseBindings.IsEmpty() {
				newBinding := mentalese.NewBinding()
				newBinding.Set(responseBindingsVar, mentalese.NewTermJson(responseBindings.ToRaw()))
				newBinding.Set(responseIndexVar, mentalese.NewTermString(strconv.Itoa(index)))
				return mentalese.InitBindingSet(newBinding)
			}
		}
	}

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

	solutionBindings := resultBindings
	resultHandler := solution.Responses[responseIndex]

	solutionBindings = messenger.ExecuteChildStackFrame(resultHandler.Preparation, resultBindings)

	// create answer relation sets by binding 'answer' to solutionBindings
	answer := base.answerer.Build(resultHandler.Answer, solutionBindings)

	base.log.AddProduction("Answer", answer.String())

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

	tokens := base.generator.Generate(grammar.GetWriteRules(), answerRelations)

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

	base.log.AddProduction("Dialog-entities", base.dialogContext.EntityBindings.String())

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) createCanned(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "vas", base.log) {
		return mentalese.NewBindingSet()
	}

	outputVar := input.Arguments[0].TermValue
	templateString := bound.Arguments[1].TermValue
	argumentString := bound.Arguments[2].TermValue

	newBinding := binding.Copy()
	newBinding.Set(outputVar, mentalese.NewTermString(common.GetString(templateString, argumentString)))
	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) translate(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return mentalese.NewBindingSet()
	}

	source := bound.Arguments[0].TermValue
	locale := bound.Arguments[1].TermValue
	translatedVar := input.Arguments[2].TermValue

	grammar, found := base.getGrammar(locale)
	if !found {
		return mentalese.NewBindingSet()
	}

	translation := grammar.GetText(source)

	newBinding := mentalese.NewBinding()
	newBinding.Set(translatedVar, mentalese.NewTermString(translation))
	return mentalese.InitBindingSet(newBinding)
}

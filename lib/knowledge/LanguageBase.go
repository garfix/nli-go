package knowledge

import (
	"encoding/json"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
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
		mentalese.PredicateAnswer: base.answer,
		mentalese.PredicateGenerate: base.generate,
		mentalese.PredicateSurface: base.surface,
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

	return newBindings
}

func (base *LanguageBase) relationize(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "jvv", base.log) {
		return mentalese.NewBindingSet()
	}

	sentenceSerialized := bound.Arguments[0].TermValue
	senseVar := input.Arguments[1].TermValue
	requestBindingVar := input.Arguments[2].TermValue

	var parseTree parse.ParseTreeNode
	jsonBytes := []byte(sentenceSerialized)
	err := json.Unmarshal(jsonBytes, &parseTree)
	if err != nil {
		base.log.AddError(err.Error())
		return mentalese.NewBindingSet()
	}

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

func (base *LanguageBase) answer(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "rjv", base.log) {
		return mentalese.NewBindingSet()
	}

	requestBindings := mentalese.Binding{}

	requestRelations := bound.Arguments[0].TermValueRelationSet
	input.Arguments[1].GetJsonValue(requestBindings)
	answerRelationVar := input.Arguments[2].TermValue

	answerRelations := base.answerer.Answer(messenger, requestRelations, mentalese.InitBindingSet(binding))
	base.log.AddProduction("Answer", answerRelations.String())
	base.log.AddProduction("Anaphora queue", base.dialogContext.AnaphoraQueue.String())

	newBinding := binding.Copy()
	newBinding.Set(answerRelationVar, mentalese.NewTermRelationSet(answerRelations))

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
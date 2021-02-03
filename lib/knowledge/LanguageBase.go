package knowledge

import (
	"encoding/json"
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
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
	}
}

func (base *LanguageBase) locale(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

func (base *LanguageBase) tokenize(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "ssv", base.log) {
		return mentalese.NewBindingSet()
	}

	locale := bound.Arguments[0].TermValue
	rawInput := bound.Arguments[1].TermValue
	tokenVar := input.Arguments[2].TermValue
	tokens := []string{}

	found := false
	for _, grammar := range base.grammars {
		if grammar.GetLocale() == locale {
			tokens = grammar.GetTokenizer().Process(rawInput)
			found = true
		}
	}

	if !found {
		return mentalese.NewBindingSet()
	}

	terms := []mentalese.Term{}
	for _, token := range tokens {
		terms = append(terms, mentalese.NewTermString(token))
	}

	newBinding := binding.Copy()
	newBinding.Set(tokenVar, mentalese.NewTermList(terms))

	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) parse(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

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

	var theGrammar parse.Grammar

	found := false
	for _, grammar := range base.grammars {
		if grammar.GetLocale() == locale {
			theGrammar = grammar
			found = true
		}
	}

	if !found {
		return mentalese.NewBindingSet()
	}

	parser := parse.NewParser(theGrammar.GetReadRules(), base.log)
	parser.SetMorphologicalAnalyzer(theGrammar.GetMorphologicalAnalyzer())
	parseTrees := parser.Parse(tokens, "s", []string{"S"})

	newBindings := mentalese.NewBindingSet()
	for _, parseTree := range parseTrees {

		jsonBytes, err := json.Marshal(parseTree)
		if err != nil {
			return mentalese.BindingSet{}
		}

		jsonString := string(jsonBytes)

		newBinding := binding.Copy()
		newBinding.Set(sentenceVar, mentalese.NewTermString(jsonString))
		newBindings.Add(newBinding)
	}

	return newBindings
}

func (base *LanguageBase) relationize(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "svv", base.log) {
		return mentalese.NewBindingSet()
	}

	sentenceSerialized := bound.Arguments[0].TermValue
	senseVar := input.Arguments[1].TermValue
//	requestBindingVar := input.Arguments[2].TermValue

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

	newBinding := binding.Merge(entityIds)

	newBinding.Set(senseVar, mentalese.NewTermRelationSet(requestRelations))
//	newBinding.Set(requestBindingVar, mentalese.NewTermRelationSet(entityIds))

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

func (base *LanguageBase) answer(input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "rv", base.log) {
		return mentalese.NewBindingSet()
	}

	requestRelations := bound.Arguments[0].TermValueRelationSet
	answerVar := input.Arguments[1].TermValue

	answerRelations := base.answerer.Answer(requestRelations, mentalese.InitBindingSet(binding))
	base.log.AddProduction("Answer", answerRelations.String())
	base.log.AddProduction("Anaphora queue", base.dialogContext.AnaphoraQueue.String())

	newBinding := binding.Copy()
	newBinding.Set(answerVar, mentalese.NewTermRelationSet(answerRelations))

	return mentalese.InitBindingSet(newBinding)
}
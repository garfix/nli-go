package knowledge

import (
	"nli-go/lib/api"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/generate"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
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
		mentalese.PredicateRespond:         base.respond,
		mentalese.PredicateDialogGetCenter: base.dialogGetCenter,
		mentalese.PredicateTranslate:       base.translate,
		mentalese.PredicateDefine:          base.define,
	}
}

func (base *LanguageBase) respond(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	if !Validate(bound, "s", base.log) {
		return mentalese.NewBindingSet()
	}

	rawInput := bound.Arguments[0].TermValue
	output := ""
	tmpOutput := ""
	resolvedRemark := ""
	score := 0
	highestScore := 0

	originalDialogContext := base.dialogContext

	locale := ""
	for _, grammar := range base.grammars {

		i18n := parse.NewI18n(&grammar)

		// set locale
		locale = grammar.GetLocale()
		messenger.SetProcessSlot("locale", mentalese.NewTermString(locale))

		// tokenize
		tokens := grammar.GetTokenizer().Process(rawInput)
		base.log.AddProduction("Tokens", strings.Join(tokens, " "))

		// parse
		parser := parse.NewParser(grammar.GetReadRules(), base.log)
		parser.SetMorphologicalAnalyzer(grammar.GetMorphologicalAnalyzer())
		parseTrees, parseResult := parser.Parse(tokens, "s", []string{"S"})
		if parseResult.Error != "" {
			output = i18n.TranslateWithParam(parseResult.Error, parseResult.ErrorArg)
			continue
		}

		for _, parseTree := range parseTrees {

			base.log.AddProduction("Parse tree", parseTree.IndentedString(""))

			// fork the dialog context
			// and make sure it is available to application predicates
			base.dialogContext = originalDialogContext.Fork()

			// work with the components of the dialog context explicitly, to expose dependencies
			clauseList := base.dialogContext.ClauseList
			deicticCenter := base.dialogContext.DeicticCenter
			entityBindings := base.dialogContext.EntityBindings
			entityTags := base.dialogContext.EntityTags
			entitySorts := base.dialogContext.EntitySorts
			entityLabels := base.dialogContext.EntityLabels
			entityDefinitions := base.dialogContext.EntityDefinitions

			// slowly forget about named references ("it")
			entityLabels.DecreaseActivation()

			// use dialog variables
			dialogizedParseTree := parse.NewDialogizer(base.dialogContext.VariableGenerator).Dialogize(parseTree, nil)
			base.log.AddProduction("Dialogized parse tree", dialogizedParseTree.IndentedString(""))

			// ellipsis
			ellipsizedParseTree, ok := parse.NewEllipsizer(clauseList.GetRootNodes(), base.log).Ellipsize(dialogizedParseTree)
			if !ok {
				break
			}
			base.log.AddProduction("Ellipsized parse tree", ellipsizedParseTree.IndentedString(""))

			// extract root clauses
			rootClauses := parse.NewRootClauseExtracter().Extract(ellipsizedParseTree)

			continueLooking := false
			for _, rootClauseTree := range rootClauses {
				tmpOutput, continueLooking, resolvedRemark, score = base.processRootClause(i18n, messenger, clauseList, deicticCenter, entityBindings, entityTags, entitySorts, entityLabels, entityDefinitions, grammar, rootClauseTree, locale, rawInput)

				if continueLooking {
					break
				}
			}

			if !continueLooking {
				// accept this tree
				output = tmpOutput
				break
			} else {
				// there's a problem, keep looking, but accept the message of the first problem (not the last one)
				if score > highestScore {
					output = tmpOutput
					highestScore = score
				}
			}
		}
	}

	if resolvedRemark != "" {
		output = resolvedRemark + "\n" + output
	}

	return base.waitForPrint(messenger, output)
}

func (base *LanguageBase) processRootClause(
	i18n *parse.I18n,
	messenger api.ProcessMessenger,
	clauseList *mentalese.ClauseList,
	deicticCenter *mentalese.DeicticCenter,
	entityBindings *mentalese.EntityBindings,
	entityTags *mentalese.TagList,
	entitySorts *mentalese.EntitySorts,
	entityLabels *mentalese.EntityLabels,
	entityDefinitions *mentalese.EntityDefinitions,
	grammar parse.Grammar,
	rootClauseTree *mentalese.ParseTreeNode,
	locale string,
	rawInput string,
) (string, bool, string, int) {

	score := 1

	// syntactic functions
	syntacticFunctions := mentalese.ExtractSyntacticFunctions(rootClauseTree)

	// request clause
	clause := mentalese.NewClause(rootClauseTree, false, syntacticFunctions)
	clauseList.AddClause(clause)

	// find tags
	tags := base.relationizer.ExtractTags(*rootClauseTree)
	entityTags.AddTags(tags)

	// write tags to dialog context db
	base.addTagsToDatabase(tags, messenger)

	// sorts
	sorts, sortFound := central.NewSortFinder(base.meta, messenger).FindSorts(rootClauseTree)
	if !sortFound {
		base.log.AddProduction("Break", "Breaking due to conflicting sorts: "+sorts.String())
		return "", true, "", score
	}
	score++
	for variable, sort := range sorts {
		entitySorts.SetSort(variable, sort)
	}

	// name resolution
	requestBinding, unresolvedName := base.resolveNames(messenger, rootClauseTree, entityBindings, entityTags, entitySorts)
	if unresolvedName != "" {
		return i18n.TranslateWithParam(common.NameNotFound, unresolvedName), true, "", score
	}
	score++

	// relationize
	requestRelations := base.relationizer.Relationize(rootClauseTree, []string{"S"})
	base.log.AddProduction("Relations", requestRelations.IndentedString(""))

	// entity definitions
	extracter := central.NewEntityDefinitionsExtracter(entityDefinitions)
	extracter.Extract(requestRelations)

	// anaphora
	resolver := central.NewAnaphoraResolver(base.log, clauseList, entityBindings, entityTags, entitySorts, entityLabels, entityDefinitions, base.meta, messenger)
	resolvedTree, resolvedRequest, resolvedBindings, resolvedOutput, resolvedRemark := resolver.Resolve(rootClauseTree, requestRelations, requestBinding)
	if resolvedOutput != "" {
		return resolvedOutput, true, "", score
	}
	score++

	// agreement
	agreementChecker := central.NewAgreementChecker()
	_, agreementOutput := agreementChecker.CheckAgreement(resolvedTree, entityTags)
	if agreementOutput != "" {
		return agreementOutput, true, resolvedRemark, score
	}
	score++

	// find intents
	intentRelations := base.relationizer.ExtractIntents(resolvedTree)
	intent, intentFound := base.answerer.FindIntent(append(resolvedRequest, intentRelations...))
	if !intentFound {
		return "No intent found", true, resolvedRemark, score
	}
	score++

	// execute intent
	executionBindings := messenger.ExecuteChildStackFrame(resolvedRequest, resolvedBindings)

	// center
	base.updateCenter(clauseList, deicticCenter)

	// response
	responseBindings, responseIndex, responseFound := base.findResponse(messenger, intent, executionBindings)
	if !responseFound {
		return "", false, resolvedRemark, score
	}
	score++

	// answer
	answerRelations, essentialBindings := base.createAnswer(messenger, intent, responseBindings, responseIndex)

	base.dialogWriteBindings(responseBindings, entityBindings, entitySorts)
	base.dialogWriteBindings(essentialBindings, entityBindings, entitySorts)

	// response clause
	base.dialogAddResponseClause(clauseList, deicticCenter, essentialBindings)

	// tokens
	tokens := base.generator.Generate(grammar.GetWriteRules(), answerRelations)

	// output
	surfacer := generate.NewSurfaceRepresentation(base.log)
	surface := surfacer.Create(tokens)

	return surface, false, resolvedRemark, score
}

func (base *LanguageBase) addTagsToDatabase(tags mentalese.RelationSet, messenger api.ProcessMessenger) {
	atomized := tags.ConvertVariablesToConstants()
	asserts := mentalese.RelationSet{}
	for _, relation := range atomized {
		asserts = append(asserts, mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
			mentalese.NewTermRelationSet(mentalese.RelationSet{relation}),
		}))
	}
	// fmt.Println(asserts.String())
	messenger.ExecuteChildStackFrame(asserts, mentalese.InitBindingSet(mentalese.NewBinding()))
}

func (base *LanguageBase) waitForPrint(messenger api.ProcessMessenger, output string) mentalese.BindingSet {
	uuid := common.CreateUuid()
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateWaitFor, []mentalese.Term{
			mentalese.NewTermRelationSet(
				mentalese.RelationSet{
					mentalese.NewRelation(false, mentalese.PredicatePrint, []mentalese.Term{
						mentalese.NewTermId(uuid, "entity"),
						mentalese.NewTermString(output),
					}),
				}),
		}),
	}
	bindings := messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(mentalese.NewBinding()))
	return bindings
}

func (base *LanguageBase) dialogAddResponseClause(clauseList *mentalese.ClauseList, deicticCenter *mentalese.DeicticCenter, essentialResponseBindings mentalese.BindingSet) {

	entities := []*mentalese.ClauseEntity{}
	for _, binding := range essentialResponseBindings.GetAll() {
		for _, variable := range binding.GetKeys() {
			entities = append(entities, mentalese.NewClauseEntity(variable, mentalese.AtomFunctionObject))
		}
	}

	clause := mentalese.NewClause(nil, true, entities)

	if len(entities) > 0 {
		deicticCenter.SetCenter(entities[0].DiscourseVariable)
	}

	clauseList.AddClause(clause)

	for _, binding := range essentialResponseBindings.GetAll() {
		for _, variable := range binding.GetKeys() {
			clause.AddEntity(variable)
		}
	}
}

func (base *LanguageBase) dialogWriteBindings(someBindings mentalese.BindingSet, entityBindings *mentalese.EntityBindings, entitySorts *mentalese.EntitySorts) {

	groupedValues := map[string][]mentalese.Term{}
	groupedSorts := map[string]string{}

	for _, someBinding := range someBindings.GetAll() {
		for key, value := range someBinding.GetAll() {
			if value.IsId() {

				_, found := groupedValues[key]
				if !found {
					groupedValues[key] = []mentalese.Term{}
					groupedSorts[key] = "entity"
				}

				alreadyAdded := false
				for _, v := range groupedValues[key] {
					if v.Equals(value) {
						alreadyAdded = true
					}
				}

				if !alreadyAdded {
					groupedValues[key] = append(groupedValues[key], value)
					groupedSorts[key] = value.TermSort
				}

			}
		}
	}

	for key, values := range groupedValues {
		if len(values) == 1 {
			entityBindings.Set(key, values[0])
		} else {
			entityBindings.Set(key, mentalese.NewTermList(values))
		}
		entitySorts.SetSort(key, groupedSorts[key])
	}
}

func (base *LanguageBase) findResponse(messenger api.ProcessMessenger, intent mentalese.Intent, resultBindings mentalese.BindingSet) (mentalese.BindingSet, int, bool) {

	for index := 0; index < len(intent.Responses); index++ {
		response := intent.Responses[index]
		if response.Condition.IsEmpty() {
			return resultBindings, index, true
		} else {
			responseBindings := messenger.ExecuteChildStackFrame(response.Condition, resultBindings)
			if !responseBindings.IsEmpty() {
				return responseBindings, index, true
			}
		}
	}

	return mentalese.NewBindingSet(), 0, false
}

func (base *LanguageBase) createAnswer(
	messenger api.ProcessMessenger,
	intent mentalese.Intent,
	resultBindings mentalese.BindingSet,
	responseIndex int,
) (mentalese.RelationSet, mentalese.BindingSet) {

	intentBindings := resultBindings
	resultHandler := intent.Responses[responseIndex]

	intentBindings = messenger.ExecuteChildStackFrame(resultHandler.Preparation, resultBindings)

	// create answer relation sets by binding 'answer' to solutionBindings
	answer := base.answerer.Build(resultHandler.Answer, intentBindings)

	base.log.AddProduction("Answer", answer.String())

	essential := mentalese.NewBindingSet()
	for _, id := range answer.GetIds() {
		newVariable := base.dialogContext.VariableGenerator.GenerateVariable("ResponseEntity")
		b := mentalese.NewBinding()
		b.Set(newVariable.TermValue, id)
		essential.Add(b)
	}
	return answer, essential
}

func (base *LanguageBase) updateCenter(clauseList *mentalese.ClauseList, deicticCenter *mentalese.DeicticCenter) {
	var previousCenter = deicticCenter.GetCenter()
	var center = ""
	var priority = 0

	priorities := map[string]int{
		"previousCenter":              100,
		mentalese.AtomFunctionSubject: 10,
		mentalese.AtomFunctionObject:  5,
	}

	clause := clauseList.GetLastClause()

	// new clause has no entities? keep existing center
	if len(clause.SyntacticFunctions) == 0 {
		center = previousCenter
	}

	for _, entity := range clause.SyntacticFunctions {
		if previousCenter != "" {
			priority = priorities["previousCenter"]
			center = entity.DiscourseVariable
			continue
		}
		prio, found := priorities[entity.SyntacticFunction]
		if found {
			if prio > priority {
				priority = prio
				center = entity.DiscourseVariable
			}
		}
	}

	deicticCenter.SetCenter(center)
}

func (base *LanguageBase) resolveNames(messenger api.ProcessMessenger, rootClauseTree *mentalese.ParseTreeNode, entityBindings *mentalese.EntityBindings, entityTags *mentalese.TagList, entitySorts *mentalese.EntitySorts) (mentalese.Binding, string) {

	names := base.nameResolver.ExtractNames(*rootClauseTree, []string{"S"})

	entityIds, nameNotFound, genderTags, numberTags := base.findNames(messenger, names, *entitySorts)
	entityTags.AddTags(genderTags)
	entityTags.AddTags(numberTags)

	requestBinding := entityBindings.Merge(entityIds)

	base.log.AddProduction("Named entities", entityIds.String())

	return requestBinding, nameNotFound
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

func (base *LanguageBase) findNames(messenger api.ProcessMessenger, names mentalese.Binding, sorts mentalese.EntitySorts) (mentalese.Binding, string, mentalese.RelationSet, mentalese.RelationSet) {

	entityIds := mentalese.NewBinding()
	nameNotFound := ""
	genderTags := mentalese.RelationSet{}
	numberTags := mentalese.RelationSet{}

	// look up entity ids by name
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
		nameInformations := base.nameResolver.ResolveName(name.TermValue, sort, messenger)
		if len(nameInformations) == 0 {
			base.log.AddProduction("Info",
				"Database lookup for name '"+name.TermValue+"'  with sort '"+sort+"' did not give any results")
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
				genderTags = append(genderTags, mentalese.NewRelation(false, mentalese.TagCategory, []mentalese.Term{
					mentalese.NewTermVariable(variable),
					mentalese.NewTermAtom(mentalese.AtomGender),
					mentalese.NewTermAtom(nameInformation.Gender),
				}))
			}
			if nameInformation.Number != "" {
				numberTags = append(numberTags, mentalese.NewRelation(false, mentalese.TagCategory, []mentalese.Term{
					mentalese.NewTermVariable(variable),
					mentalese.NewTermAtom(mentalese.AtomNumber),
					mentalese.NewTermAtom(nameInformation.Number),
				}))
			}
		}
	}

next:

	return entityIds, nameNotFound, genderTags, numberTags
}

//
// user functions
//

func (base *LanguageBase) dialogGetCenter(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	if !Validate(input, "v", base.log) {
		return mentalese.NewBindingSet()
	}

	centerVar := input.Arguments[0].TermValue

	center := mentalese.NewTermAtom("none")
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

func (base *LanguageBase) translate(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "sv", base.log) {
		return mentalese.NewBindingSet()
	}

	source := bound.Arguments[0].TermValue
	translatedVar := input.Arguments[1].TermValue

	translation := ""

	locale, localeFound := messenger.GetProcessSlot("locale")
	if !localeFound {
		translation = source
	} else {
		grammar, found := base.getGrammar(locale.TermValue)
		if !found {
			return mentalese.NewBindingSet()
		}

		i18n := parse.NewI18n(&grammar)
		translation = i18n.Translate(source)
	}

	newBinding := mentalese.NewBinding()
	newBinding.Set(translatedVar, mentalese.NewTermString(translation))
	return mentalese.InitBindingSet(newBinding)
}

func (base *LanguageBase) define(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {

	bound := input.BindSingle(binding)

	if !Validate(bound, "svr", base.log) {
		return mentalese.NewBindingSet()
	}

	word := bound.Arguments[0].TermValue
	entityVar := input.Arguments[1].TermValue
	sort := bound.Arguments[2].TermValueRelationSet

	locale, _ := messenger.GetProcessSlot("locale")
	grammar, found := base.getGrammar(locale.TermValue)
	if !found {
		return mentalese.NewBindingSet()
	}

	definition := sort.Copy()
	// definition = append(definition, specification...)
	// definition := specification.Copy()
	grammar.AddDefinition(word, entityVar, definition)

	return mentalese.InitBindingSet(binding)

}

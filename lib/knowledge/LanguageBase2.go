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

func (base *LanguageBase) reply(messenger api.ProcessMessenger, input mentalese.Relation, binding mentalese.Binding) mentalese.BindingSet {
	bound := input.BindSingle(binding)

	if !Validate(bound, "s", base.log) {
		return mentalese.NewBindingSet()
	}

	rawInput := bound.Arguments[0].TermValue

	output := ""

	// go:ignore(

	//     go:dialog_decrease_activation()
	base.dialogContext.EntityLabels.DecreaseActivation()

	//     go:find_locale(Locale)
	locale := ""
	for _, grammar := range base.grammars {
		locale = grammar.GetLocale()

		//     go:slot(locale, Locale)
		//     go:dialog_read_bindings(DialogBinding)
		dialogBinding := base.dialogContext.EntityBindings.Copy()
		//     go:tokenize(Locale, Input, InTokens)
		tokens := base.tokenize1(locale, rawInput)
		//     go:parse(Locale, InTokens, ParseTree)

		parser := parse.NewParser(grammar.GetReadRules(), base.log)
		parser.SetMorphologicalAnalyzer(grammar.GetMorphologicalAnalyzer())
		parseTrees := parser.Parse(tokens, "s", []string{"S"})

		for _, parseTree := range parseTrees {
			base.log.AddProduction("Parse tree", parseTree.IndentedString(""))

			//     /* Stop at the first successfully processed parse tree */
			//     go:cut(1,
			// go:dialogize(ParseTree, DialogParseTree)
			dialogizer := parse.NewDialogizer(base.dialogContext.VariableGenerator)
			dialogizedParseTree := dialogizer.Dialogize(&parseTree)

			base.log.AddProduction("Dialogized parse tree", dialogizedParseTree.IndentedString(""))

			//         go:ellipsize(DialogParseTree, CompletedParseTree)
			clauses := base.dialogContext.ClauseList.GetRootNodes()
			ellipsizer := parse.NewEllipsizer(clauses, base.log)
			ellipsizedParseTree, ok := ellipsizer.Ellipsize(*dialogizedParseTree)
			if !ok {
				break
			}

			base.log.AddProduction("Ellipsized parse tree", ellipsizedParseTree.IndentedString(""))

			//         go:extract_root_clauses(CompletedParseTree, RootClauseTree)
			rootClauseExtracter := parse.NewRootClauseExtracter()
			rootClauses := rootClauseExtracter.Extract(&parseTree)

			continueLooking := false
			for _, rootClauseTree := range rootClauses {

				//         go:process_root_clause(RootClauseTree, DialogBinding, Locale, RootClauseOutput, ContinueLooking)
				output, continueLooking = base.processRootClause(messenger, grammar, rootClauseTree, dialogBinding, locale, rawInput)
				//         [:Output := RootClauseOutput]
				// output = rootClauseOutput

				if continueLooking {
					break
				}

			}

			//         /* `ContinueLooking` tells us that there was a problem that should only be reported if no other parse tree succeeds */
			//         if [ContinueLooking == true] then
			//             fail
			//         end
			if !continueLooking {
				// accept this one
				break
			}
		}

	}

	uuid := common.CreateUuid()

	// )
	// go:uuid(Uuid, entity)
	// go:wait_for(
	//     go:print(Uuid, :Output)
	// )
	set := mentalese.RelationSet{
		mentalese.NewRelation(false, mentalese.PredicateWaitFor, []mentalese.Term{
			mentalese.NewTermRelationSet(
				mentalese.RelationSet{
					mentalese.NewRelation(false, mentalese.PredicatePrint, []mentalese.Term{
						mentalese.NewTermVariable(uuid),
						mentalese.NewTermVariable(output),
					}),
				}),
		}),
	}

	bindings := messenger.ExecuteChildStackFrame(set, mentalese.InitBindingSet(mentalese.NewBinding()))
	return bindings
}

func (base *LanguageBase) processRootClause(messenger api.ProcessMessenger, grammar parse.Grammar, rootClauseTree *mentalese.ParseTreeNode, dialogBinding mentalese.Binding, locale string, rawInput string) (string, bool) {
	rootClauseOutput := ""
	continueLooking := false

	// go:dialog_add_root_clause(RootClauseTree, false, ClauseVariable)
	clauseList := base.dialogContext.ClauseList
	entities := mentalese.ExtractEntities(rootClauseTree)
	clause := mentalese.NewClause(rootClauseTree, false, entities)
	clauseList.AddClause(clause)

	// rootVariable := rootClauseTree.Rule.EntityVariables[0][0]

	// go:extract_tags_and_intents(RootClauseTree)
	tags := base.relationizer.ExtractTags(*rootClauseTree)
	base.dialogContext.EntityTags.AddTags(tags)

	intentRelations := base.relationizer.ExtractIntents(*rootClauseTree)
	base.dialogContext.ClauseList.GetLastClause().SetIntents(intentRelations)

	// go:sortal_filtering(RootClauseTree)
	// extract sorts: variable => sort
	sortFinder := central.NewSortFinder(base.meta, messenger)
	sorts, sortFound := sortFinder.FindSorts(rootClauseTree)
	if !sortFound {
		// conflicting sorts
		base.log.AddProduction("Break", "Breaking due to conflicting sorts: "+sorts.String())
		return "", true
	}

	for variable, sort := range sorts {
		base.dialogContext.EntitySorts.SetSorts(variable, []string{sort})
	}

	// go:resolve_names(RootClauseTree, DialogBinding, RequestBinding, UnresolvedName)
	requestBinding, unresolvedName := base.resolveNames1(messenger, rootClauseTree, dialogBinding)
	// if [UnresolvedName != ''] then
	//     go:create_canned(Output, name_not_found, UnresolvedName)
	//     [ContinueLooking := true]
	//     return
	// end
	if unresolvedName != "" {
		rootClauseOutput = common.GetString("name_not_found", unresolvedName)
		return rootClauseOutput, true
	}

	// go:relationize(RootClauseTree, Request)
	requestRelations := base.relationizer.Relationize(*rootClauseTree, []string{"S"})

	base.log.AddProduction("Relations", requestRelations.IndentedString(""))

	extracter := central.NewEntityDefinitionsExtracter(base.dialogContext)
	extracter.Extract(requestRelations)

	// go:resolve_anaphora(RootClauseTree, Request, RequestBinding, ResolvedRequest, ResolvedBindings, ResolvedOutput)
	resolver := central.NewAnaphoraResolver(base.dialogContext, base.meta, messenger)
	resolvedRequest, resolvedBindings, resolvedOutput := resolver.Resolve(rootClauseTree, requestRelations, requestBinding)
	// if [ResolvedOutput != ''] then
	//     [Output := ResolvedOutput]
	//     [ContinueLooking := false]
	//     return
	// end
	if resolvedOutput != "" {
		return resolvedOutput, false
	}

	// go:check_agreement(RootClauseTree, AgreementOutput)
	agreementChecker := central.NewAgreementChecker()
	_, agreementOutput := agreementChecker.CheckAgreement(rootClauseTree, base.dialogContext.EntityTags)
	// if [AgreementOutput != ''] then
	//     [Output := AgreementOutput]
	//     [ContinueLooking := true]
	//     return
	// end
	if agreementOutput != "" {
		return agreementOutput, true
	}

	// /* Stop at the first successfully processed intent */
	// go:cut(1,
	//     go:detect_intent(ResolvedRequest, Intent, ClauseVariable)
	//     go:make_list(Intents, Intent)
	//     go:execute_intent(ResolvedRequest, ResolvedBindings, Intents, AnOutput, Accepted, AcceptedBindings)
	//     if [AnOutput != ''] then
	//          [Output := AnOutput]
	//          [ContinueLooking := false]
	//          return
	//     end
	// )

	// todo same variable as above?
	intentRelations2 := base.dialogContext.ClauseList.GetLastClause().GetIntents()

	conditionSubject := append(requestRelations, intentRelations2...)
	intents := base.answerer.FindIntents(conditionSubject)

	anOutput, acceptedIntent, acceptedBindings := base.executeIntent(resolvedRequest, resolvedBindings, intents)
	if anOutput != "" {
		return anOutput, false
	}

	// go:dialog_update_center()
	base.updateCenter()

	// go:find_response(Accepted, AcceptedBindings, ResponseBindings, ResponseIndex)
	responseBindings, responseIndex, responseFound := base.findResponse1(messenger, acceptedIntent, acceptedBindings)
	if !responseFound {
		return "", false
	}

	// go:create_answer(Accepted, ResponseBindings, ResponseIndex, Answer, EssentialBindings)
	answerRelations, essentialBindings := base.createAnswer1(messenger, acceptedIntent, responseBindings, responseIndex)

	// go:dialog_write_bindings(ResponseBindings)
	base.dialogWriteBindings1(resolvedBindings)

	// go:dialog_write_bindings(EssentialBindings)
	base.dialogWriteBindings1(essentialBindings)

	// go:dialog_add_response_clause(EssentialBindings)
	base.dialogAddResponseClause1(essentialBindings)

	// go:generate(Locale, Answer, OutTokens)
	tokens := base.generator.Generate(grammar.GetWriteRules(), answerRelations)

	// go:surface(OutTokens, Output)
	surfacer := generate.NewSurfaceRepresentation(base.log)
	surface := surfacer.Create(tokens)

	// [ContinueLooking := false]
	return surface, false

	// return rootClauseOutput, continueLooking
}

func (base *LanguageBase) dialogAddResponseClause1(essentialResponseBindings mentalese.BindingSet) {
	// bound := input.BindSingle(binding)

	// essentialResponseBindings := mentalese.NewBindingSet()
	// someBindingsRaw := []map[string]mentalese.Term{}
	// someBindingsRaw = bound.Arguments[0].GetBinaryValue().([]map[string]mentalese.Term)
	// essentialResponseBindings.FromRaw(someBindingsRaw)

	entities := []*mentalese.ClauseEntity{}
	for _, binding := range essentialResponseBindings.GetAll() {
		for _, variable := range binding.GetKeys() {
			entities = append(entities, mentalese.NewClauseEntity(variable, mentalese.AtomFunctionObject))
		}
	}

	clause := mentalese.NewClause(nil, true, entities)

	if len(entities) > 0 {
		base.dialogContext.DeicticCenter.SetCenter(entities[0].DiscourseVariable)
	}

	base.dialogContext.ClauseList.AddClause(clause)

	for _, binding := range essentialResponseBindings.GetAll() {
		for _, variable := range binding.GetKeys() {
			clause.AddEntity(variable)
		}
	}

	// return mentalese.InitBindingSet(binding)
}

func (base *LanguageBase) dialogWriteBindings1(someBindings mentalese.BindingSet) {
	// bound := input.BindSingle(binding)

	// someBindings := mentalese.NewBindingSet()
	// someBindingsRaw := []map[string]mentalese.Term{}
	// someBindingsRaw = bound.Arguments[0].GetBinaryValue().([]map[string]mentalese.Term)
	// someBindings.FromRaw(someBindingsRaw)

	groupedValues := map[string][]mentalese.Term{}
	groupedSorts := map[string][]string{}

	for _, someBinding := range someBindings.GetAll() {
		for key, value := range someBinding.GetAll() {
			if value.IsId() {

				_, found := groupedValues[key]
				if !found {
					groupedValues[key] = []mentalese.Term{}
					groupedSorts[key] = []string{}
				}

				alreadyAdded := false
				for _, v := range groupedValues[key] {
					if v.Equals(value) {
						alreadyAdded = true
					}
				}

				if !alreadyAdded {
					groupedValues[key] = append(groupedValues[key], value)
					groupedSorts[key] = append(groupedSorts[key], value.TermSort)
				}

			}
		}
	}

	for key, values := range groupedValues {
		if len(values) == 1 {
			base.dialogContext.EntityBindings.Set(key, values[0])
		} else {
			base.dialogContext.EntityBindings.Set(key, mentalese.NewTermList(values))
		}
		base.dialogContext.EntitySorts.SetSorts(key, groupedSorts[key])
	}

	// return mentalese.InitBindingSet(binding)
}

func (base *LanguageBase) findResponse1(messenger api.ProcessMessenger, intent mentalese.Intent, resultBindings mentalese.BindingSet) (mentalese.BindingSet, int, bool) {

	// // bound := input.BindSingle(binding)

	// // if !Validate(bound, "jjvv", base.log) {
	// // 	return mentalese.NewBindingSet()
	// // }

	// // intent := mentalese.Intent{}
	// resultBindings := mentalese.NewBindingSet()

	// intent = bound.Arguments[0].GetBinaryValue().(mentalese.Intent)

	// resultBindingsRaw := []map[string]mentalese.Term{}
	// resultBindingsRaw = bound.Arguments[1].GetBinaryValue().([]map[string]mentalese.Term)
	// resultBindings.FromRaw(resultBindingsRaw)

	// responseBindingsVar := input.Arguments[2].TermValue
	// responseIndexVar := input.Arguments[3].TermValue

	for index := 0; index < len(intent.Responses); index++ {
		response := intent.Responses[index]
		if response.Condition.IsEmpty() {
			// newBinding := mentalese.NewBinding()
			// newBinding.Set(responseBindingsVar, mentalese.NewTermBinary(resultBindings.ToRaw()))
			// newBinding.Set(responseIndexVar, mentalese.NewTermString(strconv.Itoa(index)))
			// return mentalese.InitBindingSet(newBinding)
			return resultBindings, index, true
		} else {
			responseBindings := messenger.ExecuteChildStackFrame(response.Condition, resultBindings)
			if !responseBindings.IsEmpty() {
				// newBinding := mentalese.NewBinding()
				// newBinding.Set(responseBindingsVar, mentalese.NewTermBinary(responseBindings.ToRaw()))
				// newBinding.Set(responseIndexVar, mentalese.NewTermString(strconv.Itoa(index)))
				// return mentalese.InitBindingSet(newBinding)
				return resultBindings, index, true
			}
		}
	}

	// return mentalese.NewBindingSet()
	return mentalese.NewBindingSet(), 0, false
}

func (base *LanguageBase) createAnswer1(messenger api.ProcessMessenger, intent mentalese.Intent, resultBindings mentalese.BindingSet, responseIndex int) (mentalese.RelationSet, mentalese.BindingSet) {

	// bound := input.BindSingle(binding)

	// if !Validate(bound, "jjivv", base.log) {
	// 	return mentalese.NewBindingSet()
	// }

	// intent := mentalese.Intent{}
	// resultBindings := mentalese.NewBindingSet()

	// intent = bound.Arguments[0].GetBinaryValue().(mentalese.Intent)

	// responseBindingsRaw := []map[string]mentalese.Term{}
	// responseBindingsRaw = bound.Arguments[1].GetBinaryValue().([]map[string]mentalese.Term)
	// resultBindings.FromRaw(responseBindingsRaw)

	// responseIndex, _ := bound.Arguments[2].GetIntValue()
	// answerVar := input.Arguments[3].TermValue
	// essentialVar := input.Arguments[4].TermValue

	intentBindings := resultBindings
	resultHandler := intent.Responses[responseIndex]

	intentBindings = messenger.ExecuteChildStackFrame(resultHandler.Preparation, resultBindings)

	// create answer relation sets by binding 'answer' to solutionBindings
	answer := base.answerer.Build(resultHandler.Answer, intentBindings)

	base.log.AddProduction("Answer", answer.String())

	// newBinding := mentalese.NewBinding()
	// newBinding.Set(answerVar, mentalese.NewTermRelationSet(answer))

	essential := mentalese.NewBindingSet()
	for _, id := range answer.GetIds() {
		newVariable := base.dialogContext.VariableGenerator.GenerateVariable("ResponseEntity")
		b := mentalese.NewBinding()
		b.Set(newVariable.TermValue, id)
		essential.Add(b)
	}

	// newBinding.Set(essentialVar, mentalese.NewTermBinary(essential.ToRaw()))

	// return mentalese.InitBindingSet(newBinding)
	return answer, essential
}

func (base *LanguageBase) executeIntent(messenger api.ProcessMessenger, resolvedRequest mentalese.RelationSet, resolvedBindings mentalese.BindingSet, intents []mentalese.Intent) (string, mentalese.Intent, mentalese.BindingSet) {

	// go:list_length(Intents, SolSize)
	// [LastSol := [SolSize - 1]]
	// [:Output := '']
	// [:Accepted := none]
	// [:AcceptedBindings := none]

	output := ""
	accepted := mentalese.Intent{}
	acceptedBindings := mentalese.BindingSet{}

	for index, sol := range intents {
		resultBindings := messenger.ExecuteChildStackFrame(resolvedRequest, resolvedBindings)
		if resultBindings.GetLength() > 0 {
			accepted = sol
			acceptedBindings = resultBindings
			break
		}
		if index == len(intents)-1 {
			accepted = sol
			acceptedBindings = resultBindings
			break
		}
	}

	// go:list_foreach(Intents, Index, Sol,

	//     go:solve(ResolvedRequest, ResolvedBindings, Sol, ResultBindings, ResultCount)

	//     /* If there were Intents, good! */
	//     if [ResultCount > 0] then
	//         [:Accepted := Sol]
	//         [:AcceptedBindings := ResultBindings]
	//         [:Output := '']
	//         break
	//     end
	//     /* If this is the last available intent, take it anyway */
	//     if [LastSol == Index] then
	//         [:Accepted := Sol]
	//         [:AcceptedBindings := ResultBindings]
	//         [:Output := '']
	//         break
	//     end
	// )
	// [Output := :Output]
	// [Accepted := :Accepted]
	// [AcceptedBindings := :AcceptedBindings]

	return output, accepted, acceptedBindings
}

func (base *LanguageBase) updateCenter() {
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

			panic("this shouldn't work ?!")

			// a := getValue(entity.DiscourseVariable, binding)
			// b := getValue(previousCenter, binding)
			// if a == b {
			// 	priority = priorities["previousCenter"]
			// 	center = entity.DiscourseVariable
			// 	continue
			// }
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

	// return mentalese.InitBindingSet(binding)
}

func (base *LanguageBase) resolveNames1(messenger api.ProcessMessenger, rootClauseTree *mentalese.ParseTreeNode, dialogBinding mentalese.Binding) (mentalese.Binding, string) {
	// requestBindingVar := input.Arguments[2].TermValue
	// unboundNameVar := input.Arguments[3].TermValue
	// var parseTree *mentalese.ParseTreeNode

	// parseTree = bound.Arguments[0].GetBinaryValue().(*mentalese.ParseTreeNode)

	// dialogBinding := mentalese.NewBinding()
	// dialogBindingsRaw := map[string]mentalese.Term{}
	// dialogBindingsRaw = bound.Arguments[1].GetBinaryValue().(map[string]mentalese.Term)
	// dialogBinding.FromRaw(dialogBindingsRaw)

	names := base.nameResolver.ExtractNames(*rootClauseTree, []string{"S"})

	sorts := base.dialogContext.EntitySorts

	entityIds, nameNotFound, genderTags, numberTags := base.findNames(messenger, names, *sorts)
	base.dialogContext.EntityTags.AddTags(genderTags)
	base.dialogContext.EntityTags.AddTags(numberTags)

	requestBinding := dialogBinding.Merge(entityIds)

	base.log.AddProduction("Named entities", entityIds.String())

	// newBinding := binding.Copy()

	// newBinding.Set(requestBindingVar, mentalese.NewTermBinary(requestBinding.ToRaw()))
	// newBinding.Set(unboundNameVar, mentalese.NewTermString(nameNotFound))
	return requestBinding, nameNotFound
}

func (base *LanguageBase) tokenize1(locale string, rawInput string) []string {
	grammar, found := base.getGrammar(locale)
	if !found {
		return []string{}
	}

	tokens := grammar.GetTokenizer().Process(rawInput)

	base.log.AddProduction("Tokens", strings.Join(tokens, " "))

	// terms := []mentalese.Term{}
	// for _, token := range tokens {
	// 	terms = append(terms, mentalese.NewTermString(token))
	// }

	// newBinding := binding.Copy()
	// newBinding.Set(tokenVar, mentalese.NewTermList(terms))
	return tokens
}

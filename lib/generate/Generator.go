package generate

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
)

type Generator struct {
	matcher *central.RelationMatcher
	parser  *importer.InternalGrammarParser
	state   *mentalese.GenerationState
	log     *common.SystemLog
}

func NewGenerator(log *common.SystemLog, matcher *central.RelationMatcher, state *mentalese.GenerationState) *Generator {

	return &Generator{
		matcher: matcher,
		state:   state,
		parser:  importer.NewInternalGrammarParser(),
		log:     log,
	}
}

// Creates an array of words that forms the surface representation of a mentalese sense
func (generator *Generator) Generate(grammarRules *mentalese.GrammarRules, sentenceSense mentalese.RelationSet) []string {

	generator.state.Clear()

	// canned response
	if !sentenceSense.IsEmpty() && sentenceSense[0].Predicate == mentalese.PredicateCanned {
		return []string{sentenceSense[0].Arguments[0].TermValue}
	}

	// convert variables to constants
	boundSense := sentenceSense.ConvertVariablesToConstants()

	generator.log.AddDebug("Constants", fmt.Sprintf("%v", boundSense))

	boundSense = boundSense.ExpandChildren()

	generator.log.AddDebug("Unscoped 2", fmt.Sprintf("%v", boundSense))

	return generator.GenerateNode(grammarRules, []string{}, "s", mentalese.NewTermString(""), boundSense)
}

// Creates an array of words for a syntax tree node
// antecedent: i.e. np(E1)
// antecedentBinding i.e. { E1: 1 }
func (generator *Generator) GenerateNode(grammarRules *mentalese.GrammarRules, usedRules []string, antecedentCategory string, antecedentValue mentalese.Term, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	// condition matches: grammatical_subject(E), subject(P, E)
	// rule: s(P) :- np(E), vp(P)
	rule, binding, ok := generator.findMatchingRule(grammarRules, usedRules, antecedentCategory, antecedentValue, sentenceSense)

	if ok {

		hash := generator.createRuleHash(rule, binding)
		usedRules = append(usedRules, hash)

		if generator.log.Active() {
			generator.log.AddDebug("Found", fmt.Sprintf("%v %v ", rule.String(), binding.String()))
		}

		for i, consequentCategory := range rule.GetConsequents() {
			consequentValue := generator.getConsequentValue(rule, i, binding)
			consequent := generator.generateSingleConsequent(
				grammarRules, usedRules, consequentCategory, consequentValue, rule.GetConsequentPositionType(i), sentenceSense)
			words = append(words, consequent...)

			if consequentValue.IsId() && !consequentValue.Equals(antecedentValue) {
				generator.state.MarkGenerated(consequentValue)
			}

		}
	} else {

		// place breakpoint here for debugging ;)
		//generator.findMatchingRule(messenger, grammarRules, usedRules, antecedentCategory, antecedentValue, sentenceSense)

		generator.log.AddError("No rule found for " + fmt.Sprintf("%v(%v)", antecedentCategory, antecedentValue))

		//rule, binding, ok = generator.findMatchingRule(messenger, grammarRules, usedRules, antecedentCategory, antecedentValue, sentenceSense)
	}

	return words
}

func (generator *Generator) getConsequentValue(rule mentalese.GrammarRule, i int, binding mentalese.Binding) mentalese.Term {
	consequentValue := mentalese.NewTermString("")
	found := false

	if rule.GetConsequentPositionType(i) != mentalese.PosTypeWordForm {
		consequentVariable := rule.GetConsequentVariables(i)[0]
		consequentValue, found = binding.Get(consequentVariable)
		if !found {
			consequentValue = mentalese.NewTermString("")
		}
	}

	return consequentValue
}

// From a set of rules (with a shared antecedent), find the first one whose conditions match
// antecedent: i.e. np(E1)
// bindingL i.e { E1: 3 }
func (generator *Generator) findMatchingRule(grammarRules *mentalese.GrammarRules, usedRules []string, antecedentCategory string, antecedentValue mentalese.Term, sentenceSense mentalese.RelationSet) (mentalese.GrammarRule, mentalese.Binding, bool) {

	found := false
	resultRule := mentalese.GrammarRule{}
	binding := mentalese.NewBinding()

	rules := grammarRules.FindRules(antecedentCategory, 1)

	for _, rule := range rules {

		// copy the value of the antecedent variable
		ruleAntecedentVariable := rule.GetAntecedentVariables()[0]
		binding = mentalese.NewBinding()

		if !(antecedentValue.IsString() && antecedentValue.TermValue == "") {
			binding.Set(ruleAntecedentVariable, antecedentValue)
		}

		if len(rule.Sense) == 0 {

			// no condition
			resultRule = rule
			found = true

		} else {

			// match the condition
			matchBindings, match := generator.matcher.MatchSequenceToSet(rule.Sense, sentenceSense, binding)

			if match {
				binding = matchBindings.Get(0)
				resultRule = rule
				found = true
			}
		}

		// match the goal
		if found {
			binding, found = generator.matcher.MatchTerm(antecedentValue, generator.val2term(rule.GetAntecedentVariables()[0]), binding)
		}

		if found {

			// make sure the same rule is not executed again and again
			hash := generator.createRuleHash(resultRule, binding)
			if !common.StringArrayContains(usedRules, hash) {
				break
			}
		}
	}

	return resultRule, binding, found
}

// todo: rewrite rules should consist of relations
func (generator *Generator) val2term(val string) mentalese.Term {
	return generator.parser.CreateTerm(val)
}

func (generator *Generator) createRuleHash(rule mentalese.GrammarRule, binding mentalese.Binding) string {
	return rule.BindSimple(binding).String()
}

// From one of the bound consequents of a syntactic rewrite rule, generate an array of words
// vp(P1) => married Marry
func (generator *Generator) generateSingleConsequent(grammarRules *mentalese.GrammarRules, usedRules []string, category string, value mentalese.Term, positionType string, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	if positionType == mentalese.PosTypeWordForm {
		words = append(words, category)
	} else if category == mentalese.CategoryText {
		text := value
		words = append(words, text.TermValue)
	} else {
		words = generator.GenerateNode(grammarRules, usedRules, category, value, sentenceSense)
	}

	return words
}

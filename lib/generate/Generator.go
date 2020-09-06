package generate

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
)

type Generator struct {
	matcher *mentalese.RelationMatcher
	log     *common.SystemLog
}

func NewGenerator(log *common.SystemLog, matcher *mentalese.RelationMatcher) *Generator {
	return &Generator{
		matcher: matcher,
		log: log,
	}
}

// Creates an array of words that forms the surface representation of a mentalese sense
func (generator *Generator) Generate(grammarRules *parse.GrammarRules, sentenceSense mentalese.RelationSet) []string {

	// canned response
	if !sentenceSense.IsEmpty() && sentenceSense[0].Predicate == mentalese.PredicateCanned {
		return []string{ sentenceSense[0].Arguments[0].TermValue }
	}

	// convert variables to constants
	boundSense := sentenceSense.ConvertVariablesToConstants()

	generator.log.AddProduction("Constants", fmt.Sprintf("%v", boundSense))

	return generator.GenerateNode(grammarRules, "s", []string{"S1"}, mentalese.Binding{}, boundSense)
}

// Creates an array of words for a syntax tree node
// antecedent: i.e. np(E1)
// antecedentBinding i.e. { E1: 1 }
func (generator *Generator) GenerateNode(grammarRules *parse.GrammarRules, predicate string, arguments []string, antecedentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	generator.log.StartDebug("GenerateNode", predicate, arguments, antecedentBinding)

	// condition matches: grammatical_subject(E), subject(P, E)
	// rule: s(P) :- np(E), vp(P)
	rule, conditionBinding, ok := generator.findMatchingRule(grammarRules, predicate, arguments, antecedentBinding, sentenceSense)

	if ok {

		for i, consequentPredicate := range rule.GetConsequents() {

			consequentArguments := rule.GetConsequentVariables(i)
			words = append(words, generator.generateSingleConsequent(grammarRules, consequentPredicate, consequentArguments, rule.GetConsequentPositionType(i), conditionBinding, sentenceSense)...)
		}
	} else {
		generator.log.AddError("Cannot generate response for syntax node " + predicate)
	}

	generator.log.EndDebug("GenerateNode", words)

	return words
}

// From a set of rules (with a shared antecedent), find the first one whose conditions match
// antecedent: i.e. np(E1)
// bindingL i.e { E1: 3 }
func (generator *Generator) findMatchingRule(grammarRules *parse.GrammarRules, predicate string, arguments []string, antecedentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) (parse.GrammarRule, mentalese.Binding, bool) {

	found := false
	resultRule := parse.GrammarRule{}
	conditionBinding := mentalese.Binding{}

	generator.log.StartDebug("findMatchingRule", predicate, antecedentBinding)

	rules := grammarRules.FindRules(predicate, len(arguments))

	for _, rule := range rules {

		// copy the value of the predicate
		conditionBinding = mentalese.Binding{}
		for i, argument := range arguments {
			ruleAntecedentVariable := rule.GetAntecedentVariables()[i]
			val, found := antecedentBinding[argument]
			if found {
				conditionBinding[ruleAntecedentVariable] = val
			}
		}

		if len(rule.Sense) == 0 {

			// no condition
			resultRule = rule
			found = true
			break

		} else {

			// match the condition
			matchBindings, match := generator.matcher.MatchSequenceToSet(rule.Sense, sentenceSense, conditionBinding)

			if match {
				conditionBinding = matchBindings[0]
				resultRule = rule
				found = true
				break
			}
		}
	}

	generator.log.EndDebug("findMatchingRule", resultRule, conditionBinding, found)

	return resultRule, conditionBinding, found
}

// From one of the bound consequents of a syntactic rewrite rule, generate an array of words
// vp(P1) => married Marry
func (generator *Generator) generateSingleConsequent(grammarRules *parse.GrammarRules, predicate string, arguments []string, positionType string, consequentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	generator.log.StartDebug("generateSingleConsequent", predicate, consequentBinding)

	if positionType == parse.PosTypeWordForm {
		words = append(words, predicate)
	} else if predicate == mentalese.CategoryText {
		variable := arguments[0]
		text := consequentBinding[variable]
		words = append(words, text.TermValue)
	} else {
		words = generator.GenerateNode(grammarRules, predicate, arguments, consequentBinding, sentenceSense)
	}

	generator.log.EndDebug("generateSingleConsequent", words)

	return words
}

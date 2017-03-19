package generate

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type Generator struct {
	Grammar *GenerationGrammar
	Lexicon *GenerationLexicon
	matcher mentalese.RelationMatcher
}

func NewGenerator(Grammar *GenerationGrammar, Lexicon *GenerationLexicon) *Generator {
	return &Generator{Grammar:Grammar, Lexicon:Lexicon, matcher:mentalese.RelationMatcher{}}
}

// Creates an array of words that forms the surface representation of a mentalese sense
func (generator *Generator) Generate(sentenceSense mentalese.RelationSet) []string {

	rootAntecedent := mentalese.Relation{Predicate:"s", Arguments:[]mentalese.Term{{mentalese.Term_variable, "S1", mentalese.RelationSet{}}}}

	return generator.GenerateNode(rootAntecedent, mentalese.Binding{}, sentenceSense)
}

// Creates an array of words for a syntax tree node
// antecedent: i.e. np(E1)
// antecedentBinding i.e. { E1: 1 }
func (generator *Generator) GenerateNode(antecedent mentalese.Relation, antecedentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	common.LogTree("GenerateNode", antecedent, antecedentBinding)

	// condition matches: grammatical_subject(E), subject(P, E)
	// rule: s(P) :- np(E), vp(P)
	rule, conditionBinding, ok := generator.findMatchingRule(antecedent, antecedentBinding, sentenceSense)

	if ok {

		for _, consequent := range rule.Consequents {

			consequentBinding := conditionBinding.Extract(consequent.Arguments[0].TermValue)
			words = append(words, generator.generateSingleConsequent(consequent, consequentBinding, sentenceSense)...)
		}
	}

	common.LogTree("GenerateNode", words)

	return words
}

// From a set of rules (with a shared antecedent), find the first one whose conditions match
// antecedent: i.e. np(E1)
// bindingL i.e { E1: 3 }
func (generator *Generator) findMatchingRule(antecedent mentalese.Relation, antecedentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) (GenerationGrammarRule, mentalese.Binding, bool) {

	found := false
	resultRule := GenerationGrammarRule{}
	conditionBinding := mentalese.Binding{}

	common.LogTree("findMatchingRule", antecedent, antecedentBinding)

	rules := generator.Grammar.FindRules(antecedent)

	for _, rule := range rules {

		// copy the value of the antecedent
		conditionBinding = mentalese.Binding{}
		if len(antecedentBinding) > 0 {
			parentAntecedentVariable := antecedent.Arguments[0].TermValue
			ruleAntecedentVariable := rule.Antecedent.Arguments[0].TermValue
			conditionBinding[ruleAntecedentVariable] = antecedentBinding[parentAntecedentVariable]
		}

		if len(rule.Condition) == 0 {

			// no condition
			resultRule = rule
			found = true
			break

		} else {

			// match the condition
			matchBindings, _, match := generator.matcher.MatchSequenceToSet(rule.Condition, sentenceSense, conditionBinding)

			if match {
				conditionBinding = matchBindings[0]
				resultRule = rule
				found = true
				break
			}
		}
	}

	common.LogTree("findMatchingRule", resultRule, conditionBinding, found)

	return resultRule, conditionBinding, found
}

// From one of the bound consequents of a syntactic rewrite rule, generate an array of words
// vp(P1) => married Marry
func (generator *Generator) generateSingleConsequent(consequent mentalese.Relation, consequentBinding mentalese.Binding, sentenceSense mentalese.RelationSet) []string {

	words := []string{}
	found := false

	common.LogTree("generateSingleConsequent", consequent, consequentBinding)

	boundConsequent := generator.matcher.BindSingleRelationSingleBinding(consequent, consequentBinding)

	lexItem, found := generator.Lexicon.GetLexemeForGeneration(boundConsequent, sentenceSense)
	if found {
		words = append(words, lexItem.Form)
	} else {
		words = generator.GenerateNode(consequent, consequentBinding, sentenceSense)
	}

	common.LogTree("generateSingleConsequent", words)

	return words
}
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

	rootAntecedent := mentalese.Relation{Predicate:"s", Arguments:[]mentalese.Term{{mentalese.Term_variable, "S1"}}}

	return generator.GenerateNode(rootAntecedent, sentenceSense)
}

// Creates an array of words for a syntax tree node
// antecedent: i.e. np(E1)
func (generator *Generator) GenerateNode(antecedent mentalese.Relation, sentenceSense mentalese.RelationSet) []string {

	words := []string{}

	common.LogTree("GenerateNode", antecedent)

	// condition matches: grammatical_subject(E), subject(P, E)
	// rule: s(P) :- np(E), vp(P)
	rule, binding, ok := generator.findMatchingRule(antecedent, sentenceSense)

	if ok {

		boundConsequents := generator.matcher.BindRelationSetSingleBinding(rule.Consequents, binding)

		for _, consequent:= range boundConsequents {
			words = append(words, generator.generateSingleConsequent(consequent, sentenceSense)...)
		}
	}

	common.LogTree("GenerateNode", words)

	return words
}

// From a set of rules (with a shared antecedent), find the first one whose conditions match
// antecedent: i.e. np(E1)
func (generator *Generator) findMatchingRule(antecedent mentalese.Relation, sentenceSense mentalese.RelationSet) (GenerationGrammarRule, mentalese.Binding, bool) {

	found := false
	resultRule := GenerationGrammarRule{}
	binding := mentalese.Binding{}

	common.LogTree("findMatchingRule", antecedent)

	rules := generator.Grammar.FindRules(antecedent)

	for _, rule := range rules {

		if len(rule.Condition) == 0 {

			// no condition

// note: this binding should probably be done in the else case as well (?)

			binding, _ = generator.matcher.MatchTwoRelations(rule.Antecedent, antecedent, mentalese.Binding{})
			resultRule = rule
			found = true
			break

		} else {

			bindings, _, match := generator.matcher.MatchSequenceToSet(rule.Condition, sentenceSense, mentalese.Binding{})

			if match {
				resultRule = rule
				binding = bindings[0]
				found = true
				break
			}
		}

	}

	common.LogTree("findMatchingRule", resultRule, binding, found)

	return resultRule, binding, found
}

// From one of the bound consequents of a syntactic rewrite rule, generate an array of words
// vp(P1) => married Marry
func (generator *Generator) generateSingleConsequent(consequent mentalese.Relation, sentenceSense mentalese.RelationSet) []string {

	words := []string{}
	found := false

	common.LogTree("generateSingleConsequent", consequent)

	lexItem, found := generator.Lexicon.GetLexemeForGeneration(consequent, sentenceSense)
	if found {
		words = append(words, lexItem.Form)
	} else {
		words = generator.GenerateNode(consequent, sentenceSense)
	}

	common.LogTree("generateSingleConsequent", words)

	return words
}
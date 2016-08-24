package tests

import "testing"
import (
	"fmt"
	"nli-go/lib/example1"
	"nli-go/lib/example2"
)

func TestSimpleRelationizer(test *testing.T) {

	rules := map[string][]example2.SimpleGrammarRule{
		"S": {
			// S(predication) = NP(entity) VP(predication)
			// subject(predication, entity)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"S", "NP", "VP"},
				EntityVariables:     []string{"Predication", "Entity", "Predication"},
				RelationTemplates: []example2.SimpleRelation{
					{Predicate: "subject", Arguments: []string{"Predication", "Entity"}},
				},
			},
		},
		"NP": {
			// NP(entity) = NBar(entity)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"NP", "NBar"},
				EntityVariables:     []string{"Entity", "Entity"},
				RelationTemplates:   []example2.SimpleRelation{},
			},
			// NP(entity) = DP(d1) NBar(entity)
			// determiner(entity, d1)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"NP", "DP", "NBar"},
				EntityVariables:     []string{"Entity", "Determiner", "Entity"},
				RelationTemplates: []example2.SimpleRelation{
					{Predicate: "determiner", Arguments: []string{"Entity", "Determiner"}},
				},
			},
		},
		"DP": {
			// DP(determiner) = det(determiner)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"DP", "det"},
				EntityVariables:     []string{"Determiner", "Determiner"},
				RelationTemplates:   []example2.SimpleRelation{},
			},
		},
		"NBar": {
			// NBar(entity = noun(entity)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"NBar", "noun"},
				EntityVariables:     []string{"Entity", "Entity"},
				RelationTemplates:   []example2.SimpleRelation{},
			},
		},
		"VP": {
			// VP(predication) = verb(predication) NP(entity)
			// object(predication, entity)
			example2.SimpleGrammarRule{
				SyntacticCategories: []string{"VP", "verb", "NP"},
				EntityVariables:     []string{"Predication", "Predication", "Entity"},
				RelationTemplates: []example2.SimpleRelation{
					{Predicate: "object", Arguments: []string{"Predication", "Entity"}},
				},
			},
		},
	}

	lexItems := map[string][]example2.SimpleLexItem{
		"all": {
			{PartOfSpeech: "det", RelationTemplates: []example2.SimpleRelation{
				{Predicate: "instance_of", Arguments: []string{"*", "all"}}},
			},
		},
		"horses": {
			{PartOfSpeech: "noun", RelationTemplates: []example2.SimpleRelation{
				{Predicate: "instance_of", Arguments: []string{"*", "horse"}},
				{Predicate: "number", Arguments: []string{"*", "plural"}},
			}},
		},
		"have": {
			{PartOfSpeech: "verb", RelationTemplates: []example2.SimpleRelation{
				{Predicate: "predication", Arguments: []string{"*", "have"}},
			}},
		},
		"hooves": {
			{PartOfSpeech: "noun", RelationTemplates: []example2.SimpleRelation{
				{Predicate: "instance_of", Arguments: []string{"*", "hoove"}},
				{Predicate: "number", Arguments: []string{"*", "plural"}},
			}},
		},
	}

	rawInput := "all horses have hooves"
	inputSource := example1.NewSimpleRawInputSource(rawInput)
	tokenizer := example1.NewSimpleTokenizer()
	grammar := example2.NewSimpleGrammar(rules)
	lexicon := example2.NewSimpleLexicon(lexItems)
	parser := example2.NewSimpleParser(grammar, lexicon)

	wordArray := tokenizer.Process(inputSource)
	parsedWords, relationList, ok := parser.Process(wordArray)

	if parsedWords != 4 {
		test.Error(fmt.Sprintf("Wrong number of words parsed: %d", parsedWords))
	}

	if !ok {
		test.Error("Parse failed")
	} else {

		if len(relationList) != 9 {
			test.Error(fmt.Sprintf("Wrong number of relations: %d", len(relationList)))
		}
		relationString := ""
		for i := 0; i < len(relationList); i++ {
			relationString += " " + relationList[i].ToString()
		}
		if relationString != " subject(S1, E1) determiner(E1, D1) instance_of(D1, all) instance_of(E1, horse) number(E1, plural) object(S1, E2) predication(S1, have) instance_of(E2, hoove) number(E2, plural)" {

			test.Error("Error in relations: " + relationString)
		}
	}
}

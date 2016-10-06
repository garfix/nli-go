package tests

import (
	"testing"
)

func TestSimpleRelationizer(test *testing.T) {

	//rules := map[string][]natlang.SimpleGrammarRule{
	//	"S": {
	//		// S(predication) = NP(entity) VP(predication)
	//		// subject(predication, entity)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"S", "NP", "VP"},
	//			EntityVariables:     []string{"Predication", "Entity", "Predication"},
	//			Sense: []mentalese.SimpleRelation{
	//				{Predicate: "subject", Arguments: []string{"Predication", "Entity"}},
	//			},
	//		},
	//	},
	//	"NP": {
	//		// NP(entity) = NBar(entity)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"NP", "NBar"},
	//			EntityVariables:     []string{"Entity", "Entity"},
	//			Sense:   []mentalese.SimpleRelation{},
	//		},
	//		// NP(entity) = DP(d1) NBar(entity)
	//		// determiner(entity, d1)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"NP", "DP", "NBar"},
	//			EntityVariables:     []string{"Entity", "Determiner", "Entity"},
	//			Sense: []mentalese.SimpleRelation{
	//				{Predicate: "determiner", Arguments: []string{"Entity", "Determiner"}},
	//			},
	//		},
	//	},
	//	"DP": {
	//		// DP(determiner) = det(determiner)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"DP", "det"},
	//			EntityVariables:     []string{"Determiner", "Determiner"},
	//			Sense:   []mentalese.SimpleRelation{},
	//		},
	//	},
	//	"NBar": {
	//		// NBar(entity = noun(entity)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"NBar", "noun"},
	//			EntityVariables:     []string{"Entity", "Entity"},
	//			Sense:   []mentalese.SimpleRelation{},
	//		},
	//	},
	//	"VP": {
	//		// VP(predication) = verb(predication) NP(entity)
	//		// object(predication, entity)
	//		natlang.SimpleGrammarRule{
	//			SyntacticCategories: []string{"VP", "verb", "NP"},
	//			EntityVariables:     []string{"Predication", "Predication", "Entity"},
	//			Sense: []mentalese.SimpleRelation{
	//				{Predicate: "object", Arguments: []string{"Predication", "Entity"}},
	//			},
	//		},
	//	},
	//}
	//
	//lexItems := map[string][]natlang.SimpleLexItem{
	//	"all": {
	//		{PartOfSpeech: "det", RelationTemplates: []mentalese.SimpleRelation{
	//			{Predicate: "instance_of", Arguments: []string{"*", "all"}}},
	//		},
	//	},
	//	"horses": {
	//		{PartOfSpeech: "noun", RelationTemplates: []mentalese.SimpleRelation{
	//			{Predicate: "instance_of", Arguments: []string{"*", "horse"}},
	//			{Predicate: "number", Arguments: []string{"*", "plural"}},
	//		}},
	//	},
	//	"have": {
	//		{PartOfSpeech: "verb", RelationTemplates: []mentalese.SimpleRelation{
	//			{Predicate: "predication", Arguments: []string{"*", "have"}},
	//		}},
	//	},
	//	"hooves": {
	//		{PartOfSpeech: "noun", RelationTemplates: []mentalese.SimpleRelation{
	//			{Predicate: "instance_of", Arguments: []string{"*", "hoove"}},
	//			{Predicate: "number", Arguments: []string{"*", "plural"}},
	//		}},
	//	},
	//}
	//
	//rawInput := "all horses have hooves"
	//inputSource := example1.NewSimpleRawInputSource(rawInput)
	//tokenizer := natlang.NewSimpleTokenizer()
	//grammar := natlang.NewSimpleGrammar()
	//for _, rule := range rules {
	//	grammar.AddRule(rule)
	//}
	//lexicon := natlang.NewSimpleLexicon()
	//for _, lexItem := range lexItems {
	//	lexicon.AddLexItem(lexItem)
	//}
	//parser := natlang.NewSimpleParser(grammar, lexicon)
	//
	//wordArray := tokenizer.Process(inputSource)
	//parsedWords, relationList, ok := parser.Process(wordArray)
	//
	//if parsedWords != 4 {
	//	test.Error(fmt.Sprintf("Wrong number of words parsed: %d", parsedWords))
	//}
	//
	//if !ok {
	//	test.Error("Parse failed")
	//} else {
	//
	//	if len(relationList) != 9 {
	//		test.Error(fmt.Sprintf("Wrong number of relations: %d", len(relationList)))
	//	}
	//	relationString := ""
	//	for i := 0; i < len(relationList); i++ {
	//		relationString += " " + relationList[i].String()
	//	}
	//	if relationString != " subject(S1, E1) determiner(E1, D1) instance_of(D1, all) instance_of(E1, horse) number(E1, plural) object(S1, E2) predication(S1, have) instance_of(E2, hoove) number(E2, plural)" {
	//
	//		test.Error("Error in relations: " + relationString)
	//	}
	//}
}

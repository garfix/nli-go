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
                EntityVariables: []string{"predication", "entity", "predication"},
                RelationTemplates: []example2.SimpleRelation{
                    {Predicate: "subject", Arguments: []string{"predication", "entity"}},
                },
            },
        },
        "NP": {
            // NP(entity) = NBar(entity)
            example2.SimpleGrammarRule{
                SyntacticCategories: []string{"NP", "NBar"},
                EntityVariables: []string{"entity", "entity"},
                RelationTemplates: []example2.SimpleRelation{},
            },
            // NP(entity) = DP(d1) NBar(entity)
            // determiner(entity, d1)
            example2.SimpleGrammarRule{
                SyntacticCategories: []string{"NP", "DP", "NBar"},
                EntityVariables: []string{"entity", "d1", "entity"},
                RelationTemplates: []example2.SimpleRelation{
                    {Predicate: "determiner", Arguments: []string{"entity", "d1"}},
                },
            },
        },
        "DP": {
            // DP(determiner) = det(determiner)
            example2.SimpleGrammarRule{
                SyntacticCategories: []string{"DP", "det"},
                EntityVariables: []string{"determiner", "determiner"},
                RelationTemplates: []example2.SimpleRelation{},
            },
        },
        "NBar": {
            // NBar(entity = noun(entity)
            example2.SimpleGrammarRule{
                SyntacticCategories: []string{"NBar", "noun"},
                EntityVariables: []string{"entity", "entity"},
                RelationTemplates: []example2.SimpleRelation{},
            },
        },
        "VP": {
            // VP(predication) = verb(predication) NP(entity)
            // object(predication, entity)
            example2.SimpleGrammarRule{
                SyntacticCategories: []string{"VP", "verb", "NP"},
                EntityVariables: []string{"predication", "predication", "entity"},
                RelationTemplates: []example2.SimpleRelation{
                    {Predicate: "object", Arguments: []string{"predication", "entity"}},
                },
            },
        },
    }

    lexItems := map[string][]example2.SimpleLexItem{
        "all": {
            {PartOfSpeech: "det", RelationTemplates: []example2.SimpleRelation{
                {Predicate: "instance-of", Arguments: []string{"*", "all"}}},
            },
        },
        "horses": {
            {PartOfSpeech: "noun", RelationTemplates: []example2.SimpleRelation{
                {Predicate: "instance-of", Arguments: []string{"*", "horse"}},
                {Predicate: "plural", Arguments: []string{"*"}},
            }},
        },
        "have": {
            {PartOfSpeech: "verb", RelationTemplates: []example2.SimpleRelation{
                {Predicate: "predication", Arguments: []string{"*", "have"}},
            }},
        },
        "hooves": {
            {PartOfSpeech: "noun", RelationTemplates: []example2.SimpleRelation{
                {Predicate: "instance-of", Arguments: []string{"*", "hoove"}},
                {Predicate: "plural", Arguments: []string{"*"}},
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

    if (parsedWords != 4) {
        test.Error(fmt.Sprintf("Wrong number of words parsed: %d", parsedWords))
    }

    if !ok {
        test.Error("Parse failed")
    } else {

        if len(relationList) != 9 {
            test.Error(fmt.Sprintf("Wrong number of relations: %d", len(relationList)))
        }
        relationString := "";
        for i := 0; i < len(relationList); i++ {
            relationString += " " + RelationToString(relationList[i])
        }
        if relationString != " subject(r1, e2) determiner(e2, d3) instance-of(d3, all) instance-of(e2, horse) plural(e2) object(r1, e4) predication(r1, have) instance-of(e4, hoove) plural(e4)" {
            test.Error("Error in relations: " + relationString)
        }
    }
}

func RelationToString(relation example2.SimpleRelation) string {
    text := relation.Predicate

    text += "("

    for i:= 0; i < len(relation.Arguments); i++ {

        if i > 0 {
            text += ", "
        }

        text += relation.Arguments[i]
    }

    text += ")"

    return text
}
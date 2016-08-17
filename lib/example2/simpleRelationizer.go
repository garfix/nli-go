package example2

import "fmt"

type simpleRelationizer struct {
    grammar *simpleGrammar
    lexicon *simpleLexicon
    variableIncrement int
}

func NewSimpleRelationizer(grammar *simpleGrammar, lexicon *simpleLexicon) *simpleRelationizer {
    return &simpleRelationizer{grammar: grammar, lexicon: lexicon, variableIncrement: 0}
}

func (relationizer *simpleRelationizer) GetNextVariable() string {
    relationizer.variableIncrement++
    return fmt.Sprintf("v%d", relationizer.variableIncrement)
}

// Creates meaningful relations for parseTreeRoot and its children
func (relationizer *simpleRelationizer) Process(parseTreeRoot SimpleParseTreeNode) []SimpleRelation {

    relations := relationizer.GetRelationsForNode(parseTreeRoot)
    return relations;
}

// Creates relations for parseTreeNode
func (relationizer *simpleRelationizer) GetRelationsForNode(parseTreeNode SimpleParseTreeNode) []SimpleRelation {

    children := parseTreeNode.Children
    relations := []SimpleRelation{}

    if len(children) > 0 {

        for i := 0; i < len(children); i++ {
            relations = append(relations, relationizer.GetRelationsForNode(children[i])...)
        }

    } else {
        relations = append(relations, relationizer.GetRelationsForLeafNode(parseTreeNode)...)
    }

    relations = append(relations, relationizer.GetRelationsForNonLeafNode(parseTreeNode)...)

    return relations
}

func (relationizer *simpleRelationizer) GetRelationsForNonLeafNode(parseTreeNode SimpleParseTreeNode) []SimpleRelation {

    syntacticCategories := []string{parseTreeNode.SyntacticCategory}
    for c := 0; c < len(parseTreeNode.Children); c++ {
        syntacticCategories = append(syntacticCategories, parseTreeNode.Children[c].SyntacticCategory)
    }

    rule, ok := relationizer.grammar.FindRule(syntacticCategories)
    if ok {
        relationTemplates := rule.RelationTemplates
        return relationTemplates
    } else {
        fmt.Print("Error in code!")
        fmt.Print(syntacticCategories)
    }
    return []SimpleRelation{}
}

func (relationizer *simpleRelationizer) GetRelationsForLeafNode(parseTreeNode SimpleParseTreeNode) []SimpleRelation {

    lexItem, ok := relationizer.lexicon.GetLexItem(parseTreeNode.Word, parseTreeNode.SyntacticCategory)
    if ok {

        relationTemplates := lexItem.RelationTemplates
        return relationTemplates

    } else {
        return []SimpleRelation{}
    }
}
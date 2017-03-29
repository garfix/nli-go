package tests

import (
    "testing"
    "nli-go/lib/importer"
    "nli-go/lib/mentalese"
)

func TestQuantifierScoper(t *testing.T) {

    tests := []struct {
        input string
        want string
    } {
        {
            "[have_child(S1, O1) quantification(S1, [ isa(S1, parent) ], D1, [ isa(D1, every) ]) quantification(O1, [ isa(O1, child) ], D2, [ isa(D2, 2) ])]",
            "[quant(S1, [isa(S1, parent)], D1, [isa(D1, every)], [quant(O1, [isa(O1, child)], D2, [isa(D2, 2)], [have_child(S1, O1)])])]",
        },
    }

    internalGrammarParser := importer.NewInternalGrammarParser()
    quantifierScoper := mentalese.NewQuantifierScoper()

    for _, test := range tests {

        input := internalGrammarParser.CreateRelationSet(test.input)
        result := quantifierScoper.Scope(input)

        if result.String() != test.want {
            t.Errorf("got %s, want %s", result.String(), test.want)
        }
    }
}



package tests
import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"testing"
)

func TestFunctions(t *testing.T) {

	log := common.NewSystemLog(false)
	functionBase := knowledge.NewSystemFunctionBase("name", log)
	parser := importer.NewInternalGrammarParser()
	tests := []struct {
		input      string
		binding     string
		wantBinding string
	}{
		{"split(W1, '-', S1, S2)", "{W1:'aap-noot'}", "{S1:'aap', S2:'noot', W1:'aap-noot'}"},
		{"join(W1, '-', S1, S2)", "{S1:'aap', S2:'noot'}", "{S1:'aap', S2:'noot', W1:'aap-noot'}"},
		{"concat(W1, S1, S2)", "{S1:'aap', S2:'noot'}", "{S1:'aap', S2:'noot', W1:'aapnoot'}"},
		{"greater_than(2, 1)", "{E1:1}", "{E1:1}"},
		{"greater_than(1, 2)", "{E1:1}", "{}"},
		{"less_than(E1, E2)", "{E1:1, E2:2}", "{E1:1, E2:2}"},
		{"add(E1, E2, S)", "{E1:1, E2:2}", "{E1:1, E2:2, S:'3'}"},
		{"subtract(E1, E2, S)", "{E1:1, E2:2}", "{E1:1, E2:2, S:'-1'}"},
		{"not_equals(E1, E2)", "{E1:1, E2:2}", "{E1:1, E2:2}"},
		{"not_equals(E1, E2)", "{E1:2, E2:2}", "{}"},
		{"equals(E1, E2)", "{E1:1, E2:2}", "{}"},
		{"equals(E1, E2)", "{E1:2, E2:2}", "{E1:2, E2:2}"},
		{"unify(quant(Q2, [], R2, []), quant(Q1, [], R1, []))", "{R2:5}", "{R1:5, R2:5}"},
		{"unify(X, 0)", "{Z:0}", "{X:0, Z:0}"},
		{"unify(0, Y)", "{Z:0}", "{Y:0, Z:0}"},
		{"date_subtract_years('2020-04-22', '1969-11-24', S)", "{}", "{S:'50'}"},
		{"assign(A, 2)", "{A:1}", "{A:2}"},
	}

	for _, test := range tests {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)
		wantBinding := parser.CreateBinding(test.wantBinding)

		resultBinding, _ := functionBase.Execute(input, binding)

		if !resultBinding.Equals(wantBinding) {
			t.Errorf("call %v with %v: got %v, want %v", input, binding, resultBinding, wantBinding)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

func TestAggregateFunctions(t *testing.T) {

	log := common.NewSystemLog(false)
	functionBase := knowledge.NewSystemAggregateBase("name", log)
	parser := importer.NewInternalGrammarParser()
	tests := []struct {
		input      string
		bindings     string
		wantBindings string
	}{
		{"number_of(W1, Number)", "[{W1:'aap'}{W1:'noot'}{W1:'noot'}]", "[{W1:'aap', Number:2}{W1:'noot', Number:2}{W1:'noot', Number:2}]"},
		{"number_of(W1, 2)", "[{W1:'aap'}{W1:'noot'}{W1:'noot'}]", "[{W1:'aap'}{W1:'noot'}{W1:'noot'}]"},
		{"number_of(W1, 3)", "[{W1:'aap'}{W1:'noot'}{W1:'noot'}]", "[]"},
		{"first(Name)", "[{A:1, Name:'Babbage'}{A:2, Name:'Charles B.'}{A:3, Name:'Charles Babbage'}]", "[{A:1, Name:'Babbage'}{A:2, Name:'Babbage'}{A:3, Name:'Babbage'}]"},
		{"exists()", "[{E1:1}{E1:2}]", "[{E1:1}{E1:2}]"},
		{"exists()", "[]", "[]"},
	}

	for _, test := range tests {

		input := parser.CreateRelation(test.input)
		bindings := parser.CreateBindings(test.bindings)
		wantBindings := parser.CreateBindings(test.wantBindings)

		resultBindings, _ := functionBase.Execute(input, bindings)

		if resultBindings.String() != wantBindings.String() {
			t.Errorf("call %v with %v: got %v, want %v", input, bindings, resultBindings, wantBindings)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

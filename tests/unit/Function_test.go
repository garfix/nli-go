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
		// keep extra bindings
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
		{"equals(E1, E2)", "{E1:2, E2:2}", "{E1:2, E2:2}"},
		{"date_subtract_years('2020-04-22', '1969-11-24', S)", "{}", "{S:'50'}"},
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

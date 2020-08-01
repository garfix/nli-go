package tests
import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/knowledge"
	"nli-go/lib/knowledge/nested"
	"nli-go/lib/mentalese"
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
		{"unify(quant(Q2, none, R2, none), quant(Q1, none, R1, none))", "{R2:5}", "{R1:5, R2:5}"},
		{"unify(X, 0)", "{Z:0}", "{X:0, Z:0}"},
		{"unify(0, Y)", "{Z:0}", "{Y:0, Z:0}"},
		{"date_subtract_years('2020-04-22', '1969-11-24', S)", "{}", "{S:'50'}"},
		{"date_subtract_years('2020-12-22', '1969-11-24', S)", "{}", "{S:'51'}"},
		{"date_subtract_years('2020-07-01', '2020-06-01', S)", "{}", "{S:'0'}"},
		{"date_subtract_years('2020-07-01', '2020-08-01', S)", "{}", "{S:'-1'}"},
		{"date_subtract_years('2020-07-01', '2021-01-01', S)", "{}", "{S:'-1'}"},
		{"date_subtract_years('2020-07-01', '2021-08-01', S)", "{}", "{S:'-2'}"},
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
	aggregateBase := knowledge.NewSystemAggregateBase("name", log)
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

		resultBindings, _ := aggregateBase.Execute(input, bindings)

		if resultBindings.String() != wantBindings.String() {
			t.Errorf("call %v with %v: got %v, want %v", input, bindings, resultBindings, wantBindings)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

func TestControlFunctions(t *testing.T) {

	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}

	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	functionBase := knowledge.NewSystemFunctionBase("name", log)
	solver.AddFunctionBase(functionBase)
	nestedBase := nested.NewSystemNestedStructureBase(solver, dialogContext, predicates, log)
	parser := importer.NewInternalGrammarParser()
	tests := []struct {
		input      string
		binding     string
		wantBindings string
	}{
		{"xor(_, unify(E, 1), unify(E, 2))", "{}", "[{E:1}]"},
		{"and(_, unify(E, 1), unify(E, 2))", "{}", "[]"},
		{"or(_, unify(E, 1), unify(E, 2))", "{}", "[{E:1} {E:2}]"},
		{"if_then_else(greater_than(6, 5), unify(E, 1), unify(E, 2))", "{X:3}", "[{E:1, X:3}]"},
		{"if_then_else(greater_than(5, 6), unify(E, 1), unify(E, 2))", "{X:3}", "[{E:2, X:3}]"},
	}

	for _, test := range tests {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)
		wantBindings := parser.CreateBindings(test.wantBindings)

		resultBindings := nestedBase.SolveNestedStructure(input, binding)

		if resultBindings.String() != wantBindings.String() {
			t.Errorf("call %v with %v: got %v, want %v", input, binding, resultBindings, wantBindings)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

func TestQuantFunctions(t *testing.T) {

	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	parser := importer.NewInternalGrammarParser()

	rules := parser.CreateRules(`[
		by_name(E1, E2, R) :- person(E1, Name1) person(E2, Name2) compare(Name1, Name2, R);
	]`)
	facts := parser.CreateRelationSet("" +
		"person(`:C`, 'Charles') " +
		"person(`:D`, 'Duncan') " +
		"person(`:B`, 'Bernhard') " +
		"person(`:E`, 'Edward') " +
		"person(`:A`, 'Abraham') ")
	ds2db := parser.CreateRules(`[
		person(E, Name) :- person(E, Name);
		person_named_abraham(E) :- person(E, 'Abraham');
		person_named_bernhard(E) :- person(E, 'Bernhard');
		person_named_edward(E) :- person(E, 'Edward');
	]`)
	ds2dbWrite := parser.CreateRules(`[]`)
	entities := mentalese.Entities{}

	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	factBase := knowledge.NewInMemoryFactBase("facts", facts, matcher, ds2db, ds2dbWrite, entities, log)
	solver.AddFactBase(factBase)
	functionBase := knowledge.NewSystemFunctionBase("name", log)
	solver.AddFunctionBase(functionBase)
	ruleBase := knowledge.NewInMemoryRuleBase("rules", rules, log)
	solver.AddRuleBase(ruleBase)
	nestedBase := nested.NewSystemNestedStructureBase(solver, dialogContext, predicates, log)
	tests := []struct {
		input      string
		binding     string
		wantBindings string
	}{
		{"quant_order(quant(Q1, E1), by_size, QSized)", "{}", "[{QSized: quant(Q1, E1, by_size)}]"},
		{"quant_order(or(_, quant(Q1, E1), quant(Q1, E1)), by_size, QSized)", "{}", "[{QSized: or(_, quant(Q1, E1, by_size), quant(Q1, E1, by_size))}]"},
		{`
			do(
				quant(
					quantifier(Result, Range, equals(Result, 3)),
					E,
					person(E, _)),
				List)`,
			"{}", "[{E: `:C`} {E: `:D`} {E: `:B`}]"},
		{`
			quant_ordered_list(
				quant(
					quantifier(Result, Range, equals(Result, 3)),
					E,
					person(E, _)),
				by_name,
				List)`,
			"{}", "[{List: [`:A`, `:B`, `:C`]}]"},
		{`
			quant_ordered_list(
				and(_,
					quant(
						quantifier(Result, Range, equals(Result, 3)),
						E,
						person(E, _)),
					quant(
						quantifier(Result, Range, equals(Result, 1)),
						E,
						person_named_edward(E))
				),
				by_name,
				List)`,
			"{}", "[{List: [`:A`, `:B`, `:C`, `:E`]}]"},
		{`
			quant_ordered_list(
				or(_,
					quant(
						quantifier(Result, Range, equals(Result, 3)),
						E,
						person(E, _)),
					quant(
						quantifier(Result, Range, equals(Result, 1)),
						E,
						person_named_bernhard(E))
				),
				by_name,
				List)`,
			"{}", "[{List: [`:B`]}]"},
		{`
			quant_ordered_list(
				or(_,
					quant(
						quantifier(Result, Range, equals(Result, 3)),
						E,
						person(E, _)),
					and(_,
						quant(
							quantifier(Result, Range, equals(Result, 3)),
							E,
							person_named_bernhard(E)),
						quant(
							quantifier(Result, Range, equals(Result, 1)),
							E,
							person_named_abraham(E))
					)
				),
				by_name,
				List)`,
			"{}", "[{List: [`:A`, `:B`]}]"},
		{"list_order([`:B`, `:C`, `:A`], by_name, Ordered)", "{}", "[{Ordered: [`:A`, `:B`, `:C`]}]"},
		{"list_foreach([`:B`, `:C`, `:A`], E, unify(F, E))", "{}", "[{E:`:B`, F:`:B`} {E:`:C`, F:`:C`} {E:`:A`, F:`:A`}]"},
	}

	for _, test := range tests {

		input := parser.CreateRelation(test.input)
		binding := parser.CreateBinding(test.binding)
		wantBindings := parser.CreateBindings(test.wantBindings)

		resultBindings := nestedBase.SolveNestedStructure(input, binding)

		if resultBindings.String() != wantBindings.String() {
			t.Errorf("got %v, want %v", resultBindings, wantBindings)
		}
	}

	if len(log.GetErrors()) > 0 {
		t.Errorf("errors: %v", log.String())
	}
}

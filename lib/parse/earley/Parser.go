package earley

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
	"nli-go/lib/parse"
)

// An implementation of Earley's top-down chart parsing algorithm as described in
// "Speech and Language Processing" (first edition) - Daniel Jurafsky & James H. Martin (Prentice Hall, 2000)
// It is the basic algorithm (p 381). Semantics (sense) is only calculated after the parse is complete.

type parser struct {
	grammar         *parse.Grammar
	lexicon         *parse.Lexicon
	senseBuilder    parse.SenseBuilder
}

func NewParser(grammar *parse.Grammar, lexicon *parse.Lexicon) *parser {
	return &parser{
		grammar: grammar,
		lexicon: lexicon,
		senseBuilder: parse.NewSenseBuilder(),
	}
}

// Parses words using parser.grammar and parser.lexicon
// Returns a sense, a parse tree, and a success flag
func (parser *parser) Parse(words []string) (mentalese.RelationSet, ParseTreeNode, bool) {

	common.LogTree("Parse", words);

	rootNode := ParseTreeNode{}
	sense := mentalese.RelationSet{}

	chart, ok := parser.buildChart(words)

	if ok {

		rootNode = parser.extractFirstTree(chart)
		sense = parser.extractFirstSense(chart)
	}

	common.LogTree("Parse", sense, rootNode, ok);

	return sense, rootNode, ok
}

// The body of Earley's algorithm
func (parser *parser) buildChart(words []string) (*chart, bool) {

	common.LogTree("createChart", words);

	chart := newChart(words)
	wordCount := len(words)

	initialState := newChartState(parse.NewGrammarRule([]string{"gamma", "s"}, []string{"g1", "s1"}, []mentalese.Relation{}), 1, 0, 0)
	parser.enqueue(chart, initialState, 0)

	// go through all word positions in the sentence
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a complete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			state := chart.states[i][j]

			// check if the entry is parsed completely
			if parser.isStateIncomplete(state) {

				// fetch the next consequent in the rule of the entry
				nextCat := parser.getNextCat(state)

				// is this an 'abstract' consequent like NP, VP, PP?
				if len(parser.grammar.FindRules(nextCat)) > 0 {

					// yes it is; add all entries that have this abstract consequent as their antecedent
					parser.predict(chart, state)

				} else if i < wordCount {

					// no it isn't, it is a low-level part-of-speech like noun, verb or adverb
					// if the current word in the sentence has this part-of-speech, then
					// we add a completed entry to the chart ($part-of-speech => $word)
					parser.scan(chart, state)
				}

			} else {

				// proceed all other entries in the chart that have this entry's antecedent as their next consequent
				treeComplete := parser.complete(chart, state)

				if treeComplete {

					common.LogTree("createChart", true);

					return chart, true
				}
			}
		}
	}

	common.LogTree("createChart", false);

	return chart, false
}

// Adds all entries to the chart that have the current consequent of $state as their antecedent.
func (parser *parser) predict(chart *chart, state chartState) {

	common.LogTree("predict", state);

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	endWordIndex := state.endWordIndex

	// go through all rules that have the next consequent as their antecedent
	for _, newRule := range parser.grammar.FindRules(nextConsequent) {

		predictedState := newChartState(newRule, 1, endWordIndex, endWordIndex)
		parser.enqueue(chart, predictedState, endWordIndex)
	}

	common.LogTree("predict");
}

// If the current consequent in state (which non-abstract, like noun, verb, adjunct) is one
// of the parts of speech associated with the current word in the sentence,
// then a new, completed, entry is added to the chart: (cat => word)
func (parser *parser) scan(chart *chart, state chartState) {

	common.LogTree("scan", state);

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	endWordIndex := state.endWordIndex
	endWord := chart.words[endWordIndex]

	_, lexItemFound := parser.lexicon.GetLexItem(endWord, nextConsequent)
	if lexItemFound {

		rule := parse.NewGrammarRule([]string{ nextConsequent, endWord }, []string{"a", "b"}, mentalese.RelationSet{})
		scannedState := newChartState(rule, 2, endWordIndex, endWordIndex + 1)
		parser.enqueue(chart, scannedState, endWordIndex + 1)
	}

	common.LogTree("scan", endWord, lexItemFound);
}

// This function is called whenever a state is completed.
// Its purpose is to advance other states.
//
// For example:
// - this state is NP -> noun, it has been completed
// - now proceed all other states in the chart that are waiting for an NP at the current position
func (parser *parser) complete(chart *chart, completedState chartState) bool {

	common.LogTree("complete", completedState);

	treeComplete := false;
	completedAntecedent := completedState.rule.GetAntecedent()
	for _, chartedState := range chart.states[completedState.startWordIndex] {

		dotPosition := chartedState.dotPosition
		rule := chartedState.rule

		// check if the antecedent of the completed state matches the charted state's consequent at the dot position
		if (dotPosition > rule.GetConsequentCount()) || (rule.GetConsequent(dotPosition - 1) != completedAntecedent) {
			continue;
		}

		advancedState := newChartState(rule, dotPosition + 1, chartedState.startWordIndex, completedState.endWordIndex)

		// store extra information to make it easier to extract parse trees later
		treeComplete, advancedState = parser.storeStateInfo(chart, completedState, chartedState, advancedState);
		if treeComplete {
			break;
		}

		parser.enqueue(chart, advancedState, completedState.endWordIndex)

		common.LogTree("complete");
    }

	return treeComplete
}


func (parser *parser) enqueue(chart *chart, state chartState, position int) {

	if !parser.isStateInChart(chart, state, position) {
		parser.pushState(chart, state, position)
	}
}

func (parser *parser) isStateIncomplete(state chartState) bool {

	return state.dotPosition < state.rule.GetConsequentCount() + 1
}

func (parser *parser) isStateInChart(chart *chart, state chartState, position int) bool {

	for _, presentState := range chart.states[position] {

		if  presentState.rule.Equals(state.rule) &&
			presentState.dotPosition == state.dotPosition &&
			presentState.startWordIndex == state.startWordIndex &&
			presentState.endWordIndex == state.endWordIndex {

			return true
		}

	}

	return false
}

func (parser *parser) pushState(chart *chart, state chartState, position int) {

	// index the state for later lookup
	chart.stateIdGenerator++
	state.id = chart.stateIdGenerator
	chart.indexedStates[state.id] = state

	chart.states[position] = append(chart.states[position], state)
}

func (parser *parser) getNextCat(state chartState) string {

	return state.rule.GetConsequent(state.dotPosition - 1)
}

func (parser *parser) storeStateInfo(chart *chart, completedState chartState, chartedState chartState, advancedState chartState) (bool, chartState) {

	treeComplete := false

	// store the state's "children" to ease building the parse trees from the packed forest
	advancedState.childStateIds = append(chartedState.childStateIds, completedState.id)

	// rule complete?
	if chartedState.dotPosition == chartedState.rule.GetConsequentCount() {

		// complete sentence?
		if chartedState.rule.GetAntecedent() == "gamma" {

			// that matches all words?
			if completedState.endWordIndex == len(chart.words) {

				chart.sentenceStates = append(chart.sentenceStates, advancedState)

				// set a flag to allow the parser to stop at the first complete parse
				treeComplete = true
			}
		}
	}

	return treeComplete, advancedState
}

func (parser *parser) extractFirstTree(chart *chart) ParseTreeNode {

	tree := ParseTreeNode{}

	if len(chart.sentenceStates) > 0 {

		rootStateId := chart.sentenceStates[0].childStateIds[0]
		root := chart.indexedStates[rootStateId]
		tree = parser.extractParseTreeBranch(chart, root)
	}

	return tree
}

func (parser *parser) extractParseTreeBranch(chart *chart, state chartState) ParseTreeNode {

	rule := state.rule
	branch := ParseTreeNode{ category: rule.GetAntecedent(), constituents: []ParseTreeNode{}, form: "" }

	if state.isLeafState() {

		branch.form = rule.GetConsequent(0)

	} else {

		for _, childStateId := range state.childStateIds {

			childState := chart.indexedStates[childStateId]
			branch.constituents = append(branch.constituents, parser.extractParseTreeBranch(chart, childState))
		}
	}

	return branch
}

func (parser *parser) extractFirstSense(chart *chart) mentalese.RelationSet {

	rootStateId := chart.sentenceStates[0].childStateIds[0]
	return parser.extractSenseFromState(chart, chart.indexedStates[rootStateId], parser.senseBuilder.GetNewVariable("Sentence"))
}

// Returns the sense of a state and its children
// state contains a rule with NP -> Det NBar
// antecedentVariable contains the actual variable used for the antecedent (for example: E1)
func (parser *parser) extractSenseFromState(chart *chart, state chartState, antecedentVariable string) mentalese.RelationSet {

	common.LogTree("extractSenseFromState", state, antecedentVariable)

	relations := mentalese.RelationSet{}
	rule := state.rule

	if state.isLeafState() {

		// leaf state rule: category -> word
		lexItem, _ := parser.lexicon.GetLexItem(state.rule.GetConsequent(0), state.rule.GetAntecedent())
		lexItemRelations := parser.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relations = append(relations, lexItemRelations...)

	} else {

		variableMap := parser.senseBuilder.CreateVariableMap(antecedentVariable, state.rule.EntityVariables)
		parentRelations := parser.senseBuilder.CreateGrammarRuleRelations(state.rule.Sense, variableMap)
		relations = append(relations, parentRelations...)

		// parse each of the children
		for i, _ := range rule.GetConsequents() {

			childStateId := state.childStateIds[i]
			childState := chart.indexedStates[childStateId]

			consequentVariable := variableMap[rule.EntityVariables[i + 1]]
			childRelations := parser.extractSenseFromState(chart, childState, consequentVariable)
			relations = append(relations, childRelations...)

		}
	}

	common.LogTree("extractSenseFromState", relations)

	return relations
}

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

	chart, ok := parser.createChart(words)

	if ok {

		rootNode = parser.extractFirstTree(chart)
		sense = parser.extractFirstSense(chart)
	}

	common.LogTree("Parse", sense, rootNode, ok);

	return sense, rootNode, ok
}

// The body of Earley's algorithm
func (parser *parser) createChart(words []string) (*chart, bool) {

	common.LogTree("createChart", words);

	chart := newChart(words)
	wordCount := len(words)

	parser.enqueue(chart, parser.createInitialState(words), 0)

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

func (parser *parser) createInitialState(words []string) chartState {

	rule := parse.GrammarRule{SyntacticCategories: []string{"gamma", "s"}, EntityVariables: []string{"g1", "s1"}, Sense: []mentalese.Relation{}}

	return newChartState(rule, 1, 0, 0)
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

		if
			presentState.rule.Equals(state.rule) &&
			presentState.dotPosition == state.dotPosition &&
			presentState.startWordIndex == state.startWordIndex &&
			presentState.endWordIndex == state.endWordIndex {

			return true
		}

	}
	return false
}

func (parser *parser) pushState(chart *chart, state chartState, position int) {

	chart.stateIdGenerator++
	state.id = chart.stateIdGenerator
	chart.treeInfoStates[state.id] = state
	chart.states[position] = append(chart.states[position], state)
}

func (parser *parser) getNextCat(state chartState) string {

	return state.rule.GetConsequent(state.dotPosition - 1)
}

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

func (parser *parser) scan(chart *chart, state chartState) {

	common.LogTree("scan", state);

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	endWordIndex := state.endWordIndex
	endWord := chart.words[endWordIndex]

	_, lexItemFound := parser.lexicon.GetLexItem(endWord, nextConsequent)
	if lexItemFound {

		rule := parse.GrammarRule{ SyntacticCategories: []string{ nextConsequent, endWord }, EntityVariables: []string{"a", "b"}, Sense: mentalese.RelationSet{} }
		scannedState := newChartState(rule, 2, endWordIndex, endWordIndex + 1)
		parser.enqueue(chart, scannedState, endWordIndex + 1)
	}

	common.LogTree("scan", endWord, lexItemFound);
}

func (parser *parser) complete(chart *chart, completedState chartState) bool {

	common.LogTree("complete", completedState);

	treeComplete := false;
	completedAntecedent := completedState.rule.GetAntecedent()
	for _, chartedState := range chart.states[completedState.startWordIndex] {

		dotPosition := chartedState.dotPosition
		rule := chartedState.rule
//
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

func (parser *parser) storeStateInfo(chart *chart, completedState chartState, chartedState chartState, advancedState chartState) (bool, chartState) {

	treeComplete := false

	// store the state's "children" to ease building the parse trees from the packed forest
	advancedState.children = append(chartedState.children, completedState.id)

	// rule complete?
	consequentCount := chartedState.rule.GetConsequentCount()
	if chartedState.dotPosition == consequentCount {

		// complete sentence?
		antecedent := chartedState.rule.GetAntecedent()

		if antecedent == "gamma" {

			// that matches all words?
			if completedState.endWordIndex == len(chart.words) {

				chart.treeInfoSentences = append(chart.treeInfoSentences, advancedState)

				// set a flag to allow the parser to stop at the first complete parse
				treeComplete = true
			}
		}
	}

	return treeComplete, advancedState
}

func (parser *parser) extractFirstTree(chart *chart) ParseTreeNode {

	tree := ParseTreeNode{}

	if len(chart.treeInfoSentences) > 0 {

		root := chart.treeInfoSentences[0]
		tree = parser.extractParseTreeBranch(chart, root)

	}

	return tree
}

func (parser *parser) extractParseTreeBranch(chart *chart, state chartState) ParseTreeNode {

	rule := state.rule

	antecedent := rule.GetAntecedent()
	if antecedent == "gamma" {

		constituentId := state.children[0]
		constituent := chart.treeInfoStates[constituentId]

		return parser.extractParseTreeBranch(chart, constituent)
	}

	branch := ParseTreeNode{category: antecedent, constituents: []ParseTreeNode{}, form: ""}

	if len(state.children) == 0 {

		branch.form = rule.GetConsequent(0)

	} else {

		for _, constituentId := range state.children {

			constituent := chart.treeInfoStates[constituentId]
			branch.constituents = append(branch.constituents, parser.extractParseTreeBranch(chart, constituent))
		}
	}

	return branch
}


func (parser *parser) extractFirstSense(chart *chart) mentalese.RelationSet {

	constituentId := chart.treeInfoSentences[0].children[0]
	return parser.extractSenseFromState(chart, chart.treeInfoStates[constituentId], parser.senseBuilder.GetNewVariable("Sentence"))
}

// Returns the sense of a state and its children
// state contains a rule with NP -> Det NBar
// antecedentVariable contains the actual variable used for the antecedent (for example: E1)
func (parser *parser) extractSenseFromState(chart *chart, state chartState, antecedentVariable string) mentalese.RelationSet {

	common.LogTree("extractSenseFromState", state, antecedentVariable)

	relations := mentalese.RelationSet{}
	variableMap := parser.senseBuilder.CreateVariableMap(antecedentVariable, state.rule.EntityVariables)
	rule := state.rule

	if len(state.children) > 0 {
		ruleRelations := parser.senseBuilder.CreateGrammarRuleRelations(state.rule.Sense, variableMap)
		relations = append(relations, ruleRelations...)

		// parse each of the children
		for i, _ := range rule.GetConsequents() {

			childConsequentId := state.children[i]
			childState := chart.treeInfoStates[childConsequentId]

			consequentVariable := variableMap[rule.EntityVariables[i + 1]]
			childRelations := parser.extractSenseFromState(chart, childState, consequentVariable)
			relations = append(relations, childRelations...)

		}

	} else {

		lexItem, _ := parser.lexicon.GetLexItem(state.rule.GetConsequent(0), state.rule.GetAntecedent())
		ruleRelations := parser.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relations = append(relations, ruleRelations...)
	}

	common.LogTree("extractSenseFromState", relations)

	return relations
}

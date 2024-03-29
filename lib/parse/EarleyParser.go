package parse

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"regexp"
	"strings"
)

// An implementation of Earley's top-down chart parsing algorithm as described in
// "Speech and Language Processing" (first edition) - Daniel Jurafsky & James H. Martin (Prentice Hall, 2000)
// It is the basic algorithm (p 381). Semantics (sense) is only calculated after the parse is isComplete.

type EarleyParser struct {
	grammarRules          *mentalese.GrammarRules
	morphologicalAnalyzer *MorphologicalAnalyzer
	log                   *common.SystemLog
}

type ParseResult struct {
	RootNodes []*mentalese.ParseTreeNode
	Error     string
	ErrorArg  string
}

func NewParser(grammarRules *mentalese.GrammarRules, log *common.SystemLog) *EarleyParser {
	return &EarleyParser{
		grammarRules:          grammarRules,
		morphologicalAnalyzer: nil,
		log:                   log,
	}
}

func (parser *EarleyParser) SetMorphologicalAnalyzer(morphologicalAnalyzer *MorphologicalAnalyzer) {
	parser.morphologicalAnalyzer = morphologicalAnalyzer
}

// Parses words using EarleyParser.grammar
// Returns parse tree roots
func (parser *EarleyParser) Parse(words []string, rootCategory string, rootVariables []string) ([]*mentalese.ParseTreeNode, ParseResult) {

	chart := parser.buildChart(parser.grammarRules, words, rootCategory, rootVariables)

	rootNodes := ExtractTreeRoots(chart)
	error := ""
	errorArg := ""

	if len(rootNodes) == 0 {

		nextWord := FindUnknownWord(chart)

		if nextWord != "" {
			error = common.UnknownWord
			errorArg = nextWord
		} else if len(words) == 0 {
			error = common.NoSentence
			errorArg = ""
		} else {
			error = common.NoUnderstandSentence
			errorArg = ""
		}
	}

	result := ParseResult{
		RootNodes: rootNodes,
		Error:     error,
		ErrorArg:  errorArg,
	}

	return rootNodes, result
}

// The body of Earley's algorithm
func (parser *EarleyParser) buildChart(grammarRules *mentalese.GrammarRules, words []string, rootCategory string, rootVariables []string) *chart {

	chart := NewChart(words, rootCategory, rootVariables)
	wordCount := len(words)

	chart.enqueue(chart.buildIncompleteGammaState(), 0)

	// go through all word positions in the sentence
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a isComplete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			state := chart.states[i][j]

			// check if the entry is parsed completely
			if !state.isComplete() {

				// add all entries that have this abstract consequent as their antecedent
				parser.predict(grammarRules, chart, state)

				// if the current word in the sentence has this part-of-speech, then
				// we add a completed entry to the chart (part-of-speech => word)
				if i < wordCount {
					parser.scan(chart, state)
				}

			} else {

				// proceed all other entries in the chart that have this entry's antecedent as their next consequent
				parser.complete(chart, state)
			}
		}
	}

	return chart
}

// Adds all entries to the chart that have the current consequent of $state as their antecedent.
func (parser *EarleyParser) predict(grammarRules *mentalese.GrammarRules, chart *chart, state chartState) {

	consequentIndex := state.dotPosition - 1
	nextConsequent := state.rule.GetConsequent(consequentIndex)
	nextConsequentVariables := state.rule.GetConsequentVariables(consequentIndex)
	endWordIndex := state.endWordIndex

	if parser.log.IsActive() {
		parser.log.AddDebug("predict", state.ToString(chart))
	}

	// go through all rules that have the next consequent as their antecedent
	for _, rule := range grammarRules.FindRules(nextConsequent, len(nextConsequentVariables)) {

		predictedState := newChartState(rule, 1, endWordIndex, endWordIndex)
		chart.enqueue(predictedState, endWordIndex)

		if parser.log.IsActive() {
			parser.log.AddDebug("> predicted", predictedState.ToString(chart))
		}
	}
}

// If the current consequent in state (which non-abstract, like noun, verb, adjunct) is one
// of the parts of speech associated with the current word in the sentence,
// then a new, completed, entry is added to the chart: (cat => word)
func (parser *EarleyParser) scan(chart *chart, state chartState) {

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	nextPosType := state.rule.GetConsequentPositionType(state.dotPosition - 1)
	nextVariables := state.rule.GetConsequentVariables(state.dotPosition - 1)
	endWordIndex := state.endWordIndex
	endWord := chart.words[endWordIndex]
	lexItemFound := false
	newPosType := mentalese.PosTypeRelation
	sense := mentalese.RelationSet{}

	if parser.log.IsActive() {
		parser.log.AddDebug("scan", state.ToString(chart))
	}

	// regular expression
	if nextPosType == mentalese.PosTypeRegExp {
		expression, err := regexp.Compile(nextConsequent)
		if err == nil {
			if expression.FindString(endWord) != "" {
				lexItemFound = true
				newPosType = mentalese.PosTypeRegExp
			}
		}
	}

	// proper noun
	if !lexItemFound && nextConsequent == mentalese.CategoryProperNoun {
		lexItemFound = true
	}

	// literal word form
	if !lexItemFound {
		if (nextConsequent == strings.ToLower(endWord)) && (len(nextVariables) == 0) {
			lexItemFound = true
			newPosType = mentalese.PosTypeWordForm
		}
	}

	// morphological analysis
	if !lexItemFound && nextPosType == mentalese.PosTypeRelation {
		if parser.morphologicalAnalyzer != nil {
			if !parser.grammarRules.WordOccurs(endWord, nextConsequent) {
				sense, lexItemFound = parser.morphologicalAnalyzer.Analyse(endWord, nextConsequent, nextVariables)
			}
		}
	}

	if lexItemFound {
		rule := mentalese.NewGrammarRule(
			[]string{newPosType, mentalese.PosTypeWordForm},
			[]string{nextConsequent, endWord},
			[][]string{nextVariables, {mentalese.Terminal}},
			sense)

		scannedState := newChartState(rule, 2, endWordIndex, endWordIndex+1)
		chart.enqueue(scannedState, endWordIndex+1)

		if parser.log.IsActive() {
			parser.log.AddDebug("> scanned", scannedState.ToString(chart)+" "+endWord)
		}
	}
}

// This function is called whenever a state is completed.
// Its purpose is to advance other states.
//
// For example:
// - this state is NP -> noun, it has been completed
// - now proceed all other states in the chart that are waiting for an NP at the current position
func (parser *EarleyParser) complete(chart *chart, completedState chartState) {

	completedAntecedent := completedState.rule.GetAntecedent()

	if parser.log.IsActive() {
		parser.log.AddDebug("complete", completedState.ToString(chart))
	}

	for _, chartedState := range chart.states[completedState.startWordIndex] {

		dotPosition := chartedState.dotPosition
		rule := chartedState.rule

		// check if the antecedent of the completed state matches the charted state's consequent at the dot position
		if (dotPosition > rule.GetConsequentCount()) || (rule.GetConsequent(dotPosition-1) != completedAntecedent) {
			continue
		}

		// check if the types match
		if chartedState.rule.GetConsequentPositionType(dotPosition-1) != completedState.rule.PositionTypes[0] {
			continue
		}

		// create a new state that is a dot-advancement of an older state
		advancedState := newChartState(rule, dotPosition+1, chartedState.startWordIndex, completedState.endWordIndex)

		// add this state to the index for tree extraction
		chart.updateAdvancedStatesIndex(completedState, advancedState)

		// enqueue the new state
		chart.enqueue(advancedState, completedState.endWordIndex)

		if parser.log.IsActive() {
			parser.log.AddDebug("> advanced", advancedState.ToString(chart))
		}
	}
}

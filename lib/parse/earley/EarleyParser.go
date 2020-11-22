package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"regexp"
	"strings"
)

// An implementation of Earley's top-down chart parsing algorithm as described in
// "Speech and Language Processing" (first edition) - Daniel Jurafsky & James H. Martin (Prentice Hall, 2000)
// It is the basic algorithm (p 381). Semantics (sense) is only calculated after the parse is complete.

type Parser struct {
	nameResolver *central.NameResolver
	meta         *mentalese.Meta
	log          *common.SystemLog
}

func NewParser(nameResolver *central.NameResolver, meta *mentalese.Meta, log *common.SystemLog) *Parser {
	return &Parser{
		nameResolver: nameResolver,
		meta:         meta,
		log:          log,
	}
}

// Parses words using Parser.grammar
// Returns parse tree roots
func (parser *Parser) Parse(grammarRules *parse.GrammarRules, words []string) []ParseTreeNode {

	if parser.log.Active() { parser.log.StartDebug("Parse", strings.Join(words, ",")) }

	chart := parser.buildChart(grammarRules, words)

	rootNodes := extractTreeRoots(chart)

	if len(rootNodes) == 0 {

		lastParsedWordIndex, nextWord := findLastCompletedWordIndex(chart)

		if nextWord != "" {
			parser.log.AddError("Incomplete. Could not parse word: " + nextWord)
		} else if len(words) == 0 {
			parser.log.AddError("No sentence given.")
		} else if lastParsedWordIndex == len(words)-1 {
			parser.log.AddError("All words are parsed but some word or token is missing to make the sentence complete.")
		}
	}

	if parser.log.Active() {
		str := ""
		for _, node := range rootNodes {
			str += " " + node.String()
		}
		parser.log.EndDebug("Parse", str)
	}

	return rootNodes
}

// The body of Earley's algorithm
func (parser *Parser) buildChart(grammarRules *parse.GrammarRules, words []string) (*chart) {

	if parser.log.Active() { parser.log.StartDebug("createChart", strings.Join(words, ", ")) }

	chart := newChart(words)
	wordCount := len(words)

	initialState := newChartState(chart.generateId(), parse.NewGrammarRule([]string{ parse.PosTypeRelation, parse.PosTypeRelation }, []string{"gamma", "s"}, [][]string{{"G"}, {"S"}}, mentalese.RelationSet{}), [][]string{{""}, {""}}, 1, 0, 0)

	if parser.log.Active() { parser.log.AddDebug("initial:", initialState.ToString(chart)) }

	chart.enqueue(initialState, 0)

	// go through all word positions in the sentence
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a complete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			state := chart.states[i][j]

			if parser.log.Active() { parser.log.AddDebug("do:", state.ToString(chart)) }

			// check if the entry is parsed completely
			if !state.complete() {

				// note: we make no distinction between part-of-speech and not part-of-speech; a category can be both

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

	if parser.log.Active() { parser.log.EndDebug("createChart", "") }

	return chart
}

// Adds all entries to the chart that have the current consequent of $state as their antecedent.
func (parser *Parser) predict(grammarRules *parse.GrammarRules, chart *chart, state chartState) {

	consequentIndex := state.dotPosition - 1
	nextConsequent := state.rule.GetConsequent(consequentIndex)
	nextConsequentVariables := state.rule.GetConsequentVariables(consequentIndex)
	endWordIndex := state.endWordIndex

	// go through all rules that have the next consequent as their antecedent
	for _, rule := range grammarRules.FindRules(nextConsequent, len(nextConsequentVariables)) {

		parentSSelection := state.sSelection[consequentIndex + 1]
		sSelection, allowed := combineSSelection(parser.meta, parentSSelection, rule)
		if !allowed {
			continue
		}

		predictedState := newChartState(chart.generateId(), rule, sSelection, 1, endWordIndex, endWordIndex)
		chart.enqueue(predictedState, endWordIndex)

		if parser.log.Active() { parser.log.AddDebug("predict", predictedState.ToString(chart)) }
	}
}

// If the current consequent in state (which non-abstract, like noun, verb, adjunct) is one
// of the parts of speech associated with the current word in the sentence,
// then a new, completed, entry is added to the chart: (cat => word)
func (parser *Parser) scan(chart *chart, state chartState) {

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	endWordIndex := state.endWordIndex
	endWord := chart.words[endWordIndex]
	nameInformations := []central.NameInformation{}
	lexItemFound := false

	if state.rule.GetConsequentPositionType(state.dotPosition - 1) == parse.PosTypeRegExp {
		expression, err := regexp.Compile(nextConsequent)
		if err == nil {
			if expression.FindString(endWord) != "" {
				lexItemFound = true
			}
		}
	}

	if !lexItemFound && nextConsequent == mentalese.CategoryProperNoun {
		lexItemFound, nameInformations = parser.isProperNoun(chart, state)
	}

	if !lexItemFound {
		if
		(nextConsequent == strings.ToLower(endWord)) &&
		(len(state.rule.GetConsequentVariables(state.dotPosition - 1)) == 0) {
			lexItemFound = true
		}
	}

	if lexItemFound {

		rule := parse.NewGrammarRule([]string{ parse.PosTypeRelation, parse.PosTypeWordForm }, []string{nextConsequent, endWord}, [][]string{{"a"}, {"b"}}, mentalese.RelationSet{})
		sType := state.sSelection[state.dotPosition - 1]
		scannedState := newChartState(chart.generateId(), rule, parse.SSelection{sType, sType}, 2, endWordIndex, endWordIndex+1)
		scannedState.nameInformations = nameInformations
		chart.enqueue(scannedState, endWordIndex+1)

		if parser.log.Active() { parser.log.AddDebug("scanned", scannedState.ToString(chart) + " " + endWord) }
	}
}

func (parser *Parser) isProperNoun(chart *chart, state chartState) (bool, []central.NameInformation) {

	wordIndex := state.endWordIndex
	sType := state.sSelection[state.dotPosition - 1 + 1]
	wordCount := len(state.rule.GetConsequents())

	// if the first consequent has created a match, all following words match
	if state.dotPosition > 1 {
		return true, []central.NameInformation{}
	}

	// check if it is possible to match all words in the remainder of the sentence
	if wordIndex + wordCount > len(chart.words) {
		return false, []central.NameInformation{}
	}

	// first word in proper noun consequents?  try to match all words at once
	words := chart.words[wordIndex:wordIndex + wordCount]
	wordString := strings.Join(words, " ")
	nameInformations := parser.nameResolver.ResolveName(wordString, sType[0])

	if len(nameInformations) > 0 {
		return true, nameInformations
	}

	return false, []central.NameInformation{}
}

// This function is called whenever a state is completed.
// Its purpose is to advance other states.
//
// For example:
// - this state is NP -> noun, it has been completed
// - now proceed all other states in the chart that are waiting for an NP at the current position
func (parser *Parser) complete(chart *chart, completedState chartState) {

	completedAntecedent := completedState.rule.GetAntecedent()
	for _, chartedState := range chart.states[completedState.startWordIndex] {

		dotPosition := chartedState.dotPosition
		rule := chartedState.rule
		sSelection := chartedState.sSelection

		// check if the antecedent of the completed state matches the charted state's consequent at the dot position
		if (dotPosition > rule.GetConsequentCount()) || (rule.GetConsequent(dotPosition-1) != completedAntecedent) {
			continue
		}

		advancedState := newChartState(chart.generateId(), rule, sSelection, dotPosition+1, chartedState.startWordIndex, completedState.endWordIndex)
		advancedState.children = append(parser.copyChildArray(chartedState.children), completedState)

		if advancedState.complete() {
			chart.indexChildren(advancedState)
		}

		if parser.log.Active() { parser.log.AddDebug("advanced", advancedState.ToString(chart)) }

		f := chart.enqueue(advancedState, completedState.endWordIndex)

		if parser.log.Active() {
			str := "no"
			if f { str = "yes" }
			parser.log.AddDebug("found", str)
		}
	}
}

func (parser *Parser) copyChildArray(children []chartState) []chartState {
	newArray := []chartState{}
	for _, state := range children {
		newArray = append(newArray, state)
	}
	return newArray
}


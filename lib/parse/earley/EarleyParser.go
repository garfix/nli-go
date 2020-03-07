package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strings"
)

// An implementation of Earley's top-down chart parsing algorithm as described in
// "Speech and Language Processing" (first edition) - Daniel Jurafsky & James H. Martin (Prentice Hall, 2000)
// It is the basic algorithm (p 381). Semantics (sense) is only calculated after the parse is complete.

const ProperNounCategory = "proper_noun"

type Parser struct {
	grammar      *parse.Grammar
	lexicon      *parse.Lexicon
	nameResolver *central.NameResolver
	predicates   mentalese.Predicates
	log          *common.SystemLog
}

func NewParser(grammar *parse.Grammar, lexicon *parse.Lexicon, nameResolver *central.NameResolver, predicates mentalese.Predicates, log *common.SystemLog) *Parser {
	return &Parser{
		grammar:      grammar,
		lexicon:      lexicon,
		nameResolver: nameResolver,
		predicates:   predicates,
		log:          log,
	}
}

// Parses words using Parser.grammar and Parser.lexicon
// Returns parse tree roots
func (parser *Parser) Parse(words []string) []ParseTreeNode {

	parser.log.StartDebug("Parse", words)

	chart := parser.buildChart(words)

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

	parser.log.EndDebug("Parse", rootNodes)

	return rootNodes
}

// The body of Earley's algorithm
func (parser *Parser) buildChart(words []string) (*chart) {

	parser.log.StartDebug("createChart", words)

	chart := newChart(words)
	wordCount := len(words)

	initialState := newChartState(chart.generateId(), parse.NewGrammarRule([]string{"gamma", "s"}, [][]string{{"G"}, {"S"}}, mentalese.RelationSet{}), [][]string{{""}, {""}}, 1, 0, 0)
	parser.log.EndDebug("initial:", initialState.ToString(chart))
	chart.enqueue(initialState, 0)

	// go through all word positions in the sentence
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a complete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			state := chart.states[i][j]

			parser.log.EndDebug("do:", state.ToString(chart))

			// check if the entry is parsed completely
			if state.isIncomplete() {

				// note: we make no distinction between part-of-speech and not part-of-speech; a category can be both

				// add all entries that have this abstract consequent as their antecedent
				parser.predict(chart, state)

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

	parser.log.EndDebug("createChart")

	return chart
}

// Adds all entries to the chart that have the current consequent of $state as their antecedent.
func (parser *Parser) predict(chart *chart, state chartState) {

	consequentIndex := state.dotPosition - 1
	nextConsequent := state.rule.GetConsequent(consequentIndex)
	nextConsequentVariables := state.rule.GetConsequentVariables(consequentIndex)
	endWordIndex := state.endWordIndex

	// go through all rules that have the next consequent as their antecedent
	for _, rule := range parser.grammar.FindRules(nextConsequent, len(nextConsequentVariables)) {

		parentSSelection := state.sSelection[consequentIndex + 1]
		sSelection, allowed := combineSSelection(parser.predicates, parentSSelection, rule)
		if !allowed {
			continue
		}

		predictedState := newChartState(chart.generateId(), rule, sSelection, 1, endWordIndex, endWordIndex)
		chart.enqueue(predictedState, endWordIndex)

		parser.log.EndDebug("predict:", predictedState.ToString(chart))
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

	_, lexItemFound, _ := parser.lexicon.GetLexItem(endWord, nextConsequent)
	if !lexItemFound && nextConsequent == ProperNounCategory {
		lexItemFound, nameInformations = parser.isProperNoun(chart, state)
	}
	if lexItemFound {

		rule := parse.NewGrammarRule([]string{nextConsequent, endWord}, [][]string{{"a"}, {"b"}}, mentalese.RelationSet{})
		sType := state.sSelection[state.dotPosition - 1]
		scannedState := newChartState(chart.generateId(), rule, parse.SSelection{sType, sType}, 2, endWordIndex, endWordIndex+1)
		scannedState.nameInformations = nameInformations
		chart.enqueue(scannedState, endWordIndex+1)

		parser.log.EndDebug("scanned:", scannedState.ToString(chart), endWord)
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
		advancedState.parentIds = append(common.IntArrayCopy(chartedState.parentIds), completedState.id)

		parser.log.EndDebug("advanced:", advancedState.ToString(chart))

		f := chart.enqueue(advancedState, completedState.endWordIndex)

		parser.log.EndDebug("found:", f)
	}
}

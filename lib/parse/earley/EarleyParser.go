package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"sort"
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
	relationizer EarleyRelationizer
	log          *common.SystemLog
}

func NewParser(grammar *parse.Grammar, lexicon *parse.Lexicon, nameResolver *central.NameResolver, predicates mentalese.Predicates, log *common.SystemLog) *Parser {
	return &Parser{
		grammar:      grammar,
		lexicon:      lexicon,
		nameResolver: nameResolver,
		predicates:   predicates,
		relationizer: EarleyRelationizer{},
		log:          log,
	}
}

// Parses words using Parser.grammar and Parser.lexicon
// Returns a sense, a parse tree, and a success flag
func (parser *Parser) Parse(words []string) ParseTreeNode {

	parser.log.StartDebug("Parse", words)

	rootNode := ParseTreeNode{}

	chart, ok := parser.buildChart(words)

	if ok {

		rootNode = parser.extractFirstTree(chart)

	} else {

		lastParsedWordIndex, nextWord := parser.findLastCompletedWordIndex(chart)

		if nextWord != "" {
			parser.log.AddError("Incomplete. Could not parse word: " + nextWord)
		} else if len(words) == 0 {
			parser.log.AddError("No sentence given.")
		} else if lastParsedWordIndex == len(words)-1 {
			parser.log.AddError("All words are parsed but some word or token is missing to make the sentence complete.")
		}
	}

	parser.log.EndDebug("Parse", rootNode, ok)

	return rootNode
}

// For a given sequence of words, make suggestions for the next word, and return these in sorted order
func (parser *Parser) Suggest(words []string) []string {

	suggests := []string{}
	position := len(words)

	chart, ok := parser.buildChart(words)

	if ok {

	} else if len(chart.states) < position {

	} else {

		for _, state := range chart.states[position] {
			if parser.isStateIncomplete(state) {
				category := state.rule.SyntacticCategories[state.dotPosition]
				suggests = append(suggests, parser.lexicon.GetWordForms(category)...)
			}
		}
	}

	suggests = common.StringArrayDeduplicate(suggests)
	sort.Strings(suggests)
	return suggests
}

// The body of Earley's algorithm
func (parser *Parser) buildChart(words []string) (*chart, bool) {

	parser.log.StartDebug("createChart", words)

	chart := newChart(words)
	wordCount := len(words)

	initialState := newChartState(parse.NewGrammarRule([]string{"gamma", "s"}, []string{"g1", "s1"}, mentalese.RelationSet{}), parse.SSelection{"", ""}, 1, 0, 0)
	parser.enqueue(chart, initialState, 0)

	// go through all word positions in the sentence
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a complete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			state := chart.states[i][j]

			// check if the entry is parsed completely
			if parser.isStateIncomplete(state) {

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
				treeComplete := parser.complete(chart, state)

				if treeComplete {

					parser.log.EndDebug("createChart", true)

					return chart, true
				}
			}
		}
	}

	parser.log.EndDebug("createChart", false)

	return chart, false
}

// Adds all entries to the chart that have the current consequent of $state as their antecedent.
func (parser *Parser) predict(chart *chart, state chartState) {

	parser.log.StartDebug("predict", state)

	consequentIndex := state.dotPosition - 1
	nextConsequent := state.rule.GetConsequent(consequentIndex)
	endWordIndex := state.endWordIndex

	// go through all rules that have the next consequent as their antecedent
	for _, rule := range parser.grammar.FindRules(nextConsequent) {

		parentSSelection := state.sSelection[consequentIndex + 1]
		sSelection, allowed := parser.relationizer.CombineSSelection(parser.predicates, parentSSelection, rule)
		if !allowed {
			continue
		}

		predictedState := newChartState(rule, sSelection, 1, endWordIndex, endWordIndex)
		parser.enqueue(chart, predictedState, endWordIndex)
	}

	parser.log.EndDebug("predict")
}

func (parser *Parser) collectEntityVariableTypes(state chartState) []string {
	return []string{}
}

// If the current consequent in state (which non-abstract, like noun, verb, adjunct) is one
// of the parts of speech associated with the current word in the sentence,
// then a new, completed, entry is added to the chart: (cat => word)
func (parser *Parser) scan(chart *chart, state chartState) {

	parser.log.StartDebug("scan", state)

	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)
	nextConsequentVariable := state.rule.GetConsequentVariable(state.dotPosition - 1)
	endWordIndex := state.endWordIndex
	endWord := chart.words[endWordIndex]
	nameInformations := []central.NameInformation{}

	_, lexItemFound, _ := parser.lexicon.GetLexItem(endWord, nextConsequent)
	if !lexItemFound && nextConsequent == ProperNounCategory {
		lexItemFound, nameInformations = parser.isProperNoun(chart, endWordIndex, nextConsequentVariable, state.sSelection[state.dotPosition - 1 + 1])
		lexItemFound = lexItemFound
	}
	if lexItemFound {

		rule := parse.NewGrammarRule([]string{nextConsequent, endWord}, []string{"a", "b"}, mentalese.RelationSet{})
		sType := state.sSelection[state.dotPosition - 1]
		scannedState := newChartState(rule, parse.SSelection{sType, sType}, 2, endWordIndex, endWordIndex+1)
		scannedState.nameInformations = nameInformations
		parser.enqueue(chart, scannedState, endWordIndex+1)
	}

	parser.log.EndDebug("scan", endWord, lexItemFound)
}

func (parser *Parser) isProperNoun(chart *chart, wordIndex int, nextConsequentVariable string, sType string) (bool, []central.NameInformation) {

	word := chart.words[wordIndex]

	// check memory
	for startIndex, sequences := range chart.properNounSequences {
		for _, sequence := range sequences {
			endIndex := startIndex + len(sequence)
			if wordIndex >= startIndex && wordIndex < endIndex {
				sequenceWordIndex := wordIndex - startIndex
				if word == sequence[sequenceWordIndex] {
					return true, []central.NameInformation{}
				}
			}
		}
	}

	// prime memory
	for endIndex := len(chart.words); endIndex > wordIndex; endIndex-- {
// todo max number of words should be deduced from grammar rules

		words := chart.words[wordIndex:endIndex]
		wordString := strings.Join(words, " ")
		nameInformations := parser.nameResolver.ResolveName(wordString, sType)

		if len(nameInformations) > 0 {
			_, ok := chart.properNounSequences[wordIndex]
			if !ok {
				chart.properNounSequences[wordIndex] = [][]string{}
			}
			chart.properNounSequences[wordIndex] = append(chart.properNounSequences[wordIndex], words)
			return true, nameInformations
		}
	}

	return false, []central.NameInformation{}
}

// This function is called whenever a state is completed.
// Its purpose is to advance other states.
//
// For example:
// - this state is NP -> noun, it has been completed
// - now proceed all other states in the chart that are waiting for an NP at the current position
func (parser *Parser) complete(chart *chart, completedState chartState) bool {

	parser.log.StartDebug("complete", completedState)

	treeComplete := false
	completedAntecedent := completedState.rule.GetAntecedent()
	for _, chartedState := range chart.states[completedState.startWordIndex] {

		dotPosition := chartedState.dotPosition
		rule := chartedState.rule
		sSelection := chartedState.sSelection

		// check if the antecedent of the completed state matches the charted state's consequent at the dot position
		if (dotPosition > rule.GetConsequentCount()) || (rule.GetConsequent(dotPosition-1) != completedAntecedent) {
			continue
		}

		advancedState := newChartState(rule, sSelection, dotPosition+1, chartedState.startWordIndex, completedState.endWordIndex)

		// store extra information to make it easier to extract parse trees later
		treeComplete, advancedState = parser.storeStateInfo(chart, completedState, chartedState, advancedState)
		if treeComplete {
			break
		}

		parser.enqueue(chart, advancedState, completedState.endWordIndex)
	}

	parser.log.EndDebug("complete")

	return treeComplete
}

func (parser *Parser) enqueue(chart *chart, state chartState, position int) {

	if !parser.isStateInChart(chart, state, position) {
		parser.pushState(chart, state, position)
	}
}

func (parser *Parser) isStateIncomplete(state chartState) bool {

	return state.dotPosition < state.rule.GetConsequentCount()+1
}

func (parser *Parser) isStateInChart(chart *chart, state chartState, position int) bool {

	for _, presentState := range chart.states[position] {

		if presentState.rule.Equals(state.rule) &&
			presentState.dotPosition == state.dotPosition &&
			presentState.startWordIndex == state.startWordIndex &&
			presentState.endWordIndex == state.endWordIndex {

			return true
		}
	}

	return false
}

func (parser *Parser) pushState(chart *chart, state chartState, position int) {

	// index the state for later lookup
	chart.stateIdGenerator++
	state.id = chart.stateIdGenerator
	chart.indexedStates[state.id] = state

	chart.states[position] = append(chart.states[position], state)
}

func (parser *Parser) getNextCat(state chartState) string {

	return state.rule.GetConsequent(state.dotPosition - 1)
}

func (parser *Parser) storeStateInfo(chart *chart, completedState chartState, chartedState chartState, advancedState chartState) (bool, chartState) {

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

				// set a flag to allow the Parser to stop at the first complete parse
				treeComplete = true
			}
		}
	}

	return treeComplete, advancedState
}

func (parser *Parser) extractFirstTree(chart *chart) ParseTreeNode {

	tree := ParseTreeNode{}

	if len(chart.sentenceStates) > 0 {

		rootStateId := chart.sentenceStates[0].childStateIds[0]
		root := chart.indexedStates[rootStateId]
		tree = parser.extractParseTreeBranch(chart, root)
	}

	return tree
}

func (parser *Parser) extractParseTreeBranch(chart *chart, state chartState) ParseTreeNode {

	rule := state.rule
	branch := ParseTreeNode{category: rule.GetAntecedent(), constituents: []ParseTreeNode{}, form: "", rule: state.rule, nameInformations: state.nameInformations}

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

// Returns the word that could not be parsed (or ""), and the index of the last completed word
func (parser *Parser) findLastCompletedWordIndex(chart *chart) (int, string) {

	nextWord := ""
	lastIndex := -1

	// find the last completed nextWord

	for i := len(chart.states) - 1; i >= 0; i-- {
		states := chart.states[i]
		for _, state := range states {
			if !parser.isStateIncomplete(state) {

				lastIndex = state.endWordIndex - 1
				goto done
			}
		}
	}

done:

	if lastIndex <= len(chart.words)-2 {
		nextWord = chart.words[lastIndex+1]
	}

	return lastIndex, nextWord
}

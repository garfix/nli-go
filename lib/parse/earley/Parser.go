package earley

/**
 * An implementation of Earley's top-down chart parsing algorithm as described in
 * "Speech and Language Processing" (first edition) - Daniel Jurafsky & James H. Martin (Prentice Hall, 2000)
 * It is the basic algorithm (p 381) extended with semantics (p 570)
 */

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"nli-go/lib/common"
)

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

// Parses tokens using parser.grammar and parser.lexicon
func (parser *parser) Parse(words []string) (mentalese.RelationSet, ParseTreeNode, bool) {

	common.LogTree("Parse", words);

	rootNode := ParseTreeNode{}
	set := mentalese.RelationSet{}

	chart, ok := parser.createChart(words)

	if ok {


		rootNode = parser.extractFirstTree(chart)

		//length, relationList, ok := parser.parseAllRules("s", tokens, 0, parser.getNewVariable("Sentence"))
		//
		//set = append(set, relationList...)

		common.LoggerActive = false

		set = parser.extractFirstSense(chart)

		common.LoggerActive = false

	}

	common.LogTree("Parse", set, rootNode, ok);

	return set, rootNode, ok
}

func (parser *parser) createChart(words []string) (*chart, bool) {

	common.LogTree("createChart", words);

	chart := newChart(words)

	parser.enqueue(chart, parser.createInitialState(words), 0)

	// go through all word positions in the sentence
	// $wordCount = count($this->words);
	wordCount := len(words)

	//for ($i = 0; $i <= $wordCount; $i++) {
	for i := 0; i <= wordCount; i++ {

		// go through all chart entries in this position (entries will be added while we're in the loop)
		// for ($j = 0; $j < count($this->chart[$i]); $j++) {
		for j := 0; j < len(chart.states[i]); j++ {

			// a state is a complete entry in the chart (rule, dotPosition, startWordIndex, endWordIndex)
			// $state = $this->chart[$i][$j];
			state := chart.states[i][j]

			// check if the entry is parsed completely
			// if ($this->isIncomplete($state)) {
			if parser.isStateIncomplete(state) {

				// fetch the next consequent in the rule of the entry
				// $nextCat = $this->getNextCat($state);
				nextCat := parser.getNextCat(state)

				// is this an 'abstract' consequent like NP, VP, PP?
				// NB if (!$this->Grammar->isPartOfSpeech($nextCat)) {
				if len(parser.grammar.FindRules(nextCat)) > 0 {

					// yes it is; add all entries that have this abstract consequent as their antecedent
					// $this->predict($state);
					parser.predict(chart, state)

					// } elseif ($i < $wordCount) {
				} else if i < wordCount {

					// no it isn't, it is a low-level part-of-speech like noun, verb or adverb
					// if the current word in the sentence has this part-of-speech, then
					// we add a completed entry to the chart ($part-of-speech => $word)
					// $this->scan($state);
					parser.scan(chart, state)
				}

			} else {

				// NB $this->lastParsedIndex = $i;

				// proceed all other entries in the chart that have this entry's antecedent as their next consequent
				// $treeComplete = $this->complete($state);
				treeComplete := parser.complete(chart, state)

				// NB if ($this->singleTree && $treeComplete) {
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

	//productionRule := productionRule{ antecedent: "gamma", consequents: []string{"s"}, antecedentCategory: "", consequentCategories: []string{""}}
	//rule := rule{production: productionRule, sense:mentalese.RelationSet{}}
	rule := parse.GrammarRule{SyntacticCategories: []string{"gamma", "s"}, EntityVariables: []string{"g1", "s1"}, Sense: []mentalese.Relation{}}
	initialState := newChartState(rule, 1, 0, 0)

	return initialState
}

func (parser *parser) enqueue(chart *chart, state chartState, position int) {

	ok := false

	// check for completeness
	// if ($this->isIncomplete($state)) {
	if parser.isStateIncomplete(state) {

		// if (!$this->isStateInChart($state, $position)) {
		if !parser.isStateInChart(chart, state, position) {

			// $this->pushState($state, $position);
			parser.pushState(chart, state, position)
		}

	// } elseif ($this->unifyState($state)) {
	} else if parser.unifyState(chart, state) {

		// if ($this->applySemantics($state)) {
		state, ok = parser.applySense(chart, state)
		if ok {

			// if (!$this->isStateInChart($state, $position)) {
			if !parser.isStateInChart(chart, state, position) {

				// $this->pushState($state, $position);
				parser.pushState(chart, state, position)
			}
		}
	}
}

func (parser *parser) isStateIncomplete(state chartState) bool {

	//        $consequentCount = $state['rule']->getProduction()->getConsequentCount() + 1;
	//
	//        return ($state['dotPosition'] < $consequentCount);
	return state.dotPosition < state.rule.GetConsequentCount() + 1
}

func (parser *parser) isStateInChart(chart *chart, state chartState, position int) bool {

	//        foreach ($this->chart[$position] as $presentState) {
	for _, presentState := range chart.states[position] {
	//            if (
	//                $presentState['rule'] == $state['rule'] &&
	//                $presentState['dotPosition'] == $state['dotPosition'] &&
	//                $presentState['startWordIndex'] == $state['startWordIndex'] &&
	//                $presentState['endWordIndex'] == $state['endWordIndex']
	//            ) {
	//                return true;
	//            }
		if
			presentState.rule.Equals(state.rule) &&
			presentState.dotPosition == state.dotPosition &&
			presentState.startWordIndex == state.startWordIndex &&
			presentState.endWordIndex == state.endWordIndex {

			return true
		}

	}
	//        return false;
	return false
}

func (parser *parser) pushState(chart *chart, state chartState, position int) {

	//        static $stateIDs = 0;
	//
	//        $this->showDebug('enqueue', $state);
	//
	//        $stateIDs++;
	chart.stateIdGenerator++
	//
	//        $state['id'] = $stateIDs;
	state.id = chart.stateIdGenerator
	//        $this->treeInfo['states'][$stateIDs] = $state;
	chart.treeInfoStates[state.id] = state
	//        $this->chart[$position][] = $state;
	chart.states[position] = append(chart.states[position], state)
}

/**
 * Reserved for later use: check for semantic conflicts.
 */
func (parser *parser) unifyState(chart *chart, state chartState) bool {
	return true
}

func (parser *parser) applySense(chart *chart, state chartState) (chartState, bool) {

	return state, true

//// NB ignored        $state['text'] = $text = implode(' ', $this->getWordRange($state['startWordIndex'], $state['endWordIndex'] - 1));
////	state.text = strings.Join(parser.getWordRange(chart, state.startWordIndex, state.endWordIndex - 1), " ")
//	//
//	//        $Rule = $state['rule']->getSemantics();
//	parentSense := state.rule.Sense
//	//
//	//        if ($Rule) {
//	//
//	//            $childSemantics = $this->listChildSemantics($state);
////	childSenses := parser.listChildSenses(chart, state)
//	//            $childNodeTexts = $this->listChildTexts($state);
////		childNodeTexts := parser.listChildTexts(chart, state)
//	//
//	//            // combine the semantics of the children to determine the semantics of the parent
//	//            $Applier = new SemanticApplier();
//	//            $Semantics = $Applier->apply($Rule, $childSemantics, $childNodeTexts);
//
//// todo!!
////variableMap := map[string]string{}
//
//	//sense, _ := parser.senseBuilder.Join(parentSense, childSenses, variableMap)
//	//            $state['semantics'] = $Semantics;
//	state.sense = parentSense
//
//	if len(state.children) == 0 {
//		lexItem, _ := parser.lexicon.GetLexItem(state.rule.GetConsequent(0), state.rule.GetAntecedent())
//		state.sense = lexItem.RelationTemplates
//	}
//
//	//
//	//        return true;
//	return state, true
}

func (parser *parser) getWordRange(chart *chart, startWordIndex int, endWordIndex int) []string {
// NB! return array_slice($this->words, $startIndex, $endIndex - $startIndex + 1);
	return chart.words[startWordIndex:(endWordIndex + 1) + 1]
}

//// Returns an array of one sense (relation set) per consequent
//func (parser *parser) listChildSenses(chart *chart, state chartState ) []mentalese.RelationSet {
//
//	//        $childSemantics = array();
//	childSenses := []mentalese.RelationSet{}
//	//
//	//        /** @var ProductionRule $ProductionRule */
//	//        $ProductionRule = $state['rule']->getProduction();
//	//productionRule := state.rule
//	//
//	//        foreach ($state['children'] as $i => $childNodeId) {
//	for _, childNodeId := range state.children {
//
//	//
//	//            $childState = $this->treeInfo['states'][$childNodeId];
//		childState := chart.treeInfoStates[childNodeId]
//	//
//	// $childId = $ProductionRule->getConsequent($i);
//	// category := productionRule.GetConsequent(i)
//	//
//	// $childSemantics[$childId] = $childState['semantics'];
//		childSenses = append(childSenses, childState.sense)
//	}
//	//
//	//        return $childSemantics;
//	return childSenses
//}

func (parser *parser) getNextCat(state chartState) string {

	// return $state['rule']->getProduction()->getConsequentCategory($state['dotPosition'] - 1);
	return state.rule.GetConsequent(state.dotPosition - 1)
}

func (parser *parser) predict(chart *chart, state chartState) {

	common.LogTree("predict", state);

	// $nextConsequent = $state['rule']->getProduction()->getConsequentCategory($state['dotPosition'] - 1);
	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)

	// $endWordIndex = $state['endWordIndex'];
	endWordIndex := state.endWordIndex

	// go through all rules that have the next consequent as their antecedent
	// foreach ($this->Grammar->getParseRulesForAntecedent($nextConsequent) as $newRule) {
	for _, newRule := range parser.grammar.FindRules(nextConsequent) {

		//            $predictedState = array(
		//                'rule' => $newRule,
		//                'dotPosition' => 1,
		//                'startWordIndex' => $endWordIndex,
		//                'endWordIndex' => $endWordIndex,
		//                'semantics' => null,
		//            );
		predictedState := newChartState(newRule, 1, endWordIndex, endWordIndex)

		// $this->enqueue($predictedState, $endWordIndex);
		parser.enqueue(chart, predictedState, endWordIndex)
	}

	common.LogTree("predict");
}

func (parser *parser) scan(chart *chart, state chartState) {

	common.LogTree("scan", state);

	//	$nextConsequent = $state['rule']->getProduction()->getConsequentCategory($state['dotPosition'] - 1);
	nextConsequent := state.rule.GetConsequent(state.dotPosition - 1)

	// $endWordIndex = $state['endWordIndex'];
	endWordIndex := state.endWordIndex
	// $endWord = $this->words[$endWordIndex];
	endWord := chart.words[endWordIndex]

// NB if ($this->Grammar->isWordAPartOfSpeech($endWord, $nextConsequent)) {
	_, lexItemFound := parser.lexicon.GetLexItem(endWord, nextConsequent)
	if lexItemFound {

		// $Semantics = $this->Grammar->getSemanticsForWord($endWord, $nextConsequent);
//		sense := lexItem.RelationTemplates
//
//            if ($Semantics === false) {
//                throw new SemanticsNotFoundException($endWord);
//            }

//            $Production = new ProductionRule();
//            $Production->setAntecedent($nextConsequent);
//            $Production->setConsequents(array($endWord), false);
//            $NewRule = new ParseRule();
//            $NewRule->setProduction($Production);
		rule := parse.GrammarRule{ SyntacticCategories: []string{ nextConsequent, endWord }, EntityVariables: []string{"a", "b"}, Sense: mentalese.RelationSet{} }

//            $scannedState = array(
//                'rule' => $NewRule,
//                'dotPosition' => 2,
//                'startWordIndex' => $endWordIndex,
//                'endWordIndex' => $endWordIndex + 1,
//                'semantics' => $Semantics,
//            );

		scannedState := newChartState(rule, 2, endWordIndex, endWordIndex + 1)
//
//            $this->enqueue($scannedState, $endWordIndex + 1);
		parser.enqueue(chart, scannedState, endWordIndex + 1)
	}

	common.LogTree("scan", endWord, lexItemFound);
}

func (parser *parser) complete(chart *chart, completedState chartState) bool {

	common.LogTree("complete", completedState);

//	$treeComplete = false;
	treeComplete := false;
//
//        $completedAntecedent = $completedState['rule']->getProduction()->getAntecedentCategory();
	completedAntecedent := completedState.rule.GetAntecedent()
//
//        foreach ($this->chart[$completedState['startWordIndex']] as $chartedState) {
	for _, chartedState := range chart.states[completedState.startWordIndex] {
//
//            $dotPosition = $chartedState['dotPosition'];
		dotPosition := chartedState.dotPosition
//            $rule = $chartedState['rule'];
		rule := chartedState.rule
//
		// check if the antecedent of the completed state matches the charted state's consequent at the dot position
//            if (($dotPosition > $rule->getProduction()->getConsequentCount()) || ($rule->getProduction()->getConsequentCategory($dotPosition - 1) != $completedAntecedent)) {
		if (dotPosition > rule.GetConsequentCount()) || (rule.GetConsequent(dotPosition - 1) != completedAntecedent) {
//                continue;
			continue;
		}
//
//            $advancedState = array(
//                'rule' => $rule,
//                'dotPosition' => $dotPosition + 1,
//                'startWordIndex' => $chartedState['startWordIndex'],
//                'endWordIndex' => $completedState['endWordIndex'],
//                'semantics' => null
//            );
		advancedState := newChartState(rule, dotPosition + 1, chartedState.startWordIndex, completedState.endWordIndex)
//
//            // store extra information to make it easier to extract parse trees later
//            $treeComplete = $this->storeStateInfo($completedState, $chartedState, $advancedState);
		treeComplete, advancedState = parser.storeStateInfo(chart, completedState, chartedState, advancedState);
//
//            if ($treeComplete) {
//                break;
//            }
		if treeComplete {
			break;
		}
//
//            $this->enqueue($advancedState, $completedState['endWordIndex']);
		parser.enqueue(chart, advancedState, completedState.endWordIndex)

		common.LogTree("complete");
    }
//
//        return $treeComplete;
	return treeComplete
}

func (parser *parser) storeStateInfo(chart *chart, completedState chartState, chartedState chartState, advancedState chartState) (bool, chartState) {

	//        $treeComplete = false;
	treeComplete := false
	//
	//        // store the state's "children" to ease building the parse trees from the packed forest
	//        $advancedState['children'] = !isset($chartedState['children']) ? array() : $chartedState['children'];
	advancedState.children = chartedState.children
	//        $advancedState['children'][] = $completedState['id'];
	advancedState.children = append(advancedState.children, completedState.id)
	//
	//        // rule complete?
	//
	//        $consequentCount = $chartedState['rule']->getProduction()->getConsequentCount();
	consequentCount := chartedState.rule.GetConsequentCount()
	//
	//        if ($chartedState['dotPosition'] == $consequentCount) {
	if chartedState.dotPosition == consequentCount {
	//
	//            // complete sentence?
	//            $antecedent = $chartedState['rule']->getProduction()->getAntecedentCategory();
		antecedent := chartedState.rule.GetAntecedent()
	//
	//            if ($antecedent == 'gamma') {
		if antecedent == "gamma" {
	//
	//                // that matches all words?
	//                if ($completedState['endWordIndex'] == count($this->words)) {
			if completedState.endWordIndex == len(chart.words) {
	//
	//                    $this->treeInfo['sentences'][] = $advancedState;
				chart.treeInfoSentences = append(chart.treeInfoSentences, advancedState)
	//
	//                    // set a flag to allow the parser to stop at the first complete parse
	//                    $treeComplete = true;
				treeComplete = true
			}
		}
	}
	//
	//        return $treeComplete;
	return treeComplete, advancedState
}

func (parser *parser) extractFirstTree(chart *chart) ParseTreeNode {

	tree := ParseTreeNode{}

	//		if (!empty($this->treeInfo['sentences'])) {
	if len(chart.treeInfoSentences) > 0 {
	//			$root = $this->treeInfo['sentences'][0];
		root := chart.treeInfoSentences[0]
	//			$tree = $this->extractParseTreeBranch($root);
		tree = parser.extractParseTreeBranch(chart, root)

	} else {
	//			$tree = null;

	}
	//
	//		return $tree;


	return tree
}

func (parser *parser) extractParseTreeBranch(chart *chart, state chartState) ParseTreeNode {

	//		$rule = $state['rule'];
	rule := state.rule
	//
	//		$antecedent = $rule->getProduction()->getAntecedent();
	antecedent := rule.GetAntecedent()
	//		$antecedentCategory = $rule->getProduction()->getAntecedentCategory();
	//
	//		if ($antecedent == 'gamma') {
	if antecedent == "gamma" {
	//			$constituentId = $state['children'][0];
		constituentId := state.children[0]
	//			$constituent = $this->treeInfo['states'][$constituentId];
		constituent := chart.treeInfoStates[constituentId]
	//			return $this->extractParseTreeBranch($constituent);
		return parser.extractParseTreeBranch(chart, constituent)
	}
	//
	//		$branch = array(
	//			'part-of-speech' => $antecedentCategory
	//		);

	branch := ParseTreeNode{category: antecedent, constituents: []ParseTreeNode{}, form: ""}
	//
	//		if ($this->Grammar->isPartOfSpeech($antecedent)) {
	//			$branch['word'] = $rule->getProduction()->getConsequentCategory(0);
	//		}
// NB!
	if len(state.children) == 0 {
		branch.form = rule.GetConsequent(0)
	}
	//
	//		$branch['semantics'] = $state['semantics'];
	//
	//		if (isset($state['children'])) {
	if len(state.children) > 0 {
	//
	//			$constituents = array();
	//			foreach ($state['children'] as $constituentId) {
		for _, constituentId := range state.children {
	//				$constituent = $this->treeInfo['states'][$constituentId];
			constituent := chart.treeInfoStates[constituentId]
	//				$constituents[] = $this->extractParseTreeBranch($constituent);
			branch.constituents = append(branch.constituents, parser.extractParseTreeBranch(chart, constituent))
		}
	//			$branch['constituents'] = $constituents;
	}
	//
	//		return $branch;
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
			//parser.parseSingleConsequent(syntacticCategories[i], tokens, cursor, consequentVariable)
			relations = append(relations, childRelations...)

		}

	} else {

		lexItem, _ := parser.lexicon.GetLexItem(state.rule.GetConsequent(0), state.rule.GetAntecedent())
		//state.sense = lexItem.RelationTemplates

		ruleRelations := parser.senseBuilder.CreateLexItemRelations(lexItem.RelationTemplates, antecedentVariable)
		relations = append(relations, ruleRelations...)
	}

	common.LogTree("extractSenseFromState", relations)

	return relations
}

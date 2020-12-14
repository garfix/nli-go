package morphology

import "regexp"

type SegmentationRule struct {
	antecedent SegmentNode
	consequents []SegmentNode
	regexp *regexp.Regexp
}

func NewSegmentationRule(antecedent SegmentNode, consequents []SegmentNode, regexp *regexp.Regexp) SegmentationRule {
	return SegmentationRule{
		antecedent:  antecedent,
		consequents: consequents,
		regexp: regexp,
	}
}

func (rule SegmentationRule) GetAntecedent() SegmentNode {
	return rule.antecedent
}

func (rule SegmentationRule) GetConsequents() []SegmentNode {
	return rule.consequents
}

func (rule SegmentationRule) IsTerminal() bool {
	return len(rule.consequents) == 0
}

func (rule SegmentationRule) Matches(word string) ([]string, bool) {
	results := rule.regexp.FindStringSubmatch(word)

	return results, len(results) > 0
}

func BuildRegexp(pattern []SegmentPatternCharacter, characterClasses []CharacterClass) (*regexp.Regexp, bool) {

	exp := ""

	for _, character := range pattern {
		nodeExp, ok := buildExpressionForNode(character, characterClasses)
		if !ok {
			return nil, false
		}
		exp += nodeExp
	}

	compiled, err := regexp.Compile("^" + exp + "$")

	return compiled, err == nil
}

func buildExpressionForNode(character SegmentPatternCharacter, classes []CharacterClass) (string, bool) {
	switch character.characterType {
	case CharacterTypeRest:
		return "(.*)", true
	case CharacterTypeClass:
		name := character.characterValue
		class, found := characterClassByName(classes, name)
		if !found {
			return "", false
		}
		exp := ""
		sep := ""
		for _, char := range class.characters {
			exp += sep + char.TermValue
			sep = "|"
		}
		return "(" + exp + ")", true
	case CharacterTypeLiteral:
		return "(" + character.characterValue + ")", true
	}
	return "", false
}

func characterClassByName(characterClasses []CharacterClass, name string) (CharacterClass, bool) {
	for _, class := range characterClasses {
		if class.name == name {
			return class, true
		}
	}
	return CharacterClass{}, false
}
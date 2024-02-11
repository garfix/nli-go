package mentalese

import (
	"fmt"
	"strconv"
	"strings"
)

type Term struct {
	TermType             string      `json:"type"`
	TermValue            string      `json:"value,omitempty"`
	TermSort             string      `json:"sort,omitempty"`
	TermValueRelationSet RelationSet `json:"set,omitempty"`
	TermValueRule        *Rule       `json:"rule,omitempty"`
	TermValueList        TermList    `json:"list,omitempty"`
}

const TermTypeVariable = "variable"
const TermTypePredicateAtom = "atom"
const TermTypeStringConstant = "string"
const TermTypeAnonymousVariable = "anonymous"
const TermTypeRegExp = "regexp"
const TermTypeRelationSet = "relation-set"
const TermTypeFunctionCall = "function-call"
const TermTypeRule = "rule"
const TermTypeId = "id"
const TermTypeList = "list"

func NewTermVariable(name string) Term {
	return Term{TermType: TermTypeVariable, TermValue: name, TermValueRelationSet: nil}
}

func NewTermAnonymousVariable() Term {
	return Term{TermType: TermTypeAnonymousVariable, TermValue: "", TermValueRelationSet: nil}
}

func NewTermString(value string) Term {
	return Term{TermType: TermTypeStringConstant, TermValue: value, TermValueRelationSet: nil}
}

func NewTermAtom(value string) Term {
	return Term{TermType: TermTypePredicateAtom, TermValue: value, TermValueRelationSet: nil}
}

func NewTermRelationSet(value RelationSet) Term {
	return Term{TermType: TermTypeRelationSet, TermValue: "", TermValueRelationSet: value}
}

func NewTermRule(rule Rule) Term {
	return Term{TermType: TermTypeRule, TermValue: "", TermValueRelationSet: nil, TermValueRule: &rule}
}

func NewTermId(id string, sort string) Term {
	return Term{TermType: TermTypeId, TermValue: id, TermSort: sort, TermValueRelationSet: nil}
}

func NewTermList(list TermList) Term {
	return Term{TermType: TermTypeList, TermValueList: list}
}

func (term Term) IsVariable() bool {
	return term.TermType == TermTypeVariable
}

func (term Term) IsMutableVariable() bool {
	return term.TermType == TermTypeVariable && IsMutableVariable(term.TermValue)
}

func IsMutableVariable(variable string) bool {
	return variable[0:1] == ":"
}

func IsGeneratedVariable(variable string) bool {
	return strings.Contains(variable, "$")
}

func (term Term) IsInteger() bool {
	if term.TermType != TermTypeStringConstant {
		return false
	}
	_, err := strconv.Atoi(term.TermValue)
	return err == nil
}

func (term Term) GetIntValue() (int, bool) {
	if term.TermType != TermTypeStringConstant {
		return -1, false
	}
	value, err := strconv.Atoi(term.TermValue)
	return value, err == nil
}

func (term Term) IsNumber() bool {
	if term.TermType != TermTypeStringConstant {
		return false
	}
	_, err := strconv.ParseFloat(term.TermValue, 64)
	return err == nil
}

func (term Term) GetNumber() (float64, bool) {
	if term.TermType != TermTypeStringConstant {
		return -1, false
	}
	value, err := strconv.ParseFloat(term.TermValue, 64)
	return value, err == nil
}

func (term Term) IsString() bool {
	return term.TermType == TermTypeStringConstant
}

func (term Term) IsId() bool {
	return term.TermType == TermTypeId
}

func (term Term) IsRegExp() bool {
	return term.TermType == TermTypeRegExp
}

func (term Term) IsAnonymousVariable() bool {
	return term.TermType == TermTypeAnonymousVariable
}

func (term Term) IsAtom() bool {
	return term.TermType == TermTypePredicateAtom
}

func (term Term) IsRelationSet() bool {
	return term.TermType == TermTypeRelationSet
}

func (term Term) IsFunctionCall() bool {
	return term.TermType == TermTypeFunctionCall
}

func (term Term) IsRule() bool {
	return term.TermType == TermTypeRule
}

func (term Term) IsList() bool {
	return term.TermType == TermTypeList
}

func (term Term) ListContains(t Term) bool {
	contains := false
	for _, item := range term.TermValueList {
		if item.Equals(t) {
			contains = true
			break
		}
	}
	return contains
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermSort != otherTerm.TermSort {
		return false
	}
	switch term.TermType {
	case TermTypeRelationSet:
		return term.TermValueRelationSet.Equals(otherTerm.TermValueRelationSet)
	case TermTypeRule:
		return term.TermValueRule.Equals(*otherTerm.TermValueRule)
	case TermTypeList:
		return term.TermValueList.Equals(otherTerm.TermValueList)
	default:
		return term.TermValue == otherTerm.TermValue
	}
}

func (term Term) UsesVariable(variable string) bool {
	found := false
	if term.IsVariable() {
		found = found || term.TermValue == variable
	} else if term.IsRelationSet() {
		for _, rel := range term.TermValueRelationSet {
			found = found || rel.UsesVariable(variable)
		}
	} else if term.IsRule() {
		found = found || term.TermValueRule.Goal.UsesVariable(variable)
		for _, rel := range term.TermValueRule.Pattern {
			found = found || rel.UsesVariable(variable)
		}
	} else if term.IsList() {
		found = term.TermValueList.UsesVariable(variable)
	}
	return found
}

func (term Term) GetVariableNames() []string {
	names := []string{}

	if term.IsVariable() {
		names = append(names, term.TermValue)
	} else if term.IsRelationSet() {
		names = append(names, term.TermValueRelationSet.GetVariableNames()...)
	} else if term.IsRule() {
		names = append(names, term.TermValueRule.GetVariableNames()...)
	} else if term.IsList() {
		names = append(names, term.TermValueList.GetVariableNames()...)
	}

	return names
}

func (term Term) ConvertVariablesToConstants() Term {
	if term.IsVariable() {
		return NewTermAtom(strings.ToLower(term.TermValue))
	} else if term.IsRelationSet() {
		return NewTermRelationSet(term.TermValueRelationSet.ConvertVariablesToConstants())
	} else if term.IsRule() {
		return NewTermRule(term.TermValueRule.ConvertVariablesToConstants())
	} else if term.IsList() {
		return NewTermList(term.TermValueList.ConvertVariablesToConstants())
	}
	return term
}

func (term Term) ConvertVariablesToMutables() Term {
	if term.IsVariable() && !term.IsMutableVariable() {
		return NewTermVariable(":" + term.TermValue)
	} else if term.IsRelationSet() {
		return NewTermRelationSet(term.TermValueRelationSet.ConvertVariablesToMutables())
	} else if term.IsRule() {
		return NewTermRule(term.TermValueRule.ConvertToFunction())
	} else if term.IsList() {
		return NewTermList(term.TermValueList.ConvertVariablesToMutables())
	}
	return term
}

func (term Term) ConvertVariablesToImmutables() Term {
	if term.IsVariable() && term.IsMutableVariable() {
		return NewTermVariable(strings.Replace(term.TermValue, ":", "", 1))
	} else if term.IsRelationSet() {
		return NewTermRelationSet(term.TermValueRelationSet.ConvertVariablesToImmutables())
	} else if term.IsRule() {
		return NewTermRule(term.TermValueRule.ConvertVariablesToImmutables())
	} else if term.IsList() {
		return NewTermList(term.TermValueList.ConvertVariablesToImmutables())
	}
	return term
}

func (term Term) AsKey() string {
	return fmt.Sprintf("%s/%s/%s", term.TermType, term.TermValue, term.TermSort)
}

func (term Term) Copy() Term {
	newTerm := Term{}
	newTerm.TermType = term.TermType
	newTerm.TermValue = term.TermValue
	newTerm.TermSort = term.TermSort
	if term.IsRelationSet() {
		newTerm.TermValueRelationSet = term.TermValueRelationSet.Copy()
	} else if term.IsRule() {
		copy := term.TermValueRule.Copy()
		newTerm.TermValueRule = &copy
	} else if term.IsList() {
		newTerm.TermValueList = term.TermValueList.Copy()
	}
	return newTerm
}

func (term Term) Bind(binding Binding) Term {
	arg := term
	if term.IsVariable() {
		newValue, found := binding.Get(term.TermValue)
		if found {
			arg = newValue
		}
	} else if term.IsRelationSet() {
		arg.TermValueRelationSet = term.TermValueRelationSet.BindSingle(binding)
	} else if term.IsRule() {
		bound := term.TermValueRule.BindSingle(binding)
		arg.TermValueRule = &bound
	} else if term.IsList() {
		arg.TermValueList = term.TermValueList.Bind(binding)
	}
	return arg
}

func (term Term) ReplaceTerm(from Term, to Term) Term {

	relationArgument := term

	if term.IsRelationSet() {

		relationArgument.TermValueRelationSet = relationArgument.TermValueRelationSet.ReplaceTerm(from, to)

	} else if term.IsRule() {

		newGoals := RelationSet{relationArgument.TermValueRule.Goal}.ReplaceTerm(from, to)
		newPattern := relationArgument.TermValueRule.Pattern.ReplaceTerm(from, to)
		newRule := Rule{Goal: newGoals[0], Pattern: newPattern}
		relationArgument.TermValueRule = &newRule

	} else if term.IsList() {

		relationArgument.TermValueList = relationArgument.TermValueList.ReplaceTerm(from, to)

	} else {

		if term.Equals(from) {
			relationArgument = to.Copy()
		} else {
			relationArgument = term
		}
	}

	return relationArgument
}

// If term is a variable, and occurs in binding, returns its binding
// Otherwise, return term
func (term Term) Resolve(binding Binding) Term {

	resolved := term

	if term.IsVariable() {
		value, found := binding.Get(term.TermValue)
		if found {
			resolved = value
		}
	}

	return resolved
}

func (term Term) AsSimple() interface{} {
	switch term.TermType {
	case TermTypePredicateAtom:
		return term.TermValue
	case TermTypeStringConstant:
		number, ok := term.GetNumber()
		if !ok {
			return term.TermValue
		} else {
			return number
		}
	case TermTypeRelationSet:
		list := []interface{}{}
		for _, relation := range term.TermValueRelationSet {
			list = append(list, relation.AsSimple())
		}
		return list
	case TermTypeId:
		return term.TermSort + ":" + term.TermValue
	case TermTypeList:
		list := []interface{}{}
		for _, item := range term.TermValueList {
			list = append(list, item.AsSimple())
		}
		return list
	}
	return "unsupported: " + term.TermType
}

func (term Term) String() string {

	s := ""

	switch term.TermType {
	case TermTypeVariable:
		s = term.TermValue
	case TermTypePredicateAtom:
		s = term.TermValue
	case TermTypeStringConstant:
		_, err := strconv.Atoi(term.TermValue)
		if err == nil {
			s = term.TermValue
		} else {
			s = "'" + term.TermValue + "'"
		}
	case TermTypeRegExp:
		s = "~" + term.TermValue + "~"
	case TermTypeAnonymousVariable:
		s = "_"
	case TermTypeRelationSet:
		s = term.TermValueRelationSet.String()
	case TermTypeRule:
		s = term.TermValueRule.String()
	case TermTypeId:
		s = "`" + term.TermSort + ":" + term.TermValue + "`"
	case TermTypeList:
		s = term.TermValueList.String()
	default:
		s = "<unknown>"
	}
	return s
}

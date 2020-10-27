package mentalese

import (
	"fmt"
	"strconv"
	"strings"
)

type Term struct {
	TermType             string
	TermValue            string
	TermEntityType		 string
	TermValueRelationSet RelationSet
	TermValueRule        Rule
	TermValueList		 TermList
}

const TermTypeVariable = "variable"
const TermTypePredicateAtom = "atom"
const TermTypeStringConstant = "string"
const TermTypeAnonymousVariable = "anonymous"
const TermTypeRegExp = "regexp"
const TermTypeRelationSet = "relation-set"
const TermTypeRule = "rule"
const TermTypeId = "id"
const TermTypeList = "list"

func NewTermVariable(name string) Term {
	return Term{ TermType: TermTypeVariable, TermValue: name, TermValueRelationSet: nil}
}

func NewTermAnonymousVariable() Term {
	return Term{ TermType: TermTypeAnonymousVariable, TermValue: "", TermValueRelationSet: nil}
}

func NewTermString(value string) Term {
	return Term{ TermType: TermTypeStringConstant, TermValue: value, TermValueRelationSet: nil}
}

func NewTermAtom(value string) Term {
	return Term{ TermType: TermTypePredicateAtom, TermValue: value, TermValueRelationSet: nil}
}

func NewTermRelationSet(value RelationSet) Term {
	return Term{ TermType: TermTypeRelationSet, TermValue: "", TermValueRelationSet: value}
}

func NewTermRule(rule Rule) Term {
	return Term{ TermType: TermTypeRule, TermValue: "", TermValueRelationSet: nil, TermValueRule: rule}
}

func NewTermId(id string, entityType string) Term {
	return Term{ TermType: TermTypeId, TermValue: id, TermEntityType: entityType, TermValueRelationSet: nil}
}

func NewTermList(list TermList) Term {
	return Term{ TermType: TermTypeList, TermValueList: list }
}

func (term Term) IsVariable() bool {
	return term.TermType == TermTypeVariable
}

func (term Term) IsInteger() bool {
	if term.TermType != TermTypeStringConstant {
		return false
	}
	_, err := strconv.Atoi(term.TermValue)
	return err == nil
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

func (term Term) IsRule() bool {
	return term.TermType == TermTypeRule
}

func (term Term) IsList() bool {
	return term.TermType == TermTypeList
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermEntityType != otherTerm.TermEntityType {
		return false
	}
	switch term.TermType {
	case TermTypeRelationSet:
		return term.TermValueRelationSet.Equals(otherTerm.TermValueRelationSet)
	case TermTypeRule:
		return term.TermValueRule.Equals(otherTerm.TermValueRule)
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

func (term Term) AsKey() string {
	return fmt.Sprintf("%d/%s/%s", term.TermType, term.TermValue, term.TermEntityType)
}

func (term Term) Copy() Term {
	newTerm := Term{}
	newTerm.TermType = term.TermType
	newTerm.TermValue = term.TermValue
	newTerm.TermEntityType = term.TermEntityType
	if term.IsRelationSet() {
		newTerm.TermValueRelationSet = term.TermValueRelationSet.Copy()
	} else if term.IsRule() {
		newTerm.TermValueRule = term.TermValueRule.Copy()
	} else if term.IsList() {
		newTerm.TermValueList = term.TermValueList.Copy()
	}
	return newTerm
}

var x = 1000

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
		arg.TermValueRule = term.TermValueRule.BindSingle(binding)
	} else if term.IsList() {
		arg.TermValueList = term.TermValueList.Bind(binding)
	}
	return arg
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
		s = "/" + term.TermValue + "/"
	case TermTypeAnonymousVariable:
		s = "_"
	case TermTypeRelationSet:
		s = term.TermValueRelationSet.String()
	case TermTypeRule:
		s = term.TermValueRule.String()
	case TermTypeId:
		s = "`" + term.TermEntityType + ":" + term.TermValue + "`"
	case TermTypeList:
		s = term.TermValueList.String()
	default:
		s = "<unknown>"
	}
	return s
}

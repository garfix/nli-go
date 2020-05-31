package mentalese

import (
	"fmt"
	"strconv"
)

type Term struct {
	TermType             string
	TermValue            string
	TermEntityType		 string
	TermValueRelationSet RelationSet
	TermValueRule        Rule
}

const TermVariable = "variable"
const TermPredicateAtom = "atom"
const TermStringConstant = "string"
const TermAnonymousVariable = "anonymous"
const TermRegExp = "regexp"
const TermRelationSet = "relation-set"
const TermRule = "rule"
const TermId = "id"

func NewVariable(name string) Term {
	return Term{ TermType: TermVariable, TermValue: name, TermValueRelationSet: nil}
}

func NewAnonymousVariable() Term {
	return Term{ TermType: TermAnonymousVariable, TermValue: "", TermValueRelationSet: nil}
}

func NewString(value string) Term {
	return Term{ TermType: TermStringConstant, TermValue: value, TermValueRelationSet: nil}
}

func NewPredicateAtom(value string) Term {
	return Term{ TermType: TermPredicateAtom, TermValue: value, TermValueRelationSet: nil}
}

func NewRelationSet(value RelationSet) Term {
	return Term{ TermType: TermRelationSet, TermValue: "", TermValueRelationSet: value}
}

func NewRule(rule Rule) Term {
	return Term{ TermType: TermRule, TermValue: "", TermValueRelationSet: nil, TermValueRule: rule}
}

func NewId(id string, entityType string) Term {
	return Term{ TermType: TermId, TermValue: id, TermEntityType: entityType, TermValueRelationSet: nil}
}

func (term Term) IsVariable() bool {
	return term.TermType == TermVariable
}

func (term Term) IsNumber() bool {
	if term.TermType != TermStringConstant {
		return false
	}
	_, err := strconv.Atoi(term.TermValue)
	return err == nil
}

func (term Term) IsString() bool {
	return term.TermType == TermStringConstant
}

func (term Term) IsId() bool {
	return term.TermType == TermId
}

func (term Term) IsRegExp() bool {
	return term.TermType == TermRegExp
}

func (term Term) IsAnonymousVariable() bool {
	return term.TermType == TermAnonymousVariable
}

func (term Term) IsAtom() bool {
	return term.TermType == TermPredicateAtom
}

func (term Term) IsRelationSet() bool {
	return term.TermType == TermRelationSet
}

func (term Term) IsRule() bool {
	return term.TermType == TermRule
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermEntityType != otherTerm.TermEntityType {
		return false
	}
	if term.TermType == TermRelationSet {
		return term.TermValueRelationSet.Equals(otherTerm.TermValueRelationSet)
	}
	if term.TermType == TermRule {
		return term.TermValueRule.Equals(otherTerm.TermValueRule)
	}
	return term.TermValue == otherTerm.TermValue
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
	}
	return found
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
	}
	if term.IsRule() {
		newTerm.TermValueRule = term.TermValueRule.Copy()
	}
	return newTerm
}

func (term Term) Bind(binding Binding) Term {
	arg := term
	if term.IsVariable() {
		newValue, found := binding[term.TermValue]
		if found {
			arg = newValue
		}
	} else if term.IsRelationSet() {
		arg.TermValueRelationSet = term.TermValueRelationSet.BindSingle(binding)
	} else if term.IsRule() {
		arg.TermValueRule = term.TermValueRule.BindSingle(binding)
	}
	return arg
}

// If term is a variable, and occurs in binding, returns its binding
// Otherwise, return term
func (term Term) Resolve(binding Binding) Term {

	resolved := term

	if term.IsVariable() {
		 value, found := binding[term.TermValue]
		 if found {
		 	resolved = value
		 }
	}

	return resolved
}

func (term Term) String() string {

	s := ""

	switch term.TermType {
	case TermVariable:
		s = term.TermValue
	case TermPredicateAtom:
		s = term.TermValue
	case TermStringConstant:
		_, err := strconv.Atoi(term.TermValue)
		if err == nil {
			s = term.TermValue
		} else {
			s = "'" + term.TermValue + "'"
		}
	case TermRegExp:
		s = "/" + term.TermValue + "/"
	case TermAnonymousVariable:
		s = "_"
	case TermRelationSet:
		s = term.TermValueRelationSet.String()
	case TermRule:
		s = term.TermValueRule.String()
	case TermId:
		s = "`" + term.TermEntityType + ":" + term.TermValue + "`"
	default:
		s = "<unknown>"
	}
	return s
}

package mentalese

import "fmt"

type Term struct {
	TermType             int
	TermValue            string
	TermValueRelationSet RelationSet
}

const Term_variable = 1
const Term_predicateAtom = 2
const Term_stringConstant = 3
const Term_number = 4
const Term_anonymousVariable = 5
const Term_regExp = 6
const Term_relationSet = 7
const Term_id = 8

func NewVariable(name string) Term {
	return Term{ TermType: Term_variable, TermValue: name, TermValueRelationSet: nil}
}

func NewAnonymousVariable() Term {
	return Term{ TermType: Term_anonymousVariable, TermValue: nil, TermValueRelationSet: nil}
}

func NewNumber(number string) Term {
	return Term{ TermType: Term_number, TermValue: number, TermValueRelationSet: nil}
}

func NewString(value string) Term {
	return Term{ TermType: Term_stringConstant, TermValue: value, TermValueRelationSet: nil}
}

func NewPredicateAtom(value string) Term {
	return Term{ TermType: Term_predicateAtom, TermValue: value, TermValueRelationSet: nil}
}

func NewId(id string) Term {
	return Term{ TermType: Term_id, TermValue: id, TermValueRelationSet: nil}
}

func (term Term) IsVariable() bool {
	return term.TermType == Term_variable
}

func (term Term) IsNumber() bool {
	return term.TermType == Term_number
}

func (term Term) IsRegExp() bool {
	return term.TermType == Term_regExp
}

func (term Term) IsAnonymousVariable() bool {
	return term.TermType == Term_anonymousVariable
}

func (term Term) IsRelationSet() bool {
	return term.TermType == Term_relationSet
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermType == Term_relationSet {
		return term.TermValueRelationSet.Equals(otherTerm.TermValueRelationSet)
	} else {
		return term.TermValue == otherTerm.TermValue
	}
}

func (term Term) AsKey() string {
	return fmt.Sprintf("%d/%s", term.TermType, term.TermValue)
}

func (term Term) Copy() Term {
	newTerm := Term{}
	newTerm.TermType = term.TermType
	newTerm.TermValue = term.TermValue
	if term.IsRelationSet() {
		newTerm.TermValueRelationSet = term.TermValueRelationSet.Copy()
	}
	return newTerm
}

func (term Term) String() string {

	s := ""

	switch term.TermType {
	case Term_variable:
		s = term.TermValue
	case Term_predicateAtom:
		s = term.TermValue
	case Term_stringConstant:
		s = "'" + term.TermValue + "'"
	case Term_regExp:
		s = "/" + term.TermValue + "/"
	case Term_number:
		s = term.TermValue
	case Term_anonymousVariable:
		s = "_"
	case Term_relationSet:
		s = term.TermValueRelationSet.String()
	case Term_id:
		s = term.TermValue
	default:
		s = "<unknown>"
	}
	return s
}

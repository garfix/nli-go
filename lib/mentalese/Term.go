package mentalese

import "fmt"

type Term struct {
	TermType             int
	TermValue            string
	TermValueRelationSet RelationSet
}

const TermVariable = 1
const TermPredicateAtom = 2
const TermStringConstant = 3
const TermNumber = 4
const TermAnonymousVariable = 5
const TermRegExp = 6
const TermRelationSet = 7
const TermId = 8

func NewVariable(name string) Term {
	return Term{ TermType: TermVariable, TermValue: name, TermValueRelationSet: nil}
}

func NewAnonymousVariable() Term {
	return Term{ TermType: TermAnonymousVariable, TermValue: "", TermValueRelationSet: nil}
}

func NewNumber(number string) Term {
	return Term{ TermType: TermNumber, TermValue: number, TermValueRelationSet: nil}
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

func NewId(id string) Term {
	return Term{ TermType: TermId, TermValue: id, TermValueRelationSet: nil}
}

func (term Term) IsVariable() bool {
	return term.TermType == TermVariable
}

func (term Term) IsNumber() bool {
	return term.TermType == TermNumber
}

func (term Term) IsRegExp() bool {
	return term.TermType == TermRegExp
}

func (term Term) IsAnonymousVariable() bool {
	return term.TermType == TermAnonymousVariable
}

func (term Term) IsRelationSet() bool {
	return term.TermType == TermRelationSet
}

func (term Term) Equals(otherTerm Term) bool {
	if term.TermType != otherTerm.TermType {
		return false
	}
	if term.TermType == TermRelationSet {
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
		s = "'" + term.TermValue + "'"
	case TermRegExp:
		s = "/" + term.TermValue + "/"
	case TermNumber:
		s = term.TermValue
	case TermAnonymousVariable:
		s = "_"
	case TermRelationSet:
		s = term.TermValueRelationSet.String()
	case TermId:
		s = "`" + term.TermValue + "`"
	default:
		s = "<unknown>"
	}
	return s
}

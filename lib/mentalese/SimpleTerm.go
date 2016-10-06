package mentalese

import "fmt"

type SimpleTerm struct {
	TermType  int
	TermValue string
}

const Term_variable = 1
const Term_predicateAtom = 2
const Term_stringConstant = 3
const Term_number = 4
const Term_anonymousVariable = 5

func (term *SimpleTerm) IsVariable() bool {
	return term.TermType == Term_variable
}

func (term *SimpleTerm) IsAnonymousVariable() bool {
	return term.TermType == Term_anonymousVariable
}

func (term *SimpleTerm) Equals(otherTerm SimpleTerm) bool {
	return term.TermType == otherTerm.TermType && term.TermValue == otherTerm.TermValue
}

func (term *SimpleTerm) AsKey() string {
	return fmt.Sprintf("%d/%s", term.TermType, term.TermValue)
}

func (term *SimpleTerm) String() string {
	return fmt.Sprintf("%s", term.TermValue)
}

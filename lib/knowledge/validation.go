package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

func Validate(input mentalese.Relation, format string, log *common.SystemLog) bool {

	expectedLength := len(format)

	for i, c := range format {
		if c == '*' {
			return true
		}
		if i >= len(input.Arguments) {
			log.AddError("Function '" + input.Predicate + "' expects at least " + strconv.Itoa(expectedLength) + " arguments")
			return false
		}
		arg := input.Arguments[i]
		if c == 'a' && !arg.IsAtom() {
			return false
		}
		if c == 'v' && !arg.IsVariable() {
			log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be an unbound variable")
			return false
		}
		if c == 's' && !arg.IsString() {
			log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be a string")
			return false
		}
		if c == 'i' && !arg.IsNumber() {
			//			base.log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(i + 1) + " to be a number")
			return false
		}
		if c == 'S' {
			expectedLength = len(input.Arguments)
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsString() {
					log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(j + 1) + " to be a string")
					return false
				}
			}
			break
		}
		if c == 'V' {
			expectedLength = len(input.Arguments)
			for j := i; j < len(input.Arguments); j++ {
				arg = input.Arguments[j]
				if !arg.IsVariable() {
					log.AddError("Function '" + input.Predicate + "' expects argument " + strconv.Itoa(j + 1) + " to be an unbound variable")
					return false
				}
			}
			break
		}
	}

	if expectedLength != len(input.Arguments) {
		log.AddError("Function '" + input.Predicate + "' expects " + strconv.Itoa(expectedLength) + " arguments")
		return false
	}

	return true
}
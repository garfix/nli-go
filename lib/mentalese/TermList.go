package mentalese

import (
	"sort"
	"strconv"
)

type TermList []Term

func (list TermList) Equals(otherList TermList) bool {
	if len(list) != len(otherList) {
		return false
	}
	for i, child := range list {
		if !child.Equals(otherList[i]) {
			return false
		}
	}
	return true
}

func (list TermList) Append(term Term) TermList {
	newList := list.Copy()
	newList = append(newList, term)
	return newList
}

func (list TermList) Set(index int, term Term) TermList {
	newList := list.Copy()
	newList[index] = term
	return newList
}

func (list TermList) UsesVariable(variable string) bool {
	for _, element := range list {
		if element.UsesVariable(variable) {
			return true
		}
	}
	return false
}

func (list TermList) ConvertVariablesToConstants() TermList {
	newList := TermList{}
	for _, element := range list {
		newList = append(newList, element.ConvertVariablesToConstants())
	}
	return newList
}

func (list TermList) Deduplicate() TermList {
	newList := TermList{}
	for _, element := range list {
		found := false
		for _, e := range newList {
			if element.Equals(e) {
				found = true
				break
			}
		}
		if !found {
			newList = append(newList, element)
		}
	}
	return newList
}

func (list TermList) GetTermType() (string, bool) {
	aType := ""

	for _, element := range list {
		if aType == "" {
			aType = element.TermType
		} else if aType != element.TermType {
			return "", false
		}
	}

	return aType, true
}

func (list TermList) IsIntegerList() bool {
	for _, element := range list {
		if !element.IsInteger() {
			return false
		}
	}

	return true
}

func (list TermList) GetValues() []string {
	values := []string{}
	for _, e := range list {
		values = append(values, e.TermValue)
	}
	return values
}

func (list TermList) ReplaceTerm(from Term, to Term) TermList {
	newList := TermList{}
	for _, element := range list {
		newList = append(newList, element.ReplaceTerm(from, to))
	}
	return newList
}

func (list TermList) Sort() (TermList, bool) {
	newList := TermList{}

	termType, ok := list.GetTermType()
	if !ok {
		return TermList{}, false
	}
	if termType == TermTypeStringConstant {
		if list.IsIntegerList() {
			stringValues := list.GetValues()
			values := []int{}
			for _, stringValue := range stringValues {
				i, _ := strconv.Atoi(stringValue)
				values = append(values, i)
			}
			sort.Ints(values)

			for _, value := range values {
				newList = append(newList, NewTermString(strconv.Itoa(value)))
			}
		} else {
			values := list.GetValues()
			sort.Strings(values)

			for _, value := range values {
				newList = append(newList, NewTermString(value))
			}
		}
	}
	return newList, true
}

func (list TermList) GetVariableNames() []string {
	names := []string{}
	for _, element := range list {
		names = append(names, element.GetVariableNames()...)
	}
	return names
}

func (list TermList) Copy() TermList {
	newList := TermList{}
	for _, element := range list {
		newList = append(newList, element.Copy())
	}
	return newList
}

func (list TermList) Bind(binding Binding) TermList {
	newList := TermList{}
	for _, element := range list {
		newList = append(newList, element.Bind(binding))
	}
	return newList
}

func (list TermList) String() string {
	s := ""
	sep := ""
	for _, element := range list {
		s += sep + element.String()
		sep = ", "
	}
	return "[" + s + "]"
}

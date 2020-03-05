package parse

type VariableList []string

func (list VariableList) Push(newList VariableList) VariableList {
	return append(list, newList...)
}

func (list VariableList) Pop(number int) (VariableList, bool) {
	ok := true
	newList := VariableList{}

	if len(list) < number {
		ok = false
	} else {
		newList = list[0:len(list) - number]
	}

	return newList, ok
}

func (list VariableList) Length() int {
	return len(list)
}

func (list VariableList) Empty() bool {
	return len(list) == 0
}

func (list VariableList) Copy() VariableList {
	newList := VariableList{}
	for _, variable := range list {
		newList = append(newList, variable)
	}
	return newList
}

func (list VariableList) Equals(otherList VariableList) bool {
	if len(list) != len(otherList) {
		return false
	}

	for i := range list {
		if list[i] != otherList[i] {
			return false
		}
	}

	return true
}

func (list VariableList) String() string {
	s := ""
	sep := ""

	for _, variable := range list {
		s += sep + variable
		sep = ", "
	}

	return "[" + s + "]"
}
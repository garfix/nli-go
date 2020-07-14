package mentalese

type TermList []Term

func (list TermList) Equals(otherList TermList) bool {
	if len(list) != len(otherList) { return false }
	for i, child := range list {
		if !child.Equals(otherList[i]) { return false }
	}
	return true
}

func (list TermList) UsesVariable(variable string) bool {
	for _, element := range list {
		if element.UsesVariable(variable) { return true }
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
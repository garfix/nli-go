package mentalese

type BindingSet struct {
	bindings *[]Binding
}

func NewBindingSet() BindingSet{
	return BindingSet{ bindings: &[]Binding{} }
}

func InitBindingSet(binding Binding) BindingSet{
	return BindingSet{ bindings: &[]Binding{ binding } }
}

func (set BindingSet) Add(binding Binding) {
	for _, b := range *set.bindings {
		if b.Equals(binding) {
			return
		}
	}
	*set.bindings = append(*set.bindings, binding)
}

func (set BindingSet) AddMultiple(bindingSet BindingSet) {
	for _, binding := range *bindingSet.bindings {
		set.Add(binding)
	}
}

func (set BindingSet) Copy() BindingSet {
	newSet := BindingSet{ bindings: &[]Binding{} }
	for _, binding := range *set.bindings {
		*newSet.bindings = append(*newSet.bindings, binding.Copy())
	}
	return newSet
}

func (set BindingSet) String() string {
	str := ""
	sep := ""

	for _, binding := range *set.bindings {
		str += sep + binding.String()
		sep = " "
	}

	return "[" + str + "]"
}

func (set BindingSet) GetAll() []Binding {
	return *set.bindings
}

func (set BindingSet) Get(index int) Binding {
	return (*set.bindings)[index]
}

func (set BindingSet) GetLength() int {
	return len(*set.bindings)
}

func (set BindingSet) IsEmpty() bool {
	return len(*set.bindings) == 0
}

func (set BindingSet) GetIds(variable string) []Term {
	idMap := map[string]bool{}
	ids := []Term{}

	for _, binding := range *set.bindings {
		for key, value := range binding.GetAll() {
			if key != variable {
				continue
			}
			if value.IsId() {
				found := idMap[value.String()]
				if !found {
					ids = append(ids, value)
				}
			}
		}
	}

	return ids
}

func (set BindingSet) GetDistinctValueCount(variable string) int {
	idMap := map[string]bool{}
	count := 0

	for _, binding := range *set.bindings {
		for key, value := range binding.GetAll() {
			if key != variable {
				continue
			}
			found := idMap[value.String()]
			if !found {
				count++
				idMap[value.String()] = true
			}
		}
	}

	return count
}

func (set BindingSet) GetDistinctValues(variable string) []Term {
	idMap := map[string]bool{}
	values := []Term{}

	for _, binding := range *set.bindings {
		for key, value := range binding.GetAll() {
			if key != variable {
				continue
			}
			found := idMap[value.String()]
			if !found {
				values = append(values, value)
				idMap[value.String()] = true
			}
		}
	}

	return values
}

func (s BindingSet) FilterVariablesByName(variableNames []string) BindingSet {
	newBindings := NewBindingSet()

	for _, binding := range *s.bindings {
		newBindings.Add(binding.FilterVariablesByName(variableNames))
	}

	return newBindings
}


func (s BindingSet) FilterOutVariablesByName(variableNames []string) BindingSet {
	newBindings := NewBindingSet()

	for _, binding := range *s.bindings {
		newBindings.Add(binding.FilterOutVariablesByName(variableNames))
	}

	return newBindings
}

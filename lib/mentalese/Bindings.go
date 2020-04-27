package mentalese

type Bindings []Binding

func (bindings Bindings) String() string {
	str := ""
	sep := ""

	for _, binding := range bindings {
		str += sep + binding.String()
		sep = " "
	}

	return "[" + str + "]"
}

func (bindings Bindings) GetIds(variable string) []Term {
	idMap := map[string]bool{}
	ids := []Term{}

	for _, binding := range bindings {
		for key, value := range binding {
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

func (bindings Bindings) GetDistinctValueCount(variable string) int {
	idMap := map[string]bool{}
	count := 0

	for _, binding := range bindings {
		for key, value := range binding {
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

func (bindings Bindings) GetDistinctValues(variable string) []string {
	idMap := map[string]bool{}
	values := []string{}

	for _, binding := range bindings {
		for key, value := range binding {
			if key != variable {
				continue
			}
			found := idMap[value.String()]
			if !found {
				values = append(values, value.TermValue)
				idMap[value.String()] = true
			}
		}
	}

	return values
}

func (s Bindings) FilterVariablesByName(variableNames []string) Bindings {
	newBindings := []Binding{}

	for _, binding := range s {
		newBindings = append(newBindings, binding.FilterVariablesByName(variableNames))
	}

	return newBindings
}

// Returns copy of bindings that contains each Binding only once
func (s Bindings) UniqueBindings() Bindings {
	uniqueBindings := Bindings{}
	for _, binding := range s {
		found := false
		for _, uniqueBinding := range uniqueBindings {
			if uniqueBinding.Equals(binding) {
				found = true
				break
			}
		}
		if !found {
			uniqueBindings = append(uniqueBindings, binding)
		}
	}
	return uniqueBindings
}

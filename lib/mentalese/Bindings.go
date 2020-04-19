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
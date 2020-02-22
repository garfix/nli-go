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

func (bindings Bindings) GetIds() []Term {
	idMap := map[string]bool{}
	ids := []Term{}

	for _, binding := range bindings {
		for _, value := range binding {
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
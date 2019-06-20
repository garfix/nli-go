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
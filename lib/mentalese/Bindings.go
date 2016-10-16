package mentalese

type Binding map[string]Term

func (b Binding) Merge(b2 Binding) Binding {

	result := Binding{}

	for k, v := range b {
		result[k] = v
	}

	for k, v := range b2 {
		result[k] = v
	}

	return result
}

// Binds the variables of b to the values of c
// example:
// b: A = E
//    B = 3
// c: E = 5
//    F = 6
// result:
//    A = 5
//    B = 3
// note: F is discarded
func (b Binding) Bind(c Binding) Binding {

	result := Binding{}.Merge(b)

	for bKey, bVal := range b {

		result[bKey] = bVal

		if bVal.IsVariable() {
			for cKey, cVal := range c {
				if  bVal.TermValue == cKey {
					result[bKey] = cVal
				}
			}
		}
	}

	return result
}

func (b Binding) String() string {

	s, sep := "", ""

	for k, v := range b {
		s += sep + k + ":" + v.String()
		sep = ", "
	}

	return "{" + s + "}"
}


func (b Binding) Equals(c Binding) bool {

	if len(b) != len(c) {
		return false
	}

	for key, value := range b {
		if c[key] != value {
			return false
		}
	}

	return true
}
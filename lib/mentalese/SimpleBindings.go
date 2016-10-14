package mentalese

type SimpleBinding map[string]SimpleTerm

func (b SimpleBinding) Merge(b2 SimpleBinding) SimpleBinding {

	result := SimpleBinding{}

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
func (b SimpleBinding) Bind(c SimpleBinding) SimpleBinding {

	result := SimpleBinding{}.Merge(b)

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

func (b SimpleBinding) String() string {

	s, sep := "", ""

	for k, v := range b {
		s += sep + k + "=" + v.String()
		sep = ", "
	}

	return "{" + s + "}"
}
package mentalese

import "sort"

type Binding map[string]Term

// Returns a new binding that is a copy of b, merged with b2
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

// Returns a new binding that is contains just the keys of b, and whose values may be overwritten by those of b2
func (b Binding) Intersection(b2 Binding) Binding {

	result := Binding{}

	for k, v := range b {
		result[k] = v
	}

	for k, v := range b2 {
		_, found := result[k]
		if found {
			result[k] = v
		}
	}

	return result
}

// Returns a copy
func (b Binding) Copy() Binding {

	result := Binding{}

	for k, v := range b {
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

// Returns a new binding with just key, if exists
func (b Binding) Extract(key string) Binding {
	newBinding := Binding{}

	val, found := b[key]
	if found {
		newBinding[key] = val
	}

	return newBinding
}

// Returns a string version
func (b Binding) String() string {

	s, sep := "", ""
	keys := []string{}

	for k := range b {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += sep + k + ":" + b[k].String()
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
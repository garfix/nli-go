package mentalese

import (
	"nli-go/lib/common"
	"sort"
)

type Binding struct {
	Key2vvalue map[string]Term
}

func NewBinding() Binding {
	return Binding{Key2vvalue: map[string]Term{}}
}

func (b Binding) ToRaw() map[string]Term {
	return b.Key2vvalue
}

func (p Binding) FromRaw(raw map[string]Term) {
	for key, value := range raw {
		p.Set(key, value)
	}
}

func (b *Binding) Clear() {
	b.Key2vvalue = map[string]Term{}
}

func (b Binding) ContainsVariable(variable string) bool {
	_, found := b.Key2vvalue[variable]
	return found
}

func (b Binding) Set(variable string, value Term) {
	b.Key2vvalue[variable] = value
}

func (b Binding) Get(variable string) (Term, bool) {
	value, found := b.Key2vvalue[variable]
	return value, found
}

func (b Binding) MustGet(variable string) Term {
	value, found := b.Get(variable)
	if found {
		return value
	} else {
		panic("variable not found: " + variable)
	}
}

func (b Binding) GetAll() map[string]Term {
	all := map[string]Term{}
	for key, value := range b.Key2vvalue {
		all[key] = value
	}
	return all
}

func (b Binding) GetKeys() []string {
	all := []string{}
	for key := range b.Key2vvalue {
		all = append(all, key)
	}
	return all
}

// Returns a new Binding that is a copy of b, merged with b2
func (b Binding) Merge(b2 Binding) Binding {

	result := NewBinding()

	for k, v := range b.Key2vvalue {
		result.Key2vvalue[k] = v
	}

	for k, v := range b2.Key2vvalue {
		result.Key2vvalue[k] = v
	}

	return result
}

// Returns a new Binding that contains just the keys of b, and whose values may be overwritten by those of b2
func (b Binding) Intersection(b2 Binding) Binding {

	result := NewBinding()

	for k, v := range b.Key2vvalue {
		result.Key2vvalue[k] = v
	}

	for k, v := range b2.Key2vvalue {
		_, found := result.Key2vvalue[k]
		if found {
			result.Key2vvalue[k] = v
		}
	}

	return result
}

// returns a binding with only given keys, if present
func (b Binding) Select(keys []string) Binding {
	newBinding := NewBinding()

	for _, key := range keys {
		value, found := b.Key2vvalue[key]
		if found {
			newBinding.Key2vvalue[key] = value
		}
	}

	return newBinding
}

// Returns a copy
func (b Binding) Copy() Binding {

	result := NewBinding()

	for k, v := range b.Key2vvalue {
		result.Key2vvalue[k] = v
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

	result := NewBinding().Merge(b)

	for bKey, bVal := range b.Key2vvalue {

		result.Key2vvalue[bKey] = bVal

		if bVal.IsVariable() {
			value, found := c.Key2vvalue[bVal.TermValue]
			if found {
				result.Key2vvalue[bKey] = value
			}
		}
	}

	return result
}

// Returns a version of b without the keys that have variable values
func (b Binding) RemoveVariables() Binding {

	result := NewBinding()

	for key, value := range b.Key2vvalue {
		if !value.IsVariable() {
			result.Key2vvalue[key] = value
		}
	}

	return result
}

// Returns a version of b with key and value swapped. Only variable values survive
// In:
// { A:11, B: X }
// Out:
// { X: B }
func (b Binding) Swap() Binding {

	result := NewBinding()

	for key, value := range b.Key2vvalue {
		if value.IsVariable() {
			result.Key2vvalue[value.TermValue] = Term{TermType: TermTypeVariable, TermValue: key}
		}
	}

	return result
}

func (b Binding) FilterVariablesByName(variableNames []string) Binding {
	result := NewBinding()

	for _, variableName := range variableNames {
		_, found := b.Key2vvalue[variableName]
		if found {
			result.Key2vvalue[variableName] = b.Key2vvalue[variableName]
		}
	}

	return result
}

func (b Binding) FilterOutVariablesByName(variableNames []string) Binding {
	result := NewBinding()

	for key, value := range b.Key2vvalue {
		if !common.StringArrayContains(variableNames, key) {
			result.Key2vvalue[key] = value
		}
	}

	return result
}

// Returns a new Binding with just key, if exists
func (b Binding) Extract(key string) Binding {
	newBinding := NewBinding()

	val, found := b.Key2vvalue[key]
	if found {
		newBinding.Key2vvalue[key] = val
	}

	return newBinding
}

// Returns a string version
func (b Binding) String() string {

	s, sep := "", ""
	keys := []string{}

	for k := range b.Key2vvalue {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += sep + k + ":" + b.Key2vvalue[k].String()
		sep = ", "
	}

	return "{" + s + "}"
}

func (b Binding) Equals(c Binding) bool {

	if len(b.Key2vvalue) != len(c.Key2vvalue) {
		return false
	}

	for key, bValue := range b.Key2vvalue {
		cValue, found := c.Key2vvalue[key]
		if !found {
			return false
		}
		if !cValue.Equals(bValue) {
			return false
		}
	}

	return true
}

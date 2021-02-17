package mentalese

import (
	"nli-go/lib/common"
	"sort"
)

type Binding struct {
	k2v map[string]Term
	scope *Scope
}

func NewBinding() Binding {
	return Binding{ k2v: map[string]Term{}, scope: nil }
}

func NewScopedBinding(scope *Scope) Binding {
	return Binding{ k2v: map[string]Term{}, scope: scope }
}

func (b Binding) ToRaw() map[string]Term {
	return b.k2v
}

func (p Binding) FromRaw(raw map[string]Term) {
	for key, value := range raw {
		p.Set(key, value)
	}
}

func (b Binding) GetScope() *Scope {
	return b.scope
}

func (b Binding) ContainsVariable(variable string) bool {
	_, found := b.k2v[variable]
	return found
}

func (b Binding) Set(variable string, value Term) {
	if b.scope != nil {
		variables := b.scope.GetVariables()
		if variables.ContainsVariable(variable) {
			variables.Set(variable, value)
			return
		}
	}
	b.k2v[variable] = value
}

func (b Binding) Get(variable string) (Term, bool) {

	if b.scope != nil {
		variables := b.scope.GetVariables()
		value, found := variables.Get(variable)
		if found {
			return value, true
		}
	}

	value, found := b.k2v[variable]
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
	for key, value := range b.k2v {
		all[key] = value
	}
	if b.scope != nil {
		for key, value := range b.scope.GetVariables().k2v {
			all[key] = value
		}
	}
	return all
}

// Returns a new Binding that is a copy of b, merged with b2
func (b Binding) Merge(b2 Binding) Binding {

	result := NewScopedBinding(b.scope)

	for k, v := range b.k2v {
		result.k2v[k] = v
	}

	for k, v := range b2.k2v {
		result.k2v[k] = v
	}

	if b.scope != nil && b2.scope != nil && b.scope != b2.scope {
		b.scope.variables = b.scope.GetVariables().Merge(*b2.scope.GetVariables())
	}

	return result
}

// Returns a new Binding that contains just the keys of b, and whose values may be overwritten by those of b2
func (b Binding) Intersection(b2 Binding) Binding {

	result := NewScopedBinding(b.scope)

	if b2.scope != nil {
		panic("binding has scope")
	}

	for k, v := range b.k2v {
		result.k2v[k] = v
	}

	for k, v := range b2.k2v {
		_, found := result.k2v[k]
		if found {
			result.k2v[k] = v
		}
	}

	return result
}

// returns a binding with only given keys, if present
func (b Binding) Select(keys []string) Binding {
	newBinding := NewBinding()

	for _, key := range keys {
		value, found := b.k2v[key]
		if found {
			newBinding.k2v[key] = value
		}
	}

	return newBinding
}

// Returns a copy
func (b Binding) Copy() Binding {

	result := NewScopedBinding(b.scope)

	for k, v := range b.k2v {
		result.k2v[k] = v
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

	result := NewScopedBinding(b.scope).Merge(b)

	if c.scope != nil {
		panic("binding has scope")
	}

	for bKey, bVal := range b.k2v {

		result.k2v[bKey] = bVal

		if bVal.IsVariable() {
			value, found := c.k2v[bVal.TermValue]
			if found {
				result.k2v[bKey] = value
			}
		}
	}

	return result
}

// Returns a version of b without the keys that have variable values
func (b Binding) RemoveVariables() Binding {

	result := NewScopedBinding(b.scope)

	for key, value := range b.k2v {
		if !value.IsVariable() {
			result.k2v[key] = value
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

	result := NewScopedBinding(b.scope)

	for key, value := range b.k2v {
		if value.IsVariable() {
			result.k2v[value.TermValue] = Term{TermType: TermTypeVariable, TermValue: key}
		}
	}

	return result
}

func (b Binding) FilterVariablesByName(variableNames []string) Binding {
	result := NewScopedBinding(b.scope)

	for _, variableName := range variableNames {
		_, found := b.k2v[variableName]
		if found {
			result.k2v[variableName] = b.k2v[variableName]
		}
	}

	return result
}

func (b Binding) FilterOutVariablesByName(variableNames []string) Binding {
	result := NewScopedBinding(b.scope)

	for key, value := range b.k2v {
		if !common.StringArrayContains(variableNames, key) {
			result.k2v[key] = value
		}
	}

	return result
}

// Returns a new Binding with just key, if exists
func (b Binding) Extract(key string) Binding {
	newBinding := NewBinding()

	val, found := b.k2v[key]
	if found {
		newBinding.k2v[key] = val
	}

	return newBinding
}

// Returns a string version
func (b Binding) String() string {

	s, sep := "", ""
	keys := []string{}

	for k := range b.k2v {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		s += sep + k + ":" + b.k2v[k].String()
		sep = ", "
	}

	local := ""
	if b.scope != nil {
		local = "&" + b.scope.variables.String()
	}

	return "{" + s + "}" + local
}

func (b Binding) Equals(c Binding) bool {

	if len(b.k2v) != len(c.k2v) {
		return false
	}

	for key, bValue := range b.k2v {
		cValue, found := c.k2v[key]
		if !found {
			return false
		}
		if !cValue.Equals(bValue) {
			return false
		}
	}

	return true
}

package mentalese

import (
	"encoding/json"
	"sort"
	"strconv"
)

type BindingSet struct {
	Bindings *[]Binding
	Lookup   *map[string]bool
}

func NewBindingSet() BindingSet{
	return BindingSet{ Bindings: &[]Binding{}, Lookup: &map[string]bool{} }
}

func InitBindingSet(binding Binding) BindingSet{
	return BindingSet{ Bindings: &[]Binding{binding }, Lookup: &map[string]bool{} }
}

func (set BindingSet) Add(binding Binding) {
	serialized := binding.String()
	_, found := (*set.Lookup)[serialized]
	if found { return }
	(*set.Lookup)[serialized] = true
	*set.Bindings = append(*set.Bindings, binding)
}

func (set BindingSet) AddMultiple(bindingSet BindingSet) {
	for _, binding := range *bindingSet.Bindings {
		set.Add(binding)
	}
}

func (set BindingSet) ToRaw() []map[string]Term {
	raw := []map[string]Term{}
	for _, b := range *set.Bindings {
		raw = append(raw, b.Key2vvalue)
	}
	return raw
}

func (set BindingSet) FromRaw(raw []map[string]Term) {
	for _, b := range raw {
		binding := NewBinding()
		for key, value := range b {
			binding.Set(key, value)
		}
		set.Add(binding)
	}
}

func (set BindingSet) ToJson() string {

	type aMap = map[string]string
	type array = []aMap

	arr := array{}

	for _, item := range set.GetAll() {
	i := aMap{}
	for k, v := range item.GetAll() {
		i[k] = v.String()
	}
		arr = append(arr, i)
	}

	responseRaw, _ := json.MarshalIndent(arr, "", "    ")

	return string(responseRaw)
}

func (set BindingSet) Reverse() BindingSet {
	newSet := NewBindingSet()
	lastIndex := len(*set.Bindings) - 1
	for i, _ := range *set.Bindings {
		binding := (*set.Bindings)[lastIndex - i]
		newSet.Add(binding)
	}
	return newSet
}

func (set BindingSet) GetTermType(variable string) (string, bool) {
	aType := ""

	for _, binding := range *set.Bindings {
		term, found := binding.Key2vvalue[variable]
		if !found { continue }

		if aType == "" {
			aType = term.TermType
		} else if aType != term.TermType {
			return "", false
		}
	}

	return aType, true
}


func (set BindingSet) IsIntegerSet(variable string) bool {
	for _, binding := range *set.Bindings {
		term, found := binding.Key2vvalue[variable]
		if !found { continue }

		if !term.IsInteger() {
			return false
		}
	}

	return true
}

func (set BindingSet) Sort(variable string) (BindingSet, bool) {
	newSet := NewBindingSet()

	// collect all integer terms
	// group by value
	numbers := map[float64][]Binding{}
	strings := map[string][]Binding{}

	for _, binding := range set.GetAll() {
		term := binding.Key2vvalue[variable]
		if term.IsNumber() {
			number, _ := strconv.ParseFloat(term.TermValue, 64)
			_, found := numbers[number]
			if !found {
				numbers[number] = []Binding{}
			}
			numbers[number] = append(numbers[number], binding)
		} else if term.IsString() {
			_, found := strings[term.TermValue]
			if !found {
				strings[term.TermValue] = []Binding{}
			}
			strings[term.TermValue] = append(strings[term.TermValue], binding)
		} else {
			return newSet, false
		}
	}

	if len(numbers) > 0 && len(strings) == 0 {

		sortedNumbers := []float64{}
		for integer, _ := range numbers {
			sortedNumbers = append(sortedNumbers, integer)
		}
		sort.Float64s(sortedNumbers)

		for _, integer := range sortedNumbers {
			for _, binding := range numbers[integer] {
				newSet.Add(binding)
			}
		}

	} else if len(numbers) == 0 {

		sortedStrings := []string{}
		for str, _ := range strings {
			sortedStrings = append(sortedStrings, str)
		}
		sort.Strings(sortedStrings)

		for _, str := range sortedStrings {
			for _, binding := range strings[str] {
				newSet.Add(binding)
			}
		}

	} else {
		return newSet, false
	}

	return newSet, true
}

func (set BindingSet) Copy() BindingSet {
	newSet := NewBindingSet()
	for _, binding := range *set.Bindings {
		newSet.Add(binding)
	}
	return newSet
}

func (set BindingSet) String() string {
	str := ""
	sep := ""

	for _, binding := range *set.Bindings {
		str += sep + binding.String()
		sep = " "
	}

	return "[" + str + "]"
}

func (set BindingSet) GetAll() []Binding {
	return *set.Bindings
}

func (set BindingSet) Get(index int) Binding {
	return (*set.Bindings)[index]
}

func (set BindingSet) GetLength() int {
	return len(*set.Bindings)
}

func (set BindingSet) IsEmpty() bool {
	return len(*set.Bindings) == 0
}

func (set BindingSet) GetIds(variable string) []Term {
	idMap := map[string]bool{}
	ids := []Term{}

	for _, binding := range *set.Bindings {
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

	for _, binding := range *set.Bindings {
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

func (set BindingSet) GetAllVariableValues(variable string) []Term {
	values := []Term{}

	for _, binding := range set.GetAll() {
		term, found := binding.Key2vvalue[variable]
		if !found {
			continue
		}
		values = append(values, term)
	}

	return values
}

func (set BindingSet) GetDistinctValues(variable string) []Term {
	idMap := map[string]bool{}
	values := []Term{}

	for _, binding := range *set.Bindings {
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

	for _, binding := range *s.Bindings {
		newBindings.Add(binding.FilterVariablesByName(variableNames))
	}

	return newBindings
}

func (s BindingSet) FilterOutVariablesByName(variableNames []string) BindingSet {
	newBindings := NewBindingSet()

	for _, binding := range *s.Bindings {
		newBindings.Add(binding.FilterOutVariablesByName(variableNames))
	}

	return newBindings
}

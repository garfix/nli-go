package mentalese

import "strings"

type EntitySorts map[string][]string

func NewEntitySorts() *EntitySorts {
	return &EntitySorts{}
}

func (s *EntitySorts) Clear() {
	*s = map[string][]string{}
}

func (p *EntitySorts) Copy() *EntitySorts {

	newSorts := EntitySorts{}
	for k, v := range *p {
		newSorts[k] = v
	}

	return &newSorts
}

func (s *EntitySorts) SetSorts(variable string, sorts []string) {
	(*s)[variable] = sorts
}

// an entity usually has a single value, but it can also contain a list of values
// this function always returns a list
func (s *EntitySorts) GetSorts(variable string) []string {
	sort, found := (*s)[variable]
	if found {
		return sort
	} else {
		return []string{}
	}
}

func (s *EntitySorts) ReplaceVariable(fromVariable string, toVariable string) {
	sort, found := (*s)[fromVariable]
	if found {
		delete(*s, fromVariable)
		(*s)[toVariable] = sort
	}
}

func (s *EntitySorts) String() string {
	str := ""

	for key, value := range *s {
		str += key + "=" + strings.Join(value, ",") + ";"
	}

	return str
}

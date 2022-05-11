package mentalese

type EntitySorts map[string]string

func NewEntitySorts() *EntitySorts {
	return &EntitySorts{}
}

func (s *EntitySorts) AddSort(variable string, sort string) {
	(*s)[variable] = sort
}

func (s *EntitySorts) GetSort(variable string) (bool, string) {
	sort, found := (*s)[variable]
	return found, sort
}

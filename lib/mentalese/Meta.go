package mentalese

// maps a predicate to information about a relation
type Meta struct {
	sortProperties SortProperties
}

// for each argument a sort
type ArgumentSorts struct {
	Sorts []string
}

func NewMeta() *Meta {
	return &Meta{
		sortProperties: SortProperties{},
	}
}

func (meta Meta) AddSortInfo(name string, info SortProperty) {
	meta.sortProperties[name] = info
}

func (meta Meta) GetSorts() SortProperties {
	return meta.sortProperties
}

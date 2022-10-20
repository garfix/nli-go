package mentalese

// maps a predicate to information about a relation
type Meta struct {
	argumentSorts  map[string]ArgumentSorts
	sortProperties SortProperties
	sortHierachy   map[string][]string
}

// for each argument a sort
type ArgumentSorts struct {
	Sorts []string
}

func NewMeta() *Meta {
	return &Meta{
		argumentSorts:  map[string]ArgumentSorts{},
		sortProperties: SortProperties{},
	}
}

func (meta Meta) AddPredicate(name string, sorts []string) {
	meta.argumentSorts[name] = ArgumentSorts{
		Sorts: sorts,
	}
}

func (meta Meta) AddSortInfo(name string, info SortProperty) {
	meta.sortProperties[name] = info
}

func (meta Meta) GetSort(predicate string, argumentIndex int) string {

	pred, found := meta.argumentSorts[predicate]
	if found {
		return pred.Sorts[argumentIndex]
	}
	return ""
}

func (meta Meta) GetSorts() SortProperties {
	return meta.sortProperties
}

func (meta Meta) GetSortProperty(sort string) (SortProperty, bool) {
	info, found := meta.sortProperties[sort]
	return info, found
}

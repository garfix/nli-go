package mentalese

// maps a predicate to information about a relation
type Meta struct {
	predicates map[string]PredicateInfo
	sorts      Entities
	subSorts   map[string][]string
}

// for each argument a sort
type PredicateInfo struct {
	Sorts []string
}

func NewMeta() *Meta {
	return &Meta{
		predicates: map[string]PredicateInfo{},
		sorts:      Entities{},
		subSorts:   map[string][]string{},
	}
}

func (meta Meta) AddPredicate(name string, sorts []string) {
	meta.predicates[name] = PredicateInfo{
		Sorts: sorts,
	}
}

func (meta Meta) AddSortInfo(name string, info SortInfo) {
	meta.sorts[name] = info
}

func (meta Meta) GetSort(predicate string, argumentIndex int) string {

	pred, found := meta.predicates[predicate]
	if found {
		return pred.Sorts[argumentIndex]
	}
	return ""
}

func (meta Meta) GetSorts() Entities {
	return meta.sorts
}

func (meta Meta) GetSortInfo(sort string) (SortInfo, bool) {
	info, found := meta.sorts[sort]
	return info, found
}

func (meta Meta) AddSubSort(superSort string, subSort string) {

	_, found := meta.subSorts[subSort]
	if !found {
		meta.subSorts[subSort] = []string{}
	}

	meta.subSorts[subSort] = append(meta.subSorts[subSort], superSort)
}

func (meta Meta) MatchesSort(subSort string, superSort string) bool {

	// handles cases where there is no subSorts hierarchy, and even when there are no predicate defined
	if subSort == superSort {
		return true
	}

	subSortsTried := map[string]bool{}
	return meta.matchesSortRecursive(subSort, superSort, &subSortsTried)
}

func (meta Meta) GetMostSpecific(sort1 string, sort2 string) (string, bool) {
	if meta.MatchesSort(sort1, sort2) {
		return sort1, true
	} else if meta.MatchesSort(sort2, sort1) {
		return sort2, true
	} else {
		return "", false
	}
}

func (meta Meta) matchesSortRecursive(subSort string, superSort string, subSortsTried *map[string]bool) bool {

	found := false

	_, found = meta.subSorts[subSort]
	if !found { return false }

	if subSort == superSort { return true }

	for _, super := range meta.subSorts[subSort] {
		if super == superSort {
			return true
		} else {
			found = (*subSortsTried)[super]
			if found { return false }

			(*subSortsTried)[super] = true
			return meta.matchesSortRecursive(super, superSort, subSortsTried)
		}
	}

	return false
}
package mentalese

// maps a predicate to information about a relation
type Meta struct {
	predicates map[string]PredicateInfo
	entities   Entities
	subSorts   map[string][]string
}

// for each argument an entity type
type PredicateInfo struct {
	EntityTypes []string
}

func NewMeta() *Meta {
	return &Meta{
		predicates: map[string]PredicateInfo{},
		entities:   Entities{},
		subSorts:   map[string][]string{},
	}
}

func (meta Meta) AddPredicate(name string, entityTypes []string) {
	meta.predicates[name] = PredicateInfo{
		EntityTypes: entityTypes,
	}
}

func (meta Meta) AddEntityInfo(name string, entityInfo EntityInfo) {
	meta.entities[name] = entityInfo
}

func (meta Meta) GetEntityType(predicate string, argumentIndex int) string {

	pred, found := meta.predicates[predicate]
	if found {
		return pred.EntityTypes[argumentIndex]
	}
	return ""
}

func (meta Meta) GetEntities() Entities {
	return meta.entities
}

func (meta Meta) GetSortInfo(sort string) (EntityInfo, bool) {
	info, found := meta.entities[sort]
	return info, found
}

func (meta Meta) AddSort(superSort string, subSort string) {

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
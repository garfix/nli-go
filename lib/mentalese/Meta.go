package mentalese

// maps a predicate to information about a relation
type Meta struct {
	predicates map[string]PredicateInfo
	entities Entities
	sorts map[string][]string
}

// for each argument an entity type
type PredicateInfo struct {
	EntityTypes []string
}

func NewMeta() *Meta {
	return &Meta{
		predicates: map[string]PredicateInfo{},
		entities: Entities{},
		sorts: map[string][]string{},
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

func (meta Meta) AddSort(superSort string, subSort string) {

	_, found := meta.sorts[subSort]
	if !found {
		meta.sorts[subSort] = []string{}
	}

	meta.sorts[subSort] = append(meta.sorts[subSort], superSort)
}

func (meta Meta) MatchesSort(subSort string, superSort string) bool {

	// handles cases where there is no sorts hierarchy, and even when there are no predicate defined
	if subSort == superSort {
		return true
	}

	subSortsTried := map[string]bool{}
	return meta.matchesSortRecursive(subSort, superSort, &subSortsTried)
}

func (meta Meta) matchesSortRecursive(subSort string, superSort string, subSortsTried *map[string]bool) bool {

	found := false

	_, found = meta.sorts[subSort]
	if !found { return false }

	if subSort == superSort { return true }

	for _, super := range meta.sorts[subSort] {
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
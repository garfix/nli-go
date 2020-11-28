package central

import "nli-go/lib/mentalese"

type SortFinder struct {
	meta *mentalese.Meta
}

func NewSortFinder(meta *mentalese.Meta) SortFinder {
	return SortFinder{
		meta: meta,
	}
}

// Returns variable => sort bindings, or false, if they conflict
func (finder SortFinder) FindSorts(set mentalese.RelationSet) (mentalese.Sorts, bool) {
	sorts := mentalese.Sorts{}
	return finder.findSorts(set, sorts)
}

func (finder SortFinder) findSorts(set mentalese.RelationSet, sorts mentalese.Sorts) (mentalese.Sorts, bool) {

	for _, relation := range set {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() {
				sort := finder.meta.GetSort(relation.Predicate, i)
				if sort != "" {
					variable := argument.TermValue
					existing, found := sorts[variable]
					if found {
						specific, ok := finder.meta.GetMostSpecific(existing, sort)
						if ok {
							sort = specific
						} else {
							sorts[variable] = sorts[variable] + " & " + sort
							return sorts, false
						}
					}
					sorts[variable] = sort
				}
			} else if argument.IsRelationSet() {
				newSorts, ok := finder.findSorts(argument.TermValueRelationSet, sorts)
				if ok {
					sorts = newSorts
				} else {
					return sorts, false
				}
			} else if argument.IsRule() {
				// no need to implement
			} else if argument.IsList() {
				// no need to implement
			}
		}
	}
	return sorts, true
}

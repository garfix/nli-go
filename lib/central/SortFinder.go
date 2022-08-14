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

func (finder SortFinder) FindSorts(root *mentalese.ParseTreeNode) (mentalese.Sorts, bool) {
	sorts := mentalese.Sorts{}
	ok := finder.findSortsRecursive(root, &sorts)
	return sorts, ok
}

func (finder SortFinder) findSortsRecursive(node *mentalese.ParseTreeNode, sorts *mentalese.Sorts) bool {
	for _, childNode := range node.Constituents {
		childOk := finder.findSortsRecursive(childNode, sorts)
		if !childOk {
			return false
		}
	}
	return finder.findSortsInRelations(node.Rule.Sense, sorts)
}

func (finder SortFinder) findSortsInRelations(set mentalese.RelationSet, sorts *mentalese.Sorts) bool {
	for _, relation := range set {
		for i, argument := range relation.Arguments {
			if argument.IsVariable() {
				sort := finder.meta.GetSort(relation.Predicate, i)
				if sort != "" {
					variable := argument.TermValue
					existing, found := (*sorts)[variable]
					if found {
						specific, ok := finder.meta.GetMostSpecific(existing, sort)
						if ok {
							sort = specific
						} else {
							(*sorts)[variable] = (*sorts)[variable] + " & " + sort
							return false
						}
					}
					(*sorts)[variable] = sort
				}
			} else if argument.IsRelationSet() {
				ok := finder.findSortsInRelations(argument.TermValueRelationSet, sorts)
				if !ok {
					return false
				}
			} else if argument.IsRule() {
				// no need to implement
			} else if argument.IsList() {
				// no need to implement
			}
		}
	}
	return true
}

package central

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
)

type SortFinder struct {
	messenger api.ProcessMessenger
}

func NewSortFinder(messenger api.ProcessMessenger) SortFinder {
	return SortFinder{
		messenger: messenger,
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
	return finder.findSortsInTags(node.Rule.Tag, sorts) && finder.findSortsInRelations(node.Rule.Sense, sorts)
}

func (finder SortFinder) findSortsInRelations(set mentalese.RelationSet, sorts *mentalese.Sorts) bool {
	for _, relation := range set {

		if relation.Predicate == mentalese.PredicateHasSort || relation.Predicate == mentalese.PredicateIsa {
			variable := relation.Arguments[0].TermValue
			sort := relation.Arguments[1].TermValue
			existing, found := (*sorts)[variable]
			if found {
				specific, ok := finder.getMostSpecific(existing, sort)
				if ok {
					sort = specific
				} else {
					(*sorts)[variable] = (*sorts)[variable] + " & " + sort
					return false
				}
			}
			(*sorts)[variable] = sort
		}

		for _, argument := range relation.Arguments {
			if argument.IsRelationSet() {
				ok := finder.findSortsInRelations(argument.TermValueRelationSet, sorts)
				if !ok {
					return false
				}
			}
		}
	}
	return true
}

func (finder SortFinder) findSortsInTags(tags mentalese.RelationSet, sorts *mentalese.Sorts) bool {
	for _, tag := range tags {

		if tag.Predicate == mentalese.TagSort {
			variable := tag.Arguments[0].TermValue
			sort := tag.Arguments[1].TermValue
			existing, found := (*sorts)[variable]
			if found {
				specific, ok := finder.getMostSpecific(existing, sort)
				if ok {
					sort = specific
				} else {
					(*sorts)[variable] = (*sorts)[variable] + " & " + sort
					return false
				}
			}
			(*sorts)[variable] = sort
		}
	}
	return true
}

func (finder SortFinder) getMostSpecific(sort1 string, sort2 string) (string, bool) {

	if sort1 == sort2 {
		return sort1, true
	}

	isa := mentalese.NewRelation(false, mentalese.PredicateIsa, []mentalese.Term{
		mentalese.NewTermAtom(sort1),
		mentalese.NewTermAtom(sort2),
	})

	bindings := finder.messenger.ExecuteChildStackFrame(mentalese.RelationSet{isa}, mentalese.InitBindingSet(mentalese.NewBinding()))
	if !bindings.IsEmpty() {
		return sort1, true
	}

	isa2 := mentalese.NewRelation(false, mentalese.PredicateIsa, []mentalese.Term{
		mentalese.NewTermAtom(sort2),
		mentalese.NewTermAtom(sort1),
	})

	bindings2 := finder.messenger.ExecuteChildStackFrame(mentalese.RelationSet{isa2}, mentalese.InitBindingSet(mentalese.NewBinding()))
	if !bindings2.IsEmpty() {
		return sort2, true
	}

	return "", false
}

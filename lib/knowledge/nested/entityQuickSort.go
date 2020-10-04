package nested

// https://en.wikipedia.org/wiki/Quicksort

import (
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemNestedStructureBase) entityQuickSort(ids []mentalese.Term, orderFunction string) []mentalese.Term {
	base.entityQuickSortRange(&ids, 0, len(ids) - 1, orderFunction)
	return ids
}

func (base *SystemNestedStructureBase) entityQuickSortRange(ids *[]mentalese.Term, lo int, hi int, orderFunction string) {
	if lo < hi {
		p := base.partition(ids, lo, hi, orderFunction)
		base.entityQuickSortRange(ids, lo, p, orderFunction)
		base.entityQuickSortRange(ids, p + 1, hi, orderFunction)
	}
}

func (base *SystemNestedStructureBase) partition(ids *[]mentalese.Term, lo int, hi int, orderFunction string) int {
	pivot := (hi + lo) / 2
	pivotId := (*ids)[pivot]
	i := lo - 1
	j := hi + 1
	for {
		for {
			i = i + 1
			id := (*ids)[i]
			if base.compare(id, pivotId, orderFunction) >= 0 { break }
		}
		for {
			j = j - 1
			id := (*ids)[j]
			if base.compare(id, pivotId, orderFunction) <= 0 { break }
		}
		if i >= j {
			return j
		}
		// swap id i with id j
		temp := (*ids)[i]; (*ids)[i] = (*ids)[j]; (*ids)[j] = temp
	}
}

func (base *SystemNestedStructureBase) compare(id1 mentalese.Term, id2 mentalese.Term, orderFunction string) int {

	relation := mentalese.NewRelation(true, orderFunction, []mentalese.Term{
		mentalese.NewTermVariable("E1"),
		mentalese.NewTermVariable("E2"),
		mentalese.NewTermVariable("R"),
	})

	bindings := base.solver.SolveRelationSet(mentalese.RelationSet{ relation }, mentalese.Bindings{ mentalese.Binding{
		"E1": id1,
		"E2": id2,
	} })
	values := bindings.GetDistinctValues("R")

	if len(values) != 1 {
		base.log.AddError("order compare function " + orderFunction + " returned " + strconv.Itoa(len(values)) + " values")
		return 0
	}

	r := values[0]
	value, err := strconv.Atoi(r.TermValue)

	if err != nil {
		base.log.AddError("order compare function " + orderFunction + " returned " + r.TermValue)
		return 0
	}

	return value
}
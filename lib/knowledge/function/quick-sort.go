package function

// https://en.wikipedia.org/wiki/Quicksort

import (
	"nli-go/lib/api"
	"nli-go/lib/mentalese"
	"strconv"
)

func (base *SystemSolverFunctionBase) entityQuickSort(messenger api.ProcessMessenger, ids []mentalese.Term, orderFunction string) ([]mentalese.Term, bool) {
	base.entityQuickSortRange(messenger, &ids, 0, len(ids)-1, orderFunction)
	return ids, false
}

func (base *SystemSolverFunctionBase) entityQuickSortRange(messenger api.ProcessMessenger, ids *[]mentalese.Term, lo int, hi int, orderFunction string) bool {
	if lo < hi {
		p := 0
		p, _ = base.partition(messenger, ids, lo, hi, orderFunction)
		base.entityQuickSortRange(messenger, ids, lo, p, orderFunction)
		base.entityQuickSortRange(messenger, ids, p+1, hi, orderFunction)
	}
	return false
}

func (base *SystemSolverFunctionBase) partition(messenger api.ProcessMessenger, ids *[]mentalese.Term, lo int, hi int, orderFunction string) (int, bool) {
	pivot := (hi + lo) / 2
	pivotId := (*ids)[pivot]
	i := lo - 1
	j := hi + 1
	result := 0
	for {
		for {
			i = i + 1
			id := (*ids)[i]
			result, _ = base.compare(messenger, id, pivotId, orderFunction)
			if result >= 0 {
				break
			}
		}
		for {
			j = j - 1
			id := (*ids)[j]
			result, _ = base.compare(messenger, id, pivotId, orderFunction)
			if result <= 0 {
				break
			}
			if j == 0 {
				break
			}
		}
		if i >= j {
			return j, false
		}
		// swap id i with id j
		temp := (*ids)[i]
		(*ids)[i] = (*ids)[j]
		(*ids)[j] = temp
	}
}

func (base *SystemSolverFunctionBase) compare(messenger api.ProcessMessenger, id1 mentalese.Term, id2 mentalese.Term, orderFunction string) (int, bool) {

	// special case to save time
	if id1.Equals(id2) {
		return 0, false
	}

	relation := mentalese.NewRelation(false, orderFunction, []mentalese.Term{
		mentalese.NewTermVariable("E1"),
		mentalese.NewTermVariable("E2"),
		mentalese.NewTermVariable("R"),
	})

	b := mentalese.NewBinding()
	b.Set("E1", id1)
	b.Set("E2", id2)

	bindings := messenger.ExecuteChildStackFrame(mentalese.RelationSet{relation}, mentalese.InitBindingSet(b))

	values := bindings.GetDistinctValues("R")

	if len(values) != 1 {
		base.log.AddError("order compare function " + orderFunction + " returned " + strconv.Itoa(len(values)) + " values")
		return 0, false
	}

	r := values[0]
	value, err := strconv.Atoi(r.TermValue)

	if err != nil {
		base.log.AddError("order compare function " + orderFunction + " returned " + r.TermValue)
		return 0, false
	}

	return value, false
}

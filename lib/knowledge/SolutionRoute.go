package knowledge

import "nli-go/lib/mentalese"

type SolutionRoute []RelationGroup
type SolutionRoutes []SolutionRoute

func (s SolutionRoute) Equals(t SolutionRoute) bool {
	equals := true

	if len(s) != len(t) {
		return false
	}

	for i, group := range s {
		equals = equals && group.Equals(t[i])
	}

	return equals
}

func (s SolutionRoute) Len() int {
	return len(s)
}

func (s SolutionRoute) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SolutionRoute) Less(i, j int) bool {
	return s[i].Cost < s[j].Cost
}

func (s SolutionRoute) GetCombinedRelations() mentalese.RelationSet {

	relations := mentalese.RelationSet{}

	for _, group := range s {
		relations = append(relations, group.Relations...)
	}

	return relations
}

func (s SolutionRoute) GetTotalRelationCount() int {

	count := 0

	for _, group := range s {
		count += len(group.Relations)
	}

	return count
}

func (s SolutionRoute) String() string {

	str := "["
	sep := ""

	for _, group := range s {
		str += sep + group.String()
		sep = ", "
	}

	str += "]"

	return str
}

func (s SolutionRoutes) String() string {

	str := "["
	sep := "\n "

	for _, route := range s {
		str += sep + route.String()
	}

	str += "]"

	return str
}
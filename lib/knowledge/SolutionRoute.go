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
	if s[i].Cost < s[j].Cost {
		return true
	} else if s[i].Cost > s[j].Cost {
		return false
	} else {

		if s[i].KnowledgeBaseIndex < s[j].KnowledgeBaseIndex {
			return true
		} else if s[i].KnowledgeBaseIndex > s[j].KnowledgeBaseIndex {
			return false
		} else {

			return s[i].Relations.String() < s[j].Relations.String()

		}

	}
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
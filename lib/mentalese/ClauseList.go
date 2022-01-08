package mentalese

type ClauseList struct {
	Clauses []*Clause
}

func NewClauseList() *ClauseList {
	return &ClauseList{
		Clauses: []*Clause{},
	}
}

func (e *ClauseList) Clear() {
	e.Clauses = []*Clause{}
}

func (e *ClauseList) AddClause(clause *Clause) {
	e.Clauses = append(e.Clauses, clause)
}

func (e *ClauseList) GetLastClause() *Clause {
	if len(e.Clauses) == 0 {
		return nil
	} else {
		return e.Clauses[len(e.Clauses)-1]
	}
}

package mentalese

type AnaphoraQueue struct {
	clauses []*AnaphoraQueueClause
}

func NewAnaphoraQueue() *AnaphoraQueue {
	return &AnaphoraQueue{}
}

func (q *AnaphoraQueue) GetClauses() []*AnaphoraQueueClause {
	return q.clauses
}

func (q *AnaphoraQueue) StartClause() {
	q.clauses = append(q.clauses, NewAnaphoraQueueClause())
}

func (q *AnaphoraQueue) GetActiveClause() *AnaphoraQueueClause {
	if len(q.clauses) == 0 {
		return nil
	} else {
		return q.clauses[len(q.clauses)-1]
	}
}

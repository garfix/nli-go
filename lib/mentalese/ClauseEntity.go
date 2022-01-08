package mentalese

type ClauseEntity struct {
	DiscourseVariable string
	SyntacticFunction string
}

func NewClauseEntity(variable string, function string) *ClauseEntity {
	return &ClauseEntity{
		DiscourseVariable: variable,
		SyntacticFunction: function,
	}
}

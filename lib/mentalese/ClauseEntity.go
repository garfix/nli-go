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

func (e *ClauseEntity) Replacevariable(oldVariable string, newVariable string) {
	if e.DiscourseVariable == oldVariable {
		e.DiscourseVariable = newVariable
	}
}

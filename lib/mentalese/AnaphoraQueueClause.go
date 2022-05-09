package mentalese

type AnaphoraQueueClause struct {
	discourseVariables []string
}

func NewAnaphoraQueueClause() *AnaphoraQueueClause {
	return &AnaphoraQueueClause{
		discourseVariables: []string{},
	}
}

func (c *AnaphoraQueueClause) AddDialogVariable(variable string) {
	c.discourseVariables = append(c.discourseVariables, variable)
}

func (c *AnaphoraQueueClause) GetDiscourseVariables() []string {
	return c.discourseVariables
}

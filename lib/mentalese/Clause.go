package mentalese

type Clause struct {
	AuthorIsSystem bool
	ParseTree      *ParseTreeNode
	Entities       []*ClauseEntity
	Center         *ClauseEntity
}

func NewClause(parseTree *ParseTreeNode, authorIsSystem bool, entities []*ClauseEntity) *Clause {

	return &Clause{
		AuthorIsSystem: authorIsSystem,
		ParseTree:      parseTree,
		Entities:       entities,
	}
}

func ExtractEntities(node *ParseTreeNode) []*ClauseEntity {
	entities := []*ClauseEntity{}

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == TagFunction {
			variable := tag.Arguments[0].TermValue
			value := tag.Arguments[1].TermValue
			entities = append(entities, NewClauseEntity(variable, value))
		}
	}

	for _, constituent := range node.Constituents {
		entities = append(entities, ExtractEntities(constituent)...)
	}

	return entities
}

func (c *Clause) UpdateCenter(list *ClauseList, binding *Binding) {
	var previousCenter *ClauseEntity = nil
	var center *ClauseEntity = nil
	var priority = 0

	previousClause := list.GetPreviousClause()
	if previousClause != nil {
		previousCenter = previousClause.Center
	}

	priorities := map[string]int{
		"previousCenter":    100,
		AtomFunctionSubject: 10,
		AtomFunctionObject:  5,
	}

	// new clause has no entities? keep existing center
	if len(c.Entities) == 0 {
		center = previousCenter
	}

	for _, entity := range c.Entities {
		if previousCenter != nil && getValue(entity.DiscourseVariable, binding) == getValue(previousCenter.DiscourseVariable, binding) {
			priority = priorities["previousCenter"]
			center = entity
			continue
		}
		prio, found := priorities[entity.SyntacticFunction]
		if found {
			if prio > priority {
				priority = prio
				center = entity
			}
		}
	}

	c.Center = center
}

func getValue(variable string, binding *Binding) string {
	v, found := binding.Get(variable)
	if found {
		return v.TermValue
	} else {
		return ""
	}
}

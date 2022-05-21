package mentalese

import "nli-go/lib/common"

type Clause struct {
	AuthorIsSystem   bool
	ParseTree        *ParseTreeNode
	Functions        []*ClauseEntity
	Center           string
	ResolvedEntities []string
}

func NewClause(parseTree *ParseTreeNode, authorIsSystem bool, entities []*ClauseEntity) *Clause {

	return &Clause{
		AuthorIsSystem:   authorIsSystem,
		ParseTree:        parseTree,
		Functions:        entities,
		ResolvedEntities: []string{},
	}
}

func (clause *Clause) ReplaceVariable(fromVariable string, toVariable string) {
	newTree := clause.ParseTree.ReplaceVariable(fromVariable, toVariable)
	clause.ParseTree = &newTree

	if clause.Center == fromVariable {
		clause.Center = toVariable
	}
	for _, e := range clause.Functions {
		e.Replacevariable(fromVariable, toVariable)
	}
}

func ExtractEntities(node *ParseTreeNode) []*ClauseEntity {
	variables := collectVariables(node)
	functions := collectFunctions(node)
	entities := createOrderedEntities(variables, functions)

	return entities
}

func collectVariables(node *ParseTreeNode) []string {
	variables := []string{}

	for _, entityVariables := range node.Rule.EntityVariables {
		for _, entityVariable := range entityVariables {
			if entityVariable == Terminal {
				continue
			}
			if !common.StringArrayContains(variables, entityVariable) {
				variables = append(variables, entityVariable)
			}
		}
	}

	for _, constituent := range node.Constituents {
		for _, entityVariable := range collectVariables(constituent) {
			if !common.StringArrayContains(variables, entityVariable) {
				variables = append(variables, entityVariable)
			}
		}
	}

	return variables
}

func collectFunctions(node *ParseTreeNode) map[string]string {
	functions := map[string]string{}

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == TagFunction {
			variable := tag.Arguments[0].TermValue
			function := tag.Arguments[1].TermValue
			functions[variable] = function
		}
	}

	for _, constituent := range node.Constituents {
		childFunctions := collectFunctions(constituent)
		for variable, function := range childFunctions {
			existingFunction, found := functions[variable]
			if found && existingFunction != function {
				// todo handle better
				panic(variable + " cannot be both " + existingFunction + " and " + function)
			}
			functions[variable] = function
		}
	}

	return functions
}

func createOrderedEntities(variables []string, functions map[string]string) []*ClauseEntity {
	entities := []*ClauseEntity{}

	allFunctions := []string{AtomFunctionSubject, AtomFunctionObject, AtomFunctionNone}

	for _, aFunction := range allFunctions {
		for _, variable := range variables {
			function, found := functions[variable]
			if !found {
				function = AtomFunctionNone
			}
			if function == aFunction {
				entities = append(entities, NewClauseEntity(variable, function))
			}
		}
	}

	return entities
}

func (c *Clause) AddEntity(entity string) {
	c.ResolvedEntities = append(c.ResolvedEntities, entity)
}

func (c *Clause) UpdateCenter(list *ClauseList, binding *Binding) {
	var previousCenter = ""
	var center = ""
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
	if len(c.Functions) == 0 {
		center = previousCenter
	}

	for _, entity := range c.Functions {
		if previousCenter != "" {
			a := getValue(entity.DiscourseVariable, binding)
			b := getValue(previousCenter, binding)
			if a == b {
				priority = priorities["previousCenter"]
				center = entity.DiscourseVariable
				continue
			}
		}
		prio, found := priorities[entity.SyntacticFunction]
		if found {
			if prio > priority {
				priority = prio
				center = entity.DiscourseVariable
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

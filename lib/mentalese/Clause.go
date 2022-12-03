package mentalese

import (
	"nli-go/lib/common"
)

type Clause struct {
	AuthorIsSystem     bool
	ParseTree          *ParseTreeNode
	SyntacticFunctions []*ClauseEntity
	// entities as they are encountered when resolving anaphora, used to build the anaphora queue
	QueuedEntities []string
}

func NewClause(parseTree *ParseTreeNode, authorIsSystem bool, syntacticFunctions []*ClauseEntity) *Clause {

	return &Clause{
		AuthorIsSystem:     authorIsSystem,
		ParseTree:          parseTree,
		SyntacticFunctions: syntacticFunctions,
		QueuedEntities:     []string{},
	}
}

func (clause *Clause) ReplaceVariable(fromVariable string, toVariable string) {
	newTree := clause.ParseTree.ReplaceVariable(fromVariable, toVariable)
	clause.ParseTree = newTree

	for _, e := range clause.SyntacticFunctions {
		e.Replacevariable(fromVariable, toVariable)
	}
}

func ExtractSyntacticFunctions(node *ParseTreeNode) []*ClauseEntity {
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
	c.QueuedEntities = append(c.QueuedEntities, entity)
}

package mentalese

import "strconv"

type VariableGenerator struct {
	variables map[string]int
}

func NewVariableGenerator() *VariableGenerator {
	return &VariableGenerator{
		variables: map[string]int{},
	}
}

func (gen  *VariableGenerator) GenerateVariable(initial string) Term {

	_, present := gen.variables[initial]
	if !present {
		gen.variables[initial] = 1
	} else {
		gen.variables[initial]++
	}

	return NewTermVariable(initial + "$" + strconv.Itoa(gen.variables[initial]))
}

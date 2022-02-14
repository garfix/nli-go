package mentalese

import (
	"strconv"
	"sync"
)

type VariableGenerator struct {
	Variables map[string]int
	mutex     sync.Mutex
}

func NewVariableGenerator() *VariableGenerator {
	return &VariableGenerator{
		Variables: map[string]int{},
	}
}

func (gen *VariableGenerator) Initialize() {
	gen.Variables = map[string]int{}
}

func (gen *VariableGenerator) GenerateVariable(name string) Term {

	gen.mutex.Lock()

	baseName := name
	length := len(name)

	for i := length - 1; i > 0; i-- {
		if baseName[i] >= '0' && baseName[i] <= '9' {
			baseName = baseName[0 : length-1]
		}
	}

	_, present := gen.Variables[baseName]
	if !present {
		gen.Variables[baseName] = 1
	} else {
		gen.Variables[baseName]++
	}

	gen.mutex.Unlock()

	return NewTermVariable(baseName + "$" + strconv.Itoa(gen.Variables[baseName]))
}

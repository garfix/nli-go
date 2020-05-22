package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type ShellBase struct {
	KnowledgeBaseCore
	log     	*common.SystemLog
}

func NewShellBase(name string, log *common.SystemLog) *ShellBase {
	return &ShellBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name },
		log: log,
	}
}

func (factBase *ShellBase) HandlesPredicate(predicate string) bool {
	predicates := []string{"list"}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *ShellBase) list(input mentalese.Relation, binding mentalese.Binding) mentalese.Binding {

	bound := input.BindSingle(binding)

	//if !base.validate(bound, "iiv") {
	//	return nil
	//}

	int1, _ := strconv.Atoi(bound.Arguments[0].TermValue)
	int2, _ := strconv.Atoi(bound.Arguments[1].TermValue)

	result := int1 + int2

	newBinding := binding.Copy()
	newBinding[input.Arguments[2].TermValue] = mentalese.NewString(strconv.Itoa(result))

	return newBinding
}

func (base *ShellBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	newBinding := binding
	found := true

	switch input.Predicate {
	case "list":
		newBinding = base.list(input, binding)
	default:
		found = false
	}

	return newBinding, found
}

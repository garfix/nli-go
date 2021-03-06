package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"os/exec"
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
	predicates := []string{mentalese.PredicateExec, mentalese.PredicateExecResponse}

	for _, p := range predicates {
		if p == predicate {
			return true
		}
	}
	return false
}

func (base *ShellBase) exec(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)

	if !Validate(bound, "S", base.log) {
		return binding, false
	}

	command := bound.Arguments[0].TermValue
	args := []string{}
	for i := range bound.Arguments {
		if i == 0 { continue }
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	return newBinding, true
}


func (base *ShellBase) execResponse(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool) {

	bound := input.BindSingle(binding)
	responseVar := input.Arguments[0].TermValue

	if !Validate(bound, "vS", base.log) {
		return binding, false
	}

	command := bound.Arguments[1].TermValue
	args := []string{}
	for i := range bound.Arguments {
		if i < 2 { continue }
		args = append(args, bound.Arguments[i].TermValue)
	}
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		base.log.AddError(err.Error())
	}

	newBinding := binding.Copy()

	newBinding.Set(responseVar, mentalese.NewTermString(string(output)))

	return newBinding, true
}

func (base *ShellBase) Execute(input mentalese.Relation, binding mentalese.Binding) (mentalese.Binding, bool, bool) {

	newBinding := binding
	found := true
	success := true

	switch input.Predicate {
	case mentalese.PredicateExec:
		newBinding, success = base.exec(input, binding)
	case mentalese.PredicateExecResponse:
		newBinding, success = base.execResponse(input, binding)
	default:
		found = false
	}

	return newBinding, found, success
}

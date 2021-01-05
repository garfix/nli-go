package central

import (
	"nli-go/lib/mentalese"
	"strconv"
)

type Call struct {
	relation mentalese.Relation
	singleBinding mentalese.Binding
	multipleBinding mentalese.BindingSet
	multiple bool
}

type CallStack struct {
	stack []Call
	empty []Call
}

func NewCallStack() *CallStack {
	return &CallStack{
		stack: []Call{},
		empty: []Call{},
	}
}

func (callStack *CallStack) PushSingle(relation mentalese.Relation, binding mentalese.Binding) {
	callStack.stack = append(callStack.stack, Call{relation: relation, singleBinding: binding, multiple: false})
	callStack.empty = []Call{}
}

func (callStack *CallStack) PushMultiple(relation mentalese.Relation, bindings mentalese.BindingSet) {
	callStack.stack = append(callStack.stack, Call{relation: relation, multipleBinding: bindings, multiple: true})
	callStack.empty = []Call{}
}

func (callStack *CallStack) Pop(bindings mentalese.BindingSet) {
	// empty result set? save the stack
	// except when it was saved before on a deeper level
	if bindings.IsEmpty() && len(callStack.empty) == 0 {
		callStack.empty = []Call{}
		for _, call := range callStack.stack {
			callStack.empty = append(callStack.empty, call)
		}
	}
	callStack.stack = callStack.stack[0:len(callStack.stack) - 1]
}

func (callStack *CallStack) String() string {
	str := ""

	len := len(callStack.empty)
	for i := len - 1; i >= 0; i-- {
		call := callStack.empty[i]
		str += strconv.Itoa(i + 1) + ". " + call.relation.String() + "\n"
		if call.multiple {
			str += call.multipleBinding.String() + "\n\n"
		} else {
			str += call.singleBinding.String() + "\n\n"
		}
	}

	return str
}
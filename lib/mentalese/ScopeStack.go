package mentalese

type ScopeStack struct {
	scopes []*Scope
	current *Scope
}

func NewScopeStack() *ScopeStack {
	global := NewScope()
	return &ScopeStack{
		scopes:[]*Scope{ global },
		current: global,
	}
}

func (stack *ScopeStack) Push(scope *Scope) {
	stack.scopes = append(stack.scopes, scope)
	stack.current = scope
}

func (stack *ScopeStack) Pop() {
	lastIndex := len(stack.scopes) - 1
	if lastIndex == 0 {
		panic("Cannot pop empty scope stack")
	}
	stack.scopes = stack.scopes[0 : lastIndex]
	stack.current = stack.scopes[lastIndex - 1]
}

func (stack *ScopeStack) GetCurrentScope() *Scope {
	return stack.current
}
package mentalese

type ScopeStack struct {
	scopes []*Scope
}

func NewScopeStack() *ScopeStack {
	return &ScopeStack{}
}

func (stack *ScopeStack) Push(scope *Scope) {
	stack.scopes = append(stack.scopes, scope)
}

func (stack *ScopeStack) Pop() {
	stack.scopes = stack.scopes[0 : len(stack.scopes) - 1]
}
func (stack *ScopeStack) GetCurrentScope() *Scope {
	length := len(stack.scopes)
	if length == 0 {
		return nil
	} else {
		return stack.scopes[length - 1]
	}
}
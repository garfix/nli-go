package mentalese

type Scope struct {
	variables Binding
}

func NewScope() *Scope {
	return &Scope{ variables: NewBinding() }
}

func (scope *Scope) GetVariables() *Binding {
	return &scope.variables
}

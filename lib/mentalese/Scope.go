package mentalese

type Scope struct {
	variables Binding
	isBreaked bool
}

func NewScope() *Scope {
	return &Scope{ variables: NewBinding() }
}

func (scope *Scope) GetVariables() *Binding {
	return &scope.variables
}

func (scope *Scope) SetBreaked(breaked bool) {
	scope.isBreaked = breaked
}

func (scope *Scope) IsBreaked() bool {
	return scope.isBreaked
}
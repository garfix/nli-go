package mentalese

type EntityBindings struct {
	binding Binding
}

func NewEntityBindings() *EntityBindings {
	return &EntityBindings{
		binding: NewBinding(),
	}
}

func (b *EntityBindings) Set(variable string, value Term) {
	b.binding.Set(variable, value)
}

func (b *EntityBindings) Get(variable string) (Term, bool) {
	return b.binding.Get(variable)
}

func (b *EntityBindings) Copy() *EntityBindings {
	return &EntityBindings{
		binding: b.binding.Copy(),
	}
}

func (b *EntityBindings) Clear() {
	b.binding.Clear()
}

func (b *EntityBindings) Merge(binding Binding) Binding {
	return b.binding.Merge(binding)
}

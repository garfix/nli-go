package mentalese

type FactBase interface {
	Bind(goal Relation) []Binding
}

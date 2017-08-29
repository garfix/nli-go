package mentalese

type FactBase interface {
	Bind(goal []Relation) ([]Binding, bool)
	GetMappings() []DbMapping
	GetStatistics() DbStats
}

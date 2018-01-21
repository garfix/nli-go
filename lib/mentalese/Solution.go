package mentalese

type Solution struct {
	Condition   RelationSet
	Transformations []RelationTransformation
	NoResults   ResultHandler
	SomeResults ResultHandler
}

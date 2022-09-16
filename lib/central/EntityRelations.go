package central

import "nli-go/lib/mentalese"

type EntityDefinitions struct {
	definitions map[string]mentalese.RelationSet
}

func NewEntityDefinitions() *EntityDefinitions {
	return &EntityDefinitions{
		definitions: map[string]mentalese.RelationSet{},
	}
}

func (d *EntityDefinitions) Add(variable string, definition mentalese.RelationSet) {
	d.definitions[variable] = definition
}

func (d *EntityDefinitions) Get(variable string) mentalese.RelationSet {
	definition, found := d.definitions[variable]
	if found {
		return definition
	} else {
		return mentalese.RelationSet{}
	}
}

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

func (e *EntityDefinitions) Clear() {
	e.definitions = map[string]mentalese.RelationSet{}
}

func (p *EntityDefinitions) Copy() *EntityDefinitions {

	newDefinitions := map[string]mentalese.RelationSet{}
	for k, v := range p.definitions {
		newDefinitions[k] = v
	}

	return &EntityDefinitions{
		definitions: newDefinitions,
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

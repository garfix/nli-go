package mentalese

type EntityDefinitions struct {
	definitions map[string]RelationSet
}

func NewEntityDefinitions() *EntityDefinitions {
	return &EntityDefinitions{
		definitions: map[string]RelationSet{},
	}
}

func (e *EntityDefinitions) Clear() {
	e.definitions = map[string]RelationSet{}
}

func (p *EntityDefinitions) Copy() *EntityDefinitions {

	newDefinitions := map[string]RelationSet{}
	for k, v := range p.definitions {
		newDefinitions[k] = v
	}

	return &EntityDefinitions{
		definitions: newDefinitions,
	}
}

func (d *EntityDefinitions) Add(variable string, definition RelationSet) {
	d.definitions[variable] = definition
}

func (d *EntityDefinitions) Get(variable string) RelationSet {
	definition, found := d.definitions[variable]
	if found {
		return definition
	} else {
		return RelationSet{}
	}
}

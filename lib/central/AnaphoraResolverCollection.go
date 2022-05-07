package central

import "nli-go/lib/mentalese"

type AnaphoraResolverCollection struct {
	output       string
	replacements map[string]string
	values       map[string]mentalese.Term
	sorts        map[string]mentalese.RelationSet
}

func NewAnaphoraResolverCollection() *AnaphoraResolverCollection {
	return &AnaphoraResolverCollection{
		output:       "",
		replacements: map[string]string{},
		values:       map[string]mentalese.Term{},
		sorts:        map[string]mentalese.RelationSet{},
	}
}

func (c *AnaphoraResolverCollection) AddReplacement(fromVariable string, toVariable string) {
	c.replacements[fromVariable] = toVariable
}

func (c *AnaphoraResolverCollection) AddReference(fromVariable string, value mentalese.Term) {
	c.values[fromVariable] = value
}

func (c *AnaphoraResolverCollection) AddSort(fromVariable string, sortRelations mentalese.RelationSet) {
	c.sorts[fromVariable] = sortRelations
}

package central

import "nli-go/lib/mentalese"

type AnaphoraResolverCollection struct {
	output       string
	replacements map[string]string
	references   map[string]mentalese.Term
}

func NewAnaphoraResolverCollection() *AnaphoraResolverCollection {
	return &AnaphoraResolverCollection{
		output:       "",
		replacements: map[string]string{},
		references:   map[string]mentalese.Term{},
	}
}

func (c *AnaphoraResolverCollection) AddReplacement(fromVariable string, toVariable string) {
	c.replacements[fromVariable] = toVariable
}

func (c *AnaphoraResolverCollection) AddReference(fromVariable string, value mentalese.Term) {
	c.references[fromVariable] = value
}

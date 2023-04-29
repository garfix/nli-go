package central

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type AnaphoraResolverCollection struct {
	output       string
	remark       string
	replacements map[string]string
	values       map[string]mentalese.Term
	oneAnaphors  map[string]mentalese.RelationSet
	// variables that are used in the same predication of a quantified body
	coArguments map[string][]string
}

func NewAnaphoraResolverCollection() *AnaphoraResolverCollection {
	return &AnaphoraResolverCollection{
		output:       "",
		remark:       "",
		replacements: map[string]string{},
		values:       map[string]mentalese.Term{},
		oneAnaphors:  map[string]mentalese.RelationSet{},
		coArguments:  map[string][]string{},
	}
}

func (c *AnaphoraResolverCollection) AddCoArgument(variable1 string, variable2 string) {
	_, found1 := c.coArguments[variable1]
	if !found1 {
		c.coArguments[variable1] = []string{}
	}
	_, found2 := c.coArguments[variable2]
	if !found2 {
		c.coArguments[variable2] = []string{}
	}

	c.coArguments[variable1] = append(c.coArguments[variable1], variable2)
	c.coArguments[variable2] = append(c.coArguments[variable2], variable1)
}

func (c *AnaphoraResolverCollection) IsCoArgument(variable1 string, variable2 string) bool {
	coArguments, found := c.coArguments[variable1]
	if !found {
		return false
	}
	return common.StringArrayContains(coArguments, variable2)
}

func (c *AnaphoraResolverCollection) AddReplacement(fromVariable string, toVariable string) {
	c.replacements[fromVariable] = toVariable
}

func (c *AnaphoraResolverCollection) AddReference(fromVariable string, value mentalese.Term) {
	c.values[fromVariable] = value
}

func (c *AnaphoraResolverCollection) AddOneAnaphor(fromVariable string, value mentalese.RelationSet) {
	c.oneAnaphors[fromVariable] = value
}

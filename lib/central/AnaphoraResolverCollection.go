package central

type AnaphoraResolverCollection struct {
	output       string
	replacements map[string]string
}

func NewAnaphoraResolverCollection() *AnaphoraResolverCollection {
	return &AnaphoraResolverCollection{
		output:       "",
		replacements: map[string]string{},
	}
}

func (c *AnaphoraResolverCollection) AddReplacement(fromVariable string, toVariable string) {
	c.replacements[fromVariable] = toVariable
}

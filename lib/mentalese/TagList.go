package mentalese

type TagList struct {
	tags map[string]RelationSet
}

func NewTagList() *TagList {
	return &TagList{
		tags: map[string]RelationSet{},
	}
}

func (p *TagList) Clear() {
	p.tags = map[string]RelationSet{}
}

func (p *TagList) Copy() *TagList {

	newTags := map[string]RelationSet{}
	for k, v := range p.tags {
		newTags[k] = v
	}

	return &TagList{
		tags: newTags,
	}
}

func (p *TagList) ReplaceVariable(from string, to string) {
	newTags := map[string]RelationSet{}
	for variable, tagSet := range p.tags {
		if variable == from {
			newTags[to] = tagSet.ReplaceTerm(NewTermVariable(from), NewTermVariable(to))
		} else {
			newTags[variable] = tagSet.ReplaceTerm(NewTermVariable(from), NewTermVariable(to))
		}
	}
	p.tags = newTags
}

func (p *TagList) AddTags(tags RelationSet) {
	for _, tag := range tags {
		variable := tag.Arguments[0].TermValue
		tags, found := p.tags[variable]
		if found {
			p.tags[variable] = append(tags, tag)
		} else {
			p.tags[variable] = RelationSet{tag}
		}
	}
}

func (p *TagList) GetTags(variable string) RelationSet {
	tags, found := p.tags[variable]
	if found {
		return tags
	} else {
		return RelationSet{}
	}
}

func (p *TagList) GetTagPredicates(variable string) []string {
	tags, found := p.tags[variable]
	if found {
		predicates := []string{}
		for _, tag := range tags {
			predicates = append(predicates, tag.Predicate)
		}
		return predicates
	} else {
		return []string{}
	}
}

func (p *TagList) GetTagsByPredicate(variable string, predicate string) RelationSet {
	predicates := RelationSet{}
	tags, found := p.tags[variable]
	if found {
		for _, tag := range tags {
			if tag.Predicate == predicate {
				predicates = append(predicates, tag)
			}
		}
	}
	return predicates
}

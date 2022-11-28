package central

import "nli-go/lib/mentalese"

type TagList struct {
	tags map[string]mentalese.RelationSet
}

func NewTagList() *TagList {
	return &TagList{
		tags: map[string]mentalese.RelationSet{},
	}
}

func (p *TagList) Clear() {
	p.tags = map[string]mentalese.RelationSet{}
}

func (p *TagList) Copy() *TagList {

	newTags := map[string]mentalese.RelationSet{}
	for k, v := range p.tags {
		newTags[k] = v
	}

	return &TagList{
		tags: newTags,
	}
}

func (p *TagList) ReplaceVariable(from string, to string) {
	newTags := map[string]mentalese.RelationSet{}
	for variable, tagSet := range p.tags {
		if variable == from {
			newTags[to] = tagSet.ReplaceTerm(mentalese.NewTermVariable(from), mentalese.NewTermVariable(to))
		} else {
			newTags[variable] = tagSet.ReplaceTerm(mentalese.NewTermVariable(from), mentalese.NewTermVariable(to))
		}
	}
	p.tags = newTags
}

func (p *TagList) AddTags(tags mentalese.RelationSet) {
	for _, tag := range tags {
		variable := tag.Arguments[0].TermValue
		tags, found := p.tags[variable]
		if found {
			p.tags[variable] = append(tags, tag)
		} else {
			p.tags[variable] = mentalese.RelationSet{tag}
		}
	}
}

func (p *TagList) GetTags(variable string) mentalese.RelationSet {
	tags, found := p.tags[variable]
	if found {
		return tags
	} else {
		return mentalese.RelationSet{}
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

func (p *TagList) GetTagsByPredicate(variable string, predicate string) mentalese.RelationSet {
	predicates := mentalese.RelationSet{}
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

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

func (p *TagList) ReplaceVariable(fromVariable string, toVariable string) {
	tags, found := p.tags[fromVariable]
	if found {
		delete(p.tags, fromVariable)
		p.tags[toVariable] = tags
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

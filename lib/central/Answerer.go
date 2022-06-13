package central

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

// The answerer takes a relation set in domain format
// and returns a relation set in domain format
// It uses Intent structures to determine how to act
type Answerer struct {
	intents []mentalese.Intent
	matcher *RelationMatcher
	log     *common.SystemLog
}

func NewAnswerer(matcher *RelationMatcher, log *common.SystemLog) *Answerer {

	return &Answerer{
		intents: []mentalese.Intent{},
		matcher: matcher,
		log:     log,
	}
}

func (answerer *Answerer) AddIntents(intents []mentalese.Intent) {
	answerer.intents = append(answerer.intents, intents...)
}

// Returns the solutions whose condition matches the goal, and a set of bindings per solution
func (answerer Answerer) FindIntents(goal mentalese.RelationSet) []mentalese.Intent {

	var intents []mentalese.Intent

	for _, anIntent := range answerer.intents {

		unScopedGoal := goal.UnScope()

		bindings, found := answerer.matcher.MatchSequenceToSet(anIntent.Condition, unScopedGoal, mentalese.NewBinding())
		if found {

			for _, binding := range bindings.GetAll() {
				boundIntent := anIntent.BindSingle(binding)
				intents = append(intents, boundIntent)
			}
		}
	}

	return intents
}

func (answerer Answerer) Build(template mentalese.RelationSet, bindings mentalese.BindingSet) mentalese.RelationSet {

	newSet := mentalese.RelationSet{}

	if bindings.IsEmpty() {
		newSet = template
	} else {

		sets := template.BindRelationSetMultipleBindings(bindings)

		newSet = mentalese.RelationSet{}
		for _, set := range sets {
			newSet = newSet.Merge(set)
		}
	}

	return newSet
}

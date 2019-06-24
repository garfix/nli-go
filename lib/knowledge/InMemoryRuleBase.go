package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
)

type InMemoryRuleBase struct {
	KnowledgeBaseCore
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewInMemoryRuleBase(name string, rules []mentalese.Rule, log *common.SystemLog) *InMemoryRuleBase {
	return &InMemoryRuleBase{
		KnowledgeBaseCore: KnowledgeBaseCore{ Name: name},
		rules: rules,
		log: log,
	}
}

func (ruleBase *InMemoryRuleBase) GetMatchingGroups(set mentalese.RelationSet, keyCabinet *mentalese.KeyCabinet) []RelationGroup {

	matchingGroups := []RelationGroup{}

	for _, rule := range ruleBase.rules {
		for _, setRelation := range set {
			if rule.Goal.Predicate == setRelation.Predicate {
// TODO calculate real costs
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, ruleBase.Name, worst_cost})
				break
			}
		}
	}

	return matchingGroups
}

func (ruleBase *InMemoryRuleBase) Bind(goal mentalese.Relation) ([]mentalese.RelationSet, mentalese.Bindings) {

	ruleBase.log.StartDebug("RuleBase BindSingle", goal)

	matcher := mentalese.NewRelationMatcher(ruleBase.log)
	subgoalRelationSets := []mentalese.RelationSet{}
	subgoalBindings := mentalese.Bindings{}

	for _, rule := range ruleBase.rules {

		// match goal
		aBinding, match := matcher.MatchTwoRelations(rule.Goal, goal, mentalese.Binding{})
		if match {
			subgoalRelationSets = append(subgoalRelationSets, rule.Pattern)
			subgoalBindings = append(subgoalBindings, aBinding)
		}
	}

	ruleBase.log.EndDebug("RuleBase BindSingle", subgoalRelationSets, subgoalBindings)

	return subgoalRelationSets, subgoalBindings
}

package mentalese

type Rules []Rule

func (rules Rules) Copy() Rules {
	newRules := []Rule{}

	for _, rule := range rules {
		newRules = append(newRules, rule.Copy())
	}

	return newRules
}
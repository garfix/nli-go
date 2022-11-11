package central

import "nli-go/lib/mentalese"

type AgreementChecker struct {
}

type Categories map[string]map[string]mentalese.Term

func (c *Categories) Get(entityVariable string) map[string]mentalese.Term {
	categories, found := (*c)[entityVariable]
	if !found {
		categories = map[string]mentalese.Term{}
		(*c)[entityVariable] = categories
	}
	return categories
}

func (c *Categories) Set(entityVariable string, agreementType string, agreementValue mentalese.Term) {
	_, found := (*c)[entityVariable]
	if !found {
		categories := map[string]mentalese.Term{}
		(*c)[entityVariable] = categories
	}
	(*c)[entityVariable][agreementType] = agreementValue
}

// --------------------------------

func NewAgreementChecker() *AgreementChecker {
	return &AgreementChecker{}
}

func (c *AgreementChecker) CheckAgreement(root *mentalese.ParseTreeNode, tagList *TagList) (bool, string) {
	categories := Categories{}
	output := ""

	entityVariables := root.GetVariablesRecursive()

	// collect categories for all entities,and make sure they don't conflict
	agree := c.checkWithinEntityAgreement(&categories, entityVariables, tagList)
	if !agree {
		return false, ""
	}

	// check between-entity agreements
	agree, output = c.checkBetweenEntityAgreement(root, &categories)

	return agree, output
}

func (c *AgreementChecker) checkWithinEntityAgreement(categories *Categories, entityVariables []string, tagList *TagList) bool {
	for _, variable := range entityVariables {
		for _, tag := range tagList.GetTagsByPredicate(variable, mentalese.TagCategory) {
			agreementType := tag.Arguments[1].TermValue
			agreementValue := tag.Arguments[2]
			agreed := c.matchCategories(categories, variable, agreementType, agreementValue)
			if !agreed {
				return false
			}
		}
	}

	return true
}

func (c *AgreementChecker) matchCategories(categories *Categories, variable string, agreementType string, agreementValue mentalese.Term) bool {
	entityCategories := categories.Get(variable)
	existingValue, found := entityCategories[agreementType]
	if found {
		if !existingValue.Equals(agreementValue) {
			return false
		}
	} else {
		categories.Set(variable, agreementType, agreementValue)
	}

	return true
}

func (c *AgreementChecker) checkBetweenEntityAgreement(node *mentalese.ParseTreeNode, categories *Categories) (bool, string) {

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == mentalese.TagAgree {
			variable1 := tag.Arguments[0].TermValue
			variable2 := tag.Arguments[1].TermValue
			v1Categories := categories.Get(variable1)
			v2Categories := categories.Get(variable2)

			for cat1, val1 := range v1Categories {
				val2, found := v2Categories[cat1]
				if found && !val2.Equals(val1) {
					return false, "Agreement mismatch: " + val1.TermValue + " / " + val2.TermValue
				}
			}
		}
	}

	for _, child := range node.Constituents {
		agreed, childOutput := c.checkBetweenEntityAgreement(child, categories)
		if !agreed {
			return false, childOutput
		}
	}

	return true, ""
}

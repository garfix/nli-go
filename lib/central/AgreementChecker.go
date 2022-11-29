package central

import "nli-go/lib/mentalese"

type AgreementChecker struct {
}

func NewAgreementChecker() *AgreementChecker {
	return &AgreementChecker{}
}

func (c *AgreementChecker) CheckAgreement(root *mentalese.ParseTreeNode, tagList *mentalese.TagList) (bool, string) {
	for _, variable := range root.GetVariablesRecursive() {
		// check for categories with multiple, conflicting values
		conflict, _, _ := c.CheckForCategoryConflictWithin(variable, tagList)
		if !conflict {
			return true, ""
		}
	}

	// check between-entity agreements
	agree, _, values := c.checkAgreementInTree(root, tagList)
	if !agree {
		return false, "Agreement mismatch: " + values[0].TermValue + " / " + values[1].TermValue
	}

	return true, ""
}

func (c *AgreementChecker) checkAgreementInTree(node *mentalese.ParseTreeNode, tagList *mentalese.TagList) (bool, string, []mentalese.Term) {

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == mentalese.TagAgree {
			variable1 := tag.Arguments[0].TermValue
			variable2 := tag.Arguments[1].TermValue
			v1Tags := tagList.GetTagsByPredicate(variable1, mentalese.TagCategory)
			v2Tags := tagList.GetTagsByPredicate(variable2, mentalese.TagCategory)

			for _, tag1 := range v1Tags {
				cat1 := tag1.Arguments[1]
				value1 := tag1.Arguments[2]
				for _, tag2 := range v2Tags {
					cat2 := tag2.Arguments[1]
					value2 := tag2.Arguments[2]

					if cat1.Equals(cat2) {
						if !value1.Equals(value2) {
							return false, cat1.TermValue, []mentalese.Term{value1, value2}
						}
					}
				}
			}
		}
	}

	for _, child := range node.Constituents {
		agreed, childCat, childValues := c.checkAgreementInTree(child, tagList)
		if !agreed {
			return false, childCat, childValues
		}
	}

	return true, "", []mentalese.Term{}
}

func (c *AgreementChecker) CheckForCategoryConflictWithin(variable string, tagList *mentalese.TagList) (bool, string, []mentalese.Term) {

	categoryTags := tagList.GetTagsByPredicate(variable, mentalese.TagCategory)
	categories := map[string]mentalese.Term{}

	for _, tag := range categoryTags {
		cat := tag.Arguments[1].TermValue
		value := tag.Arguments[2]
		existingValue, found := categories[cat]
		if found && !existingValue.Equals(value) {
			return false, cat, []mentalese.Term{existingValue, value}
		}
	}

	return true, "", []mentalese.Term{}
}

func (c *AgreementChecker) CheckForCategoryConflictBetween(variable1 string, variable2 string, tagList *mentalese.TagList) (bool, string, []mentalese.Term) {

	categoryTags1 := tagList.GetTagsByPredicate(variable1, mentalese.TagCategory)
	categoryTags2 := tagList.GetTagsByPredicate(variable2, mentalese.TagCategory)

	for _, tag1 := range categoryTags1 {
		cat1 := tag1.Arguments[1]
		value1 := tag1.Arguments[2]
		for _, tag2 := range categoryTags2 {
			cat2 := tag2.Arguments[1]
			value2 := tag2.Arguments[2]
			if cat1.Equals(cat2) {
				if !value1.Equals(value2) {
					return false, cat1.TermValue, []mentalese.Term{value1, value2}
				}
			}
		}
	}

	return true, "", []mentalese.Term{}
}

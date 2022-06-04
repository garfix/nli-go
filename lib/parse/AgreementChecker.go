package parse

import "nli-go/lib/mentalese"

type AgreementChecker struct {
}

type Agreements map[string]map[string]mentalese.Term

func NewAgreementChecker() *AgreementChecker {
	return &AgreementChecker{}
}

func (c *AgreementChecker) CheckAgreement(root *mentalese.ParseTreeNode) bool {
	agreements := Agreements{}
	return c.checkNode(root, &agreements)
}

func (c *AgreementChecker) checkNode(node *mentalese.ParseTreeNode, agreements *Agreements) bool {

	for _, tag := range node.Rule.Tag {
		if tag.Predicate == mentalese.TagAgree {
			variable := tag.Arguments[0].TermValue
			agreementType := tag.Arguments[1].TermValue
			agreementValue := tag.Arguments[2]

			agreed := c.checkAgreement(agreements, variable, agreementType, agreementValue)
			if !agreed {
				return false
			}
		}
	}

	for _, child := range node.Constituents {
		agreed := c.checkNode(child, agreements)
		if !agreed {
			return false
		}
	}

	return true
}

func (c *AgreementChecker) checkAgreement(agreements *Agreements, variable string, agreementType string, agreementValue mentalese.Term) bool {
	_, found1 := (*agreements)[variable]
	if found1 {
		existingValue, found2 := (*agreements)[variable][agreementType]
		if found2 {
			if !existingValue.Equals(agreementValue) {
				return false
			}
		} else {
			(*agreements)[variable][agreementType] = agreementValue
		}
	} else {
		(*agreements)[variable] = map[string]mentalese.Term{}
		(*agreements)[variable][agreementType] = agreementValue
	}

	return true
}
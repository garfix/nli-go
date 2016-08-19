package types

type Grammar interface {
	FindRules(antecedent string) [][]string
}

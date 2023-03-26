package parse

import "fmt"

type I18n struct {
	grammar *Grammar
}

func NewI18n(grammar *Grammar) *I18n {
	return &I18n{
		grammar: grammar,
	}
}

func (i18n *I18n) Translate(source string) string {
	translation := i18n.grammar.GetText(source)
	return translation
}

func (i18n *I18n) TranslateWithParam(source string, argument string) string {
	translation := i18n.grammar.GetText(source)
	return fmt.Sprintf(translation, argument)
}

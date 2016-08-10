package types

type Lexicon interface {
    CheckPartOfSpeech(word string, parOfSpeech string) bool
}

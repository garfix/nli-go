package types

type Tokenizer interface {
    Process(rawInput string) []string
}

package importer

import (
	"strings"
	"fmt"
)

const (
	service_tokenizer = "tokenizer"
	service_parser = "parser"
)

type ParseResult struct {
	Ok         bool
	Service    string
	Source     string
	LineNumber int
}

func (result ParseResult) String() string {
	lines := strings.Split(result.Source, "\n")
	lines[result.LineNumber - 1] = "* " + lines[result.LineNumber - 1]
	errorString := strings.Join(lines, "\n")

	return fmt.Sprintf("%s failed in line %d:\n%s", result.Service, result.LineNumber, errorString);
}

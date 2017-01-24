package importer

import (
	"strings"
	"fmt"
)

const (
	file_read = "file_read"
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
	line := result.LineNumber - 1
	if line < 0 {
		line = 0
	}
	lines[line] = "* " + lines[line]
	errorString := strings.Join(lines, "\n")

	return fmt.Sprintf("%s failed in line %d:\n%s", result.Service, result.LineNumber, errorString);
}

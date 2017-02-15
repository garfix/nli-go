package generate

import (
	"strings"
	"fmt"
)

type SurfaceRepresentation struct {
}

func NewSurfaceRepresentation() *SurfaceRepresentation {
	return &SurfaceRepresentation{}
}

func (r *SurfaceRepresentation) Create(tokens []string) string {

	result := ""

	for i, token := range tokens {

		if token == "." {
			result += token
		} else if i == 0 {
			result += fmt.Sprintf("%s%s", strings.Title(string(token[0])), token[1:])
		} else {
			result += " " + token
		}
	}

	return result
}
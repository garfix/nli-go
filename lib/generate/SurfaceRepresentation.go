package generate

import (
	"strings"
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
			result += strings.Title(token)
		} else {
			result += " " + token
		}
	}

	return result
}
package generate

import (
	"fmt"
	"nli-go/lib/common"
	"strings"
)

type SurfaceRepresentation struct {
	log *common.SystemLog
}

func NewSurfaceRepresentation(log *common.SystemLog) *SurfaceRepresentation {
	return &SurfaceRepresentation{log: log}
}

func (r *SurfaceRepresentation) Create(tokens []string) string {

	result := ""

	for i, token := range tokens {

		if token == "." {
			result += token
		} else if token == "," {
			result += token
		} else if i == 0 {
			if len(token) == 0 {
				result += "<empty>"
			} else if len(token) == 1 {
				result += strings.Title(string(token[0]))
			} else {
				result += fmt.Sprintf("%s%s", strings.Title(string(token[0])), token[1:])
			}
		} else {
			result += " " + token
		}
	}

	return result
}

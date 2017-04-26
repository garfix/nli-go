package tests

import (
	"testing"
    "nli-go/lib/global"
    "nli-go/lib/common"
)

func TestRelationships(t *testing.T) {

    log := global.NewSystemLog()
    system := global.NewSystem(common.Dir() + "/../../resources/relationships/config.json", log)

    if !log.IsOk() {
        t.Errorf(log.String())
    }

	var tests = []struct {
		question string
		answer   string
	} {
		{"Who married Jacqueline de Boer?", "Mark van Dongen married her"},
        {"Did Mark van Dongen marry Jacqueline de Boer?", "Yes"},
        {"Did Jacqueline de Boer marry Gerard van As?", "No"},
        {"Are Mark van Dongen and Suzanne van Dongen siblings?", "Yes"},
        {"Are Mark van Dongen and John van Dongen siblings?", "No"},
        {"Which children has John van Dongen?", "Mark van Dongen, Suzanne van Dongen, Dirk van Dongen and Durkje van Dongen"},
        {"How many children has John van Dongen?", "He has 4 children"},
        {"Does every parent have 4 children?", "Yes"},
        {"Does every parent have 3 children?", "No"},
	}

	for _, test := range tests {

		answer := system.Answer(test.question)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
		}
	}
}

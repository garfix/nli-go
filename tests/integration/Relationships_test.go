package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestRelationships(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/relationships/config.json", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"Who married Jacqueline de Boer?", "Mark van Dongen married her"},
		{"Did Mark van Dongen marry Jacqueline de Boer?", "Yes"},
		////{"Did Jacqueline de Boer marry Gerard van As?", "Name not found in any knowledge base: Gerard van As"},
		{"Are Mark van Dongen and Suzanne van Dongen siblings?", "Yes"},
		{"Are Mark van Dongen and John van Dongen siblings?", "No"},
		{"Which children has John van Dongen?", "Mark van Dongen, Suzanne van Dongen, Dirk van Dongen and Durkje van Dongen"},
		{"How many children has John van Dongen?", "He has 4 children"},
		{"Does every parent have 4 children?", "Yes"},
		{"Does every parent have 3 children?", "No"},
	}

	for _, test := range tests {

		log.Clear()

		answer, _ := system.Answer(test.question)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.String())
		}
	}
}

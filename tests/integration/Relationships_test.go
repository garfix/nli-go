package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestRelationships(t *testing.T) {

	log := common.NewSystemLog()
	//log.SetDebug(true)
	//log.SetPrint(true)
	system := global.NewSystem(common.Dir()+"/../../resources/relationships", "", common.Dir()+"/../../var", log, nil)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		// {"Who married Jacqueline de Boer?", "Mark van Dongen married her"},
		// {"Did Mark van Dongen marry Jacqueline de Boer?", "Yes"},
		// {"Did Jacqueline de Boer marry Gerard van As?", "Name not found: Gerard van As"},
		// {"Are Mark van Dongen and Suzanne van Dongen siblings?", "Yes"},
		// {"Are Mark van Dongen and John van Dongen siblings?", "No"},
		// {"Which children has John van Dongen?", "Mark van Dongen, Suzanne van Dongen, Dirk van Dongen and Durkje van Dongen"},
		// {"How many children has John van Dongen?", "He has 4 children"},

		// {"Does every parent have 4 children?", "Yes"},
		// {"Does every parent have 3 children?", "No"},

		// // assert
		// {"Suzanne van Dongen is married to Henk Smit", "Ok"},
		// {"Did Suzanne van Dongen marry Henk Smit?", "Yes"},

		// // "he" refers to "a man", not to the subject of the previous sentence
		// // thanks to gender agreement
		// {"Is Jacqueline de Boer married to a man?", "Yes"},
		// {"Who is he?", "Mark van Dongen"},

		// // // syntactic ambiguity; resolved by sortal restrictions
		// {"Was Jacqueline de Boer married by a river?", "Yes"},
		// {"Was Jacqueline de Boer married by Father Brown?", "Yes"},

		// // reflective pronouns
		// {"Is Jacqueline de Boer in love with herself?", "No"},
		// {"Is Henk Smit in love with himself?", "Yes"},

		// // number agreement mismatch
		{"Are Jacqueline de Boer married to a man?", "Agreement mismatch: plural / singular"},
	}

	for _, test := range tests {

		log.Clear()

		println(test.question)

		answer, _ := system.Answer(test.question)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.GetErrors())
			break
		}
	}
}

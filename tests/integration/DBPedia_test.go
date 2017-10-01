package tests

import (
	"testing"
	"nli-go/lib/common"
	"nli-go/lib/global"
)

func TestDBPedia(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/dbpedia/config-online.json", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"Who married Lord Byron?", "Anne Isabella Byron married him"},
		{"Who married Anne Isabella Milbanke?", "Lord Byron married her"},
		{"Who married Michael Jackson?", "Lisa Marie Presley and Debbie Rowe married him"},
	}

	for _, test := range tests {

		log.Clear()

		answer := system.Answer(test.question)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.String())
		}
	}
}
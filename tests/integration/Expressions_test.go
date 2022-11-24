package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestExpressions(t *testing.T) {

	log := common.NewSystemLog()
	system := global.NewSystem(common.Dir()+"/../../resources/expressions", "", common.Dir()+"/../../var", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"What is three plus four minus five", "2"},
		{"What is 3 plus 4 minus 5", "2"},
		{"What is 8 times 5", "40"},
		{"What is 18 minus 5 times 2", "8"},
		{"What is 5 times 4 minus 5", "15"},
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

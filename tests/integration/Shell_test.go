package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestShell(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/shell/config.json", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"List files page by page", "OK"},
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

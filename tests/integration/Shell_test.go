package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"strings"
	"testing"
)

// Test of a minimal NLI-GO application
func TestShell(t *testing.T) {

	log := common.NewSystemLog()
	system := global.NewSystem(common.Dir() + "/../../resources/shell", common.Dir() + "/../../var", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question string
		answer   string
	}{
		{"List files", "OK"},
		{"List files resources", "Shell_test.txt"},
	}

	for _, test := range tests {

		log.Clear()

		answer, _ := system.Answer(test.question)

		if !strings.Contains(answer, test.answer) {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.String())
		}
	}
}

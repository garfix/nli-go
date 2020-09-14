package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestExpressions(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/expressions", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	answer2, _ := system.Answer("What is 3 plus 4 minus five")

	if answer2 != "2" {
		t.Errorf(log.String())
	}
}

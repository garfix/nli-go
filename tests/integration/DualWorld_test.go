package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestDualWorld(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/dualworld", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	answer, _ := system.Answer("Which poem wrote the grandfather of Charles Darwin?")

	if answer != "The Loves of the Plants" {
		t.Errorf(log.String())
	}
}

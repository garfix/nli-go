package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestHelloWorld(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/helloworld", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	answer, _ := system.Answer("Hello world")

	if answer != "Welcome!" {
		t.Errorf(log.String())
	}
}

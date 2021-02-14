package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

// Test of a minimal NLI-GO application
func TestHelloWorld(t *testing.T) {

	log := common.NewSystemLog()
	system := global.NewSystem(common.Dir() + "/../../resources/helloworld", "", common.Dir() + "/../../var", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	answer, _ := system.AnswerAsync("Hello world")

	if answer != "Welcome!" {
		t.Errorf(log.String())
	}
}

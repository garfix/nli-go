package tests

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestDomain(t *testing.T) {

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
			{"Test A", "1"},
		},
	}

	log := common.NewSystemLog()
	log.SetDebug(true)
	//log.SetPrint(true)
	system := global.NewSystem(common.Dir() + "/../../resources/run", "run-demo", common.Dir() + "/../../var", log)

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	for _, session := range tests {

		system.ResetSession()

		for _, test := range session {

			log.Clear()

			system.CreateAnswerGoal(test.question)
			system.Run()
			//actions := system.ReadActions("print")

			answer := ""//actions[0].Text()

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				fmt.Println(log.String())
			}

		}
	}
}
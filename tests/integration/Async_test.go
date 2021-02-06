package tests

import (
	"fmt"
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestAsync(t *testing.T) {

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
			{"Test A", "1"},
			//{"Test B", "3"},
		},
	}

	log := common.NewSystemLog()
	log.SetDebug(true)
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
			actions := system.ReadActions("print")

			answer := ""
			if actions.GetLength() > 0 {
				answer = actions.Get(0).MustGet("Content").TermValue
			}

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				fmt.Println(log.String())
			}
		}
	}
}
package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"testing"
)

// Mimics some of SHRDLU's functions, but in the cli-go way

func TestBlocksWorld(t *testing.T) {
	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/shrdlu/config.json", log)
	sessionId := "1"
	actualSessionPath := common.AbsolutePath(common.Dir(), "sessions/" + sessionId + ".json")

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = [][]struct {
		question      string
		answer        string
		inSessionName string
		outSessionName string
	}{
		{
			{"Pick up a big red block", "OK", "", ""},
		},
	}

	for _, session := range tests {

		os.Remove(actualSessionPath)

		for _, test := range session {

			log.Clear()

			if test.inSessionName == "" {
				system.ClearDialogContext()
			} else {
				inSessionPath := common.AbsolutePath(common.Dir(), "resources/" + test.inSessionName)
				inSession, _ := common.ReadFile(inSessionPath)
				common.WriteFile(actualSessionPath, inSession)
				system.PopulateDialogContext(actualSessionPath)
			}

			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			system.StoreDialogContext(actualSessionPath)

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				t.Error(log.String())
			}

			if test.outSessionName != "" {
				outSessionPath := common.AbsolutePath(common.Dir(), "resources/"+test.outSessionName)
				expected, err := common.ReadFile(outSessionPath)

				if err != nil {
					t.Errorf("Test relationships: error reading %v", outSessionPath)
				}

				actual, err := common.ReadFile(actualSessionPath)

				if err != nil {
					t.Errorf("Test relationships: error reading %v", actualSessionPath)
				}

				if expected != actual {
					t.Errorf("Test relationships: got %v, want %v", actual, expected)
				}
			}
		}
	}
}

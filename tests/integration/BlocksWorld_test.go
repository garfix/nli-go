package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"testing"
)

// Mimics some of SHRDLU's functions, but in the nli-go way

func TestBlocksWorld(t *testing.T) {
	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/blocks/config.json", log)
	sessionId := "1"
	actualSessionPath := common.AbsolutePath(common.Dir(), "sessions/" + sessionId + ".json")

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
			// todo: move the green block on top of it away
			{"Pick up a big red block", "OK"},
			// todo "I don't understand which pyramid you mean"
			{"Grasp the pyramid", "I don't understand which one you mean"},
			// todo "By "it", I assume you mean the block which is taller than the one I am holding"
			{"Find a block which is taller than the one you are holding and put it into the box.", "OK"},
			// todo: the names of the objects could be generated; now they are explicitly added
			{"What does the box contain?", "The blue pyramid and the blue block"},
			{"What is the pyramid supported by?", "The box"},
			{"How many blocks are not in the box?", "Four of them"},
		},
		{
		},
	}

	os.Remove(actualSessionPath)

	for _, session := range tests {

		for _, test := range session {

			log.Clear()

			system.PopulateDialogContext(actualSessionPath, false)

			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			system.StoreDialogContext(actualSessionPath)

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				t.Error(log.String())
			}
		}
	}
}

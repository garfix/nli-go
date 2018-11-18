package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"os"
	"testing"
)

func TestDBPedia(t *testing.T) {

	log := common.NewSystemLog(false)
	system := global.NewSystem(common.Dir() + "/../../resources/dbpedia/config-online.json", log)
	sessionId := "1"
	absSessionPath := common.AbsolutePath(common.Dir(), "sessions/" + sessionId + ".json")

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	var tests = []struct {
		question    string
		answer      string
		sessionName string
	}{
		//{"Who married Lord Byron?", "Anne Isabella Byron married him", ""},
		//{"Who married Anne Isabella Milbanke?", "Lord Byron married her", ""},
		//{"Who married Michael Jackson?", "Lisa Marie Presley and Debbie Rowe married him", ""},
		//{"Who married Xyz Abc?", "I do not know", ""},
		//{"How many children had Lord Byron?", "He has 2 children", ""}, // Ada and Allegra
		//{"How many children has Madonna?", "She has 4 children", ""},
		//{"Who was Ada Lovelace's father?", "Lord Byron was her father", ""},
		//{"Who was Ada Lovelace's mother?", "Anne Isabella Byron was her mother", ""},
		//{"Who was Percy Florence Shelley's father?", "Percy Bysshe Shelley was his father", ""},
		//{"When was Lord Byron born?", "Which one? [dbpedia/http://dbpedia.org/resource/Lord_Byron] person; birth_date: 1788-01-22; birth_place: London [dbpedia/http://dbpedia.org/resource/Lord_Byron_(umpire)] person; birth_date: 1872-09-18; birth_place: New York City", ""},
		{"dbpedia/http://dbpedia.org/resource/Lord_Byron", "He was born on January 23, 1788", "session-1.json"},
	}

	for _, test := range tests {

		log.Clear()

		os.Remove(absSessionPath)

		if test.sessionName != "" {
			sessionPath := common.AbsolutePath(common.Dir(), "resources/" + test.sessionName)
			os.Link(sessionPath, absSessionPath)
		}

		system.PopulateDialogContext(absSessionPath)

		answer := system.Answer(test.question)

		system.StoreDialogContext(absSessionPath)

		if answer != test.answer {
			t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
			t.Error(log.String())
		}
	}
}
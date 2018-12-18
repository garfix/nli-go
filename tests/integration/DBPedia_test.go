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
		//{
		//	{"Who married Anne Isabella Milbanke?", "Lord Byron married her", "", ""},
		//	{"Who was Ada Lovelace's father?", "Lord Byron was her father", "", ""},
		//	{"Who was Ada Lovelace's mother?", "Anne Isabella Byron was her mother", "", ""},
		//	{"Who was Percy Florence Shelley's father?", "Percy Bysshe Shelley was his father", "", ""},
		//	{"Who married Xyz Abc?", "Name not found in any knowledge base: Xyz Abc", "", ""},
		//},
		{
			{"Who married Michael Jackson?", "Which one? [dbpedia/http://dbpedia.org/resource/Mariléia_dos_Santos] person; birth_date: 1963-11-19; birth_place: Brazil [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(footballer,_born_1980)] person; birth_date: 1980-06-26; birth_place: Cheltenham [dbpedia/http://dbpedia.org/resource/Michael_Jackson] person; birth_date: 1958-8-29; birth_place: City of Gary [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(American_Revolution)] person; birth_date: 1734-12-18 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(English_singer)] person; birth_date: 1964-1-1 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(TV_executive)] person; birth_date: 1958-2-11 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(actor)] person; birth_date: 1970-11-08; birth_place: Ottawa [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(anthropologist)] person; birth_date: 1940-1-1; birth_place: Nelson City [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(basketball)] person; birth_date: 1964-07-13; birth_place: Fairfax, Virginia [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(bishop)] person; birth_date: 1956-05-24 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(journalist)] person [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(linebacker)] person; birth_date: 1957-07-15; birth_place: Pasco, Washington [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(radio_commentator)] person; birth_date: 1934-04-16; birth_place: London [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(rugby_league)] person; birth_date: 1969-10-11 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(wide_receiver)] person; birth_date: 1969-04-12; birth_place: Tangipahoa, Louisiana [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(writer)] person; birth_date: 1942-03-27; birth_place: Wetherby [dbpedia/http://dbpedia.org/resource/2000–01_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/2002–03_Swansea_City_A.F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(footballer,_born_1973)] person; birth_date: 1973-12-04; birth_place: Runcorn [dbpedia/http://dbpedia.org/resource/1996–97_Bury_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/1996–97_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_A._Jackson] person; birth_date: 1936-1-1 [dbpedia/http://dbpedia.org/resource/1999–2000_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(engineer)] person", "", ""},
//			{"dbpedia/http://dbpedia.org/resource/Michael_Jackson", "Lisa Marie Presley and Debbie Rowe married him", "session-3.json", ""},
		},
		//{
		//	{"How many children has Madonna?", "She has 4 children", "", ""},
		//},
		//{
		//	{"When was Lord Byron born?", "Which one? [dbpedia/http://dbpedia.org/resource/Lord_Byron] person; birth_date: 1788-01-22; birth_place: London [dbpedia/http://dbpedia.org/resource/Lord_Byron_(umpire)] person; birth_date: 1872-09-18; birth_place: New York City", "", "session-1.json"},
		//	{"dbpedia/http://dbpedia.org/resource/Lord_Byron", "He was born on January 22, 1788", "session-1.json", "session-2.json"},
		//	{"Who married Lord Byron?", "Anne Isabella Byron married him", "session-2.json", ""},
		//	{"How many children had Lord Byron?", "He has 2 children", "session-2.json", ""}, // Ada and Allegra
		//},
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

			answer := system.Answer(test.question)

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
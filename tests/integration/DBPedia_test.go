package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestDBPedia(t *testing.T) {

	var tests = [][]struct {
		question      string
		answer        string
	}{
		{
			{"How many children has Madonna?", "She has 4 children"},
			{"How old is she", "She is 62 years old ( born on August 16, 1958 )"},
			//{"When is her birthday?", "August 16, 1958", "", ""},
		},
		{
			{"How many children has Madonna?", "She has 4 children"},
			{"When is Madonna's birthday?", "August 16, 1958"},
			{"Who is Madonna's husband?", "Sean Penn and Guy Ritchie"},
			{"Who is Sean Penn's wife?", "Robin Wright and Madonna"},
			{"How old is Madonna?", "She is 62 years old ( born on August 16, 1958 )"},
			{"How old is percy florence shelley?", "He was 70 years old ( born on November 12, 1819 ; died on December 05, 1889 )"},
			{"Who married Anne Isabella Milbanke?", "Lord Byron married her"},
			{"Who was Ada Lovelace's father?", "Lord Byron was her father"},
			{"Who was Ada Lovelace's mother?", "Anne Isabella Byron was her mother"},
			{"What is the name of Ada Lovelace's father?", "Lord Byron was her father"},
			{"Who was Percy Florence Shelley's father?", "Percy Bysshe Shelley was his father"},
			{"Who is Lisa Marie Presley?", "Lisa Marie Presley (born February 1, 1968) is an American singer-songwriter. She is the daughter of musician-actor Elvis Presley and actress and business magnate Priscilla Presley, and is Elvis' only child. Sole heir to her father's estate, she has developed a career in the music business and has issued three albums. Presley has been married four times, including to singer Michael Jackson and actor Nicolas Cage, before marrying music producer Michael Lockwood, father of her twin girls."},
			{"What is the capital of Iraq?", "Baghdad"},
			{"What is the capital of Iran?", "Tehran"},
			{"What is the population of Iraq?", "37056169"},
		},
		{
			{"Who married Michael Jackson?", "Which one? [dbpedia/http://dbpedia.org/resource/Mariléia_dos_Santos] person; description: Women's footballer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(footballer,_born_1980)] person; description: English footballer, born 1980 [dbpedia/http://dbpedia.org/resource/Michael_Jackson] person; description: American singer and dancer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(American_Revolution)] person; description: American revolutionary officer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(English_singer)] person; description: UK male singer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(TV_executive)] person; description: British television producer and executive [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(actor)] person; description: Canadian actor, grip and gaffer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(anthropologist)] person; description: New Zealand poet and anthropologist [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(basketball)] person; description: retired American professional basketball player [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(bishop)] person; description: Church of Ireland Archbishop of Dublin [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(journalist)] person; description: Niuean journalist and former politician [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(linebacker)] person; description: former professional American football player, born 1957 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(radio_commentator)] person; description: American talk radio host [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(rugby_league)] person; description: retired professional rugby league footballer [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(wide_receiver)] person; description: former American professional football player, born 1969 [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(writer)] person; description: English writer and journalist [dbpedia/http://dbpedia.org/resource/2000–01_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/2002–03_Swansea_City_A.F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(footballer,_born_1973)] person; description: retired English professional football defender, born 1973 [dbpedia/http://dbpedia.org/resource/1996–97_Bury_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/1996–97_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_A._Jackson] person; description: British computer scientist [dbpedia/http://dbpedia.org/resource/1999–2000_Preston_North_End_F.C._season__Michael_Jackson__1] person [dbpedia/http://dbpedia.org/resource/Michael_Jackson_(engineer)] person"},
			{"dbpedia/http://dbpedia.org/resource/Michael_Jackson", "Lisa Marie Presley and Debbie Rowe married him"},
			{"Was Michael Jackson married to Elvis Presley's daughter?", "Yes"},
			{"Was Michael Jackson married to the daughter of Elvis Presley?", "Yes"},
		},
		{
			{"When was Lord Byron born?", "Which one? [dbpedia/http://dbpedia.org/resource/Lord_Byron] person; description: English poet and a leading figure in the Romantic movement [dbpedia/http://dbpedia.org/resource/Lord_Byron_(umpire)] person; description: Major League Baseball umpire"},
			{"dbpedia/http://dbpedia.org/resource/Lord_Byron", "He was born on January 22, 1788"},
			{"Who married Lord Byron?", "Anne Isabella Byron married him"},
			{"How many children had Lord Byron?", "He has 2 children"}, // Ada and Allegra
			{"When did Lord Byron die?", "He died on April 19, 1824"},
		},
		{
			{"How many countries have population above 130000000", "8"},
		 	{"What is the largest state of America by area?", "Alaska"},
			{"What are the two largest states of america by area?", "Texas and Alaska"},
			{"What is the second largest state of america by area?", "Texas"},
		},
		{
		},
	}

	log := common.NewSystemLog()
	log.SetDebug(true)
	system := global.NewSystem(common.Dir() + "/../../resources/dbpedia", common.Dir() + "/../../var", log)
	sessionId := "dbpedia-demo"

	if !log.IsOk() {
		t.Errorf(log.String())
		return
	}

	system.RemoveDialogContext(sessionId)

	for _, session := range tests {

		system.ClearDialogContext()

		for _, test := range session {

			log.Clear()

			system.PopulateDialogContext(sessionId, false)

			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			system.StoreDialogContext(sessionId)

			if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				t.Error(log.String())
			}

		}
	}
}
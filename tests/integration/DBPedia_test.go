package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/global"
	"testing"
)

func TestDBPedia(t *testing.T) {

	var tests = [][]struct {
		question string
		answer   string
	}{
		{
			{"How many children has Madonna?", "She has 4 children"},
			{"How old is she", "She is 64 years old ( born on August 16, 1958 )"},
			//{"When is her birthday?", "August 16, 1958", "", ""},
		},
		{
			{"When is Madonna's birthday?", "August 16, 1958"},
			{"Who is Madonna's husband?", "Sean Penn and Guy Ritchie"},
			{"Who is Sean Penn's wife?", "Robin Wright and Madonna"},
			{"How old is Madonna?", "She is 64 years old ( born on August 16, 1958 )"},
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
			{"Who married Michael Jackson?", " [0] Women's footballer [1] English footballer, born 1980 [2] American singer and dancer [3] American revolutionary officer [4] UK male singer [5] British television producer and executive [6] Canadian actor, grip and gaffer [7] New Zealand poet and anthropologist [8] retired American professional basketball player [9] Church of Ireland Archbishop of Dublin [10] Niuean journalist and former politician [11] former professional American football player, born 1957 [12] American talk radio host [13] retired professional rugby league footballer [14] former American professional football player, born 1969 [15] English writer and journalist [16]  [17]  [18] retired English professional football defender, born 1973 [19]  [20]  [21] British computer scientist [22]  [23] "},
			{"2", "Lisa Marie Presley and Debbie Rowe married him"},
			{"Was Michael Jackson married to Elvis Presley's daughter?", "Yes"},
			{"Was Michael Jackson married to the daughter of Elvis Presley?", "Yes"},
		},
		{
			{"When was Lord Byron born?", " [0] English poet and a leading figure in the Romantic movement [1] Major League Baseball umpire"},
			{"0", "He was born on January 22, 1788"},
			{"Who married Lord Byron?", "Anne Isabella Byron married him"},
			{"How many children had Lord Byron?", "He has 2 children"}, // Ada and Allegra
			{"When did Lord Byron die?", "He died on April 19, 1824"},
			{"Who is Lord Byron's youngest daughter?", "Allegra Byron"},
			{"Who is Lord Byron's oldest daughter?", "Ada Lovelace"},
		},
		{
			{"How many countries have population above 130000000", "8"},
			{"What is the largest state of America by area?", "Alaska"},
			{"What are the two largest states of america by area?", "Alaska and Texas"},
			{"What is the second largest state of america by area?", "Texas"},
			{"What is the largest state of America by population?", "California"},
			{"What is america's largest state by population?", "California"},
		},
		{},
	}

	log := common.NewSystemLog()
	//log.SetDebug(true)
	//log.SetPrint(true)

	for _, session := range tests {

		system := global.NewSystem(common.Dir()+"/../../resources/dbpedia", "dbpedia-demo", common.Dir()+"/../../var", log, nil)

		if !log.IsOk() {
			t.Errorf("error...")
			return
		}

		for _, test := range session {

			log.Clear()
			println(test.question)
			answer, options := system.Answer(test.question)

			if options.HasOptions() {
				answer += options.String()
			}

			if !log.IsOk() {
				t.Errorf(log.GetErrors()[0])
				t.Errorf("\n%s", log.String())
				break
			} else if answer != test.answer {
				t.Errorf("Test relationships: got %v, want %v", answer, test.answer)
				t.Errorf("\n%s", log.String())
				break
			}

		}
	}
}

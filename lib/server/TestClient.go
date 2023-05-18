package server

import (
	"fmt"
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type TestClient struct {
	conn           *websocket.Conn
	applicationDir string
	workDir        string
}

const SESSION_ID = "test123"

func CreateTestClient(applicationDir string, workDir string) *TestClient {

	address := "ws://localhost:3333/"
	conn, err := websocket.Dial(address, "", address)
	if err != nil {
		panic("Could not connect to server: " + err.Error())
	}

	return &TestClient{
		conn:           conn,
		applicationDir: applicationDir,
		workDir:        workDir,
	}
}

func (c *TestClient) Close() {
	println("Client closed")
	c.conn.Close()
}

func (c *TestClient) Run(tests []Test) {

	var err error

	for _, test := range tests {

		println("TEST: " + test.H)

		request := mentalese.Request{
			SessionId: SESSION_ID,
			// todo: just send the application's name; this is insecure information
			ApplicationDir: c.applicationDir,
			WorkDir:        c.workDir,
			Command:        "send",
			Message: mentalese.RelationSet{
				mentalese.NewRelation(false, "go_respond", []mentalese.Term{
					mentalese.NewTermString(test.H),
				}),
			},
		}

		err = websocket.JSON.Send(c.conn, request)
		if err != nil {
			panic("Could not send to server: " + err.Error())
		}

		ok := true

		for true {

			response := mentalese.Response{}
			var err error = nil

			err = websocket.JSON.Receive(c.conn, &response)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("result: %v\n", response)

			first := response.Message[0]
			if first.Predicate == mentalese.PredicatePrint {

				c.Send(
					mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
						mentalese.NewTermRelationSet([]mentalese.Relation{first}),
					}),
				)

				answer := first.Arguments[1].TermValue
				println("      " + answer)
				if answer != test.C {
					ok = false
					println("ERROR expected " + test.C + ", got: " + answer)
				}

				// break
			}
			if first.Predicate == "dom_action_move_to" {

				println("move")
				for _, relation := range response.Message {
					c.Send(
						mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
							mentalese.NewTermRelationSet([]mentalese.Relation{relation}),
						}),
					)
				}

			}
			if first.Predicate == "go_processlist_clear" {
				println("cool! notification of all empty processes!")
				break
			}
		}

		if !ok {
			break
		}

	}
}

func (c *TestClient) Send(message mentalese.Relation) {
	request := mentalese.Request{
		SessionId:      SESSION_ID,
		ApplicationDir: c.applicationDir,
		WorkDir:        c.workDir,
		Command:        "send",
		Message:        []mentalese.Relation{message},
	}

	err := websocket.JSON.Send(c.conn, request)
	if err != nil {
		panic("Could not send to server: " + err.Error())
	}
}

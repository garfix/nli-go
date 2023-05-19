package server

import (
	"fmt"
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type TestClient struct {
	conn       *websocket.Conn
	systemName string
}

const SESSION_ID = "test123"

func CreateTestClient(systemName string) *TestClient {

	address := "ws://localhost:3333/"
	conn, err := websocket.Dial(address, "", address)
	if err != nil {
		panic("Could not connect to server: " + err.Error())
	}

	return &TestClient{
		conn:       conn,
		systemName: systemName,
	}
}

func (c *TestClient) Close() {
	println("Client closed")
	c.conn.Close()
}

func (c *TestClient) Run(tests []Test) {

	var err error

	for _, test := range tests {

		clarificationIndex := 0

		println("TEST: " + test.H)

		request := mentalese.Request{
			System:    c.systemName,
			SessionId: SESSION_ID,
			Command:   "send",
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
			if first.Predicate == "go_user_select" {
				answer := first.Copy()
				answer.Arguments[2] = mentalese.NewTermString(test.Clarifications[clarificationIndex])
				clarificationIndex++
				c.Send(
					mentalese.NewRelation(false, mentalese.PredicateAssert, []mentalese.Term{
						mentalese.NewTermRelationSet([]mentalese.Relation{answer}),
					}),
				)
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
		SessionId: SESSION_ID,
		System:    c.systemName,
		Command:   "send",
		Message:   []mentalese.Relation{message},
	}

	err := websocket.JSON.Send(c.conn, request)
	if err != nil {
		panic("Could not send to server: " + err.Error())
	}
}

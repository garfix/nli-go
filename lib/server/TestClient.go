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
			System:      c.systemName,
			MessageType: mentalese.MessageRespond,
			Message:     test.H,
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

			if response.MessageType == mentalese.PredicatePrint {

				var relations []interface{} = (response.Message).([]interface{})
				first := relations[0].(mentalese.Relation)

				c.Send(mentalese.MessageAnswer, "OK")

				answer := first.Arguments[1].TermValue
				println("      " + answer)
				if answer != test.C {
					ok = false
					println("ERROR expected " + test.C + ", got: " + answer)
				}

				// break
			}
			if response.MessageType == "move_to" {

				println("move")
				c.Send(mentalese.MessageAnswer, "OK")

			}
			if response.MessageType == "choice" {
				var relations []interface{} = (response.Message).([]interface{})
				first := relations[0].(mentalese.Relation)
				answer := first.Copy()
				answer.Arguments[2] = mentalese.NewTermString(test.Clarifications[clarificationIndex])
				clarificationIndex++
				c.Send(mentalese.MessageAnswer, test.Clarifications[clarificationIndex])
			}
			if response.MessageType == "done" {
				println("cool! notification of all empty processes!")
				break
			}
		}

		if !ok {
			break
		}

	}
}

func (c *TestClient) Send(messageType string, message string) {
	request := mentalese.Request{
		System:      c.systemName,
		MessageType: messageType,
		Message:     message,
	}

	err := websocket.JSON.Send(c.conn, request)
	if err != nil {
		panic("Could not send to server: " + err.Error())
	}
}

package server

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
)

type TestClient struct {
	conn       *websocket.Conn
	systemName string
}

const SESSION_ID = "test123"

func CreateTestClient(systemName string) *TestClient {

	address := "ws://localhost:3334/"
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

	for _, test := range tests {

		clarificationIndex := 0

		println("TEST: " + test.H)

		c.Send(central.LANGUAGE_PROCESS, mentalese.MessageRespond, test.H)

		ok := true

		for true {

			response := mentalese.Response{}
			var err error = nil

			err = websocket.JSON.Receive(c.conn, &response)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("result: %v\n", response)

			if response.MessageType == mentalese.MessagePrint {

				var answer string = (response.Message).(string)

				c.Send(central.LANGUAGE_PROCESS, mentalese.MessageAcknowledge, "")

				println("      " + answer)
				if answer != test.C {
					ok = false
					println("ERROR expected " + test.C + ", got: " + answer)
				}

				// break
			}
			if response.MessageType == "move_to" {

				println("move")
				c.Send(central.ROBOT_PROCESS, mentalese.MessageAcknowledge, "")

			}
			if response.MessageType == mentalese.MessageChoose {
				c.Send(central.LANGUAGE_PROCESS, mentalese.MessageChosen, test.Clarifications[clarificationIndex])
				clarificationIndex++
			}
			if response.MessageType == mentalese.MessageProcessListClear {
				println("cool! notification of all empty processes!")
				break
			}
		}

		if !ok {
			break
		}

	}
}

func (c *TestClient) Send(processType string, messageType string, message string) {
	request := mentalese.Request{
		System:      c.systemName,
		ProcessType: processType,
		MessageType: messageType,
		Message:     message,
	}

	err := websocket.JSON.Send(c.conn, request)
	if err != nil {
		panic("Could not send to server: " + err.Error())
	}
}

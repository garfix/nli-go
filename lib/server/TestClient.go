package server

import (
	"fmt"
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"

	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v2"
)

type TestClient struct {
	conn *websocket.Conn
}

const SESSION_ID = "test123"

func CreateTestClient() *TestClient {

	address := "ws://localhost:3334/"
	conn, err := websocket.Dial(address, "", address)
	if err != nil {
		panic("Could not connect to server: " + err.Error())
	}

	return &TestClient{
		conn: conn,
	}
}

func (c *TestClient) Close() {
	println("Client closed")
	c.conn.Close()
}

type Test struct {
	H              string   `yaml:"H"`
	C              string   `yaml:"C"`
	Clarifications []string `yaml:"Clarifications"`
	Send           string   `yaml:"Send"`
}

func (c *TestClient) RunTests(system string, filename string) {
	yml, err := common.ReadFile(filename)
	if err != nil {
		println("Error reading " + filename)
		return
	}

	tests := []Test{}

	err = yaml.Unmarshal([]byte(yml), &tests)
	if err != nil {
		println("Error parsing " + filename + ": " + err.Error())
	}

	c.Run(system, tests)
}

func (c *TestClient) Run(system string, tests []Test) {

	response := mentalese.Response{}
	var err error = nil

	c.Send(system, central.NO_RESOURCE, mentalese.MessageReset, "")

	err = websocket.JSON.Receive(c.conn, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, test := range tests {

		clarificationIndex := 0

		if test.Send != "" {

			println(test.Send)
			c.Send(system, central.RESOURCE_LANGUAGE, test.Send, "")

			err = websocket.JSON.Receive(c.conn, &response)
			if err != nil {
				fmt.Println(err)
			}

		} else {

			println(test.H)
			c.Send(system, central.RESOURCE_LANGUAGE, mentalese.MessageRespond, test.H)

			ok := true

			for {

				err = websocket.JSON.Receive(c.conn, &response)
				if err != nil {
					fmt.Println(err)
				}

				if response.MessageType == mentalese.MessagePrint {

					var answer string = (response.Message).(string)

					c.Send(system, central.RESOURCE_LANGUAGE, mentalese.MessageAcknowledge, "")

					println("  " + answer)
					if answer != test.C {
						ok = false
						println("ERROR expected \"" + test.C + "\", got: \"" + answer + "\"")
						break
					}
				}
				if response.MessageType == "move_to" {
					c.Send(system, central.RESOURCE_ROBOT, mentalese.MessageAcknowledge, "")
				}
				if response.MessageType == mentalese.MessageChoose {
					if clarificationIndex >= len(test.Clarifications) {
						ok = false
						println("Missing clarification for " + response.Message.(string))
						break
					}
					c.Send(system, central.RESOURCE_LANGUAGE, mentalese.MessageChosen, test.Clarifications[clarificationIndex])
					clarificationIndex++
				}
				if response.MessageType == mentalese.MessageProcessListClear {
					break
				}
			}

			if !ok {
				break
			}
		}

	}
}

func (c *TestClient) Send(system string, resource string, messageType string, message string) {
	request := mentalese.Request{
		System:      system,
		Resource:    resource,
		MessageType: messageType,
		Message:     message,
	}

	err := websocket.JSON.Send(c.conn, request)
	if err != nil {
		panic("Could not send to server: " + err.Error())
	}
}

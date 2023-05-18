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

		request := mentalese.Request{
			SessionId: SESSION_ID,
			// todo: just send the application's name; this is insecure information
			ApplicationDir: c.applicationDir,
			WorkDir:        c.workDir,
			Command:        "send",
			Message: mentalese.RelationSet{
				mentalese.NewRelation(false, "go_tell", []mentalese.Term{
					mentalese.NewTermString(test.H),
				}),
			},
		}

		err = websocket.JSON.Send(c.conn, request)
		if err != nil {
			panic("Could not connect to server: " + err.Error())
		}

		for true {

			response := mentalese.Response{}

			websocket.JSON.Receive(c.conn, &response)

			fmt.Printf("result: %v", response)

			break
		}

	}
}

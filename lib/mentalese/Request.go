package mentalese

type Request struct {
	System      string
	Resource    string
	MessageType string
	Message     interface{}
}

package mentalese

type Request struct {
	System      string
	ProcessType string
	MessageType string
	Message     interface{}
}

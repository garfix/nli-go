package mentalese

type Response struct {
	Resource    string
	MessageType string
	Success     bool
	ErrorLines  []string
	Productions []string
	Message     interface{}
}

type ResponseAnswer struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Answer      string
}

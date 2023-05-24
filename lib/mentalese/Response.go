package mentalese

type Response struct {
	ProcessType string
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

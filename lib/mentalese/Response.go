package mentalese

type Response struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Message     RelationSet
}

type ResponseAnswer struct {
	Success     bool
	ErrorLines  []string
	Productions []string
	Answer      string
}

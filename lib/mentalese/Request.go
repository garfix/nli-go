package mentalese

type Request struct {
	SessionId   string
	System      string
	Command     string
	Message     RelationSet
	Query       string
	ProcessType string
}

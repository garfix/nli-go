package mentalese

type Request struct {
	SessionId      string
	ApplicationDir string
	WorkDir        string
	Command        string
	Message        RelationSet
	Query          string
	ProcessType    string
}

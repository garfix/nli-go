package api

// This struct is a working environment of a single step in a stack frame
// It is needed by relations that have child stack frames:
// when a stack frame has finished, the parent relation is re-entered and continued

type ProcessCursor interface {
	GetState() string
	GetType() string
	SetType(string)
}

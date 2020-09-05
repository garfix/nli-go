package global

type Entities map[string]EntityInfo

type EntityInfo struct {
	Name string
	Knownby map[string]string
}

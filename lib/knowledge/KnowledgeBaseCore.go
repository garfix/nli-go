package knowledge

type KnowledgeBaseCore struct {
	Name string
}

func (core KnowledgeBaseCore) GetName() string {
	return core.Name
}

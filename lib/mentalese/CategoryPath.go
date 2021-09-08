package mentalese

type CategoryPath []CategoryPathNode

func (c CategoryPath) BindSingle(binding Binding) CategoryPath {
	newPath := CategoryPath{}
	for _, node := range c {
		newPath = append(newPath, node.BindSingle(binding))
	}
	return newPath
}

func (c CategoryPath) Copy() CategoryPath {
	newPath := CategoryPath{}
	for _, node :=  range c {
		newPath = append(newPath, node.Copy())
	}

	return newPath
}

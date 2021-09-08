package mentalese

type CategoryPathList []CategoryPath

func (list CategoryPathList) Copy() CategoryPathList {
	newEllipsis := CategoryPathList{}
	for _, path := range list {
		newEllipsis = append(newEllipsis, path.Copy())
	}
	return newEllipsis
}

func (list CategoryPathList) BindSingle(binding Binding) CategoryPathList {
	newList := CategoryPathList{}
	for _, path := range list {
		newList = append(newList, path.BindSingle(binding))
	}
	return newList
}
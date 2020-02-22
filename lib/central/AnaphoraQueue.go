package central

type AnaphoraQueue []EntityReferenceGroup

func (queue *AnaphoraQueue) AddReferenceGroup(entityReferenceGroup EntityReferenceGroup) {

	// empty group? ignore
	if len(entityReferenceGroup) == 0 {
		return
	}

	// same element again? ignore
	if len(*queue) > 0 && (*queue)[0].Equals(entityReferenceGroup) {
		return
	}

	// remove any existing reference
	for i := range *queue {
		if (*queue)[i].Equals(entityReferenceGroup) {
			*queue = append((*queue)[0:i], (*queue)[i+1:]...)
			break
		}
	}

	// prepend the reference
	*queue = append([]EntityReferenceGroup{entityReferenceGroup}, *queue...)

	// queue too long: remove the last element
	if len(*queue) > MaxSizeAnaphoraQueue {
		*queue = (*queue)[0:MaxSizeAnaphoraQueue]
	}
}

func (queue *AnaphoraQueue) String() string {
	list := ""
	sep := ""
	for _, group := range *queue {
		list += sep + group.String()
		sep = " "
	}
	return "[" + list + "]"
}
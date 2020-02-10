package central

type AnaphoraQueue []EntityReference

func (queue *AnaphoraQueue) AddEntityReference(entityReference EntityReference) {

	// same element again? ignore
	if len(*queue) > 0 && (*queue)[0].Equals(entityReference) {
		return
	}

	// remove any existing reference
	for i := range *queue {
		if (*queue)[i].Equals(entityReference) {
			*queue = append((*queue)[0:i], (*queue)[i+1:]...)
			break
		}
	}

	// prepend the reference
	*queue = append([]EntityReference{entityReference}, *queue...)

	// queue too long: remove the last element
	if len(*queue) > MaxSizeAnaphoraQueue {
		*queue = (*queue)[0:MaxSizeAnaphoraQueue]
	}
}

func (queue *AnaphoraQueue) String() string {
	list := ""
	sep := ""
	for _, ref := range *queue {
		list += sep + ref.String()
		sep = " "
	}
	return "[" + list + "]"
}
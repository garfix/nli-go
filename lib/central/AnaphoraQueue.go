package central

type AnaphoraQueue []EntityReferenceGroup

func NewAnaphoraQueue() *AnaphoraQueue {
	return &AnaphoraQueue{}
}

func (queue *AnaphoraQueue) Initialize() {
	(*queue) = []EntityReferenceGroup{}
}

func (queue *AnaphoraQueue) RemoveVariables() {
	for i, group := range *queue {
		for j := range group {
			(*queue)[i][j].Variable = ""
		}
	}
}

func (queue *AnaphoraQueue) AddReferenceGroup(entityReferenceGroup EntityReferenceGroup) {

	// empty group? ignore
	if len(entityReferenceGroup) == 0 {
		return
	}

	// remove doubles
	entityReferenceGroup = entityReferenceGroup.Deduplicate()

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

func (queue *AnaphoraQueue) FormattedString() string {
	str := ""

	for _, group := range *queue {
		sep := ""
		for _, item := range group {
			str += sep + item.String()
			sep = ", "
		}

		str += "\n"
	}

	return str
}
package central

type AnaphoraQueueElementValue struct {
	Sort string
	Id   string
}

type AnaphoraQueueElement struct {
	Variable string
	values   []AnaphoraQueueElementValue
}

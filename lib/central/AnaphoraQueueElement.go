package central

type AnaphoraQueueElementValue struct {
	Sort  string
	Id    string
	Score int
}

type AnaphoraQueueElement struct {
	Variable string
	values   []AnaphoraQueueElementValue
}

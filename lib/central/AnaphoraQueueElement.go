package central

type AnaphoraQueueElementValue struct {
	Sort string
	Id   string
}

type AnaphoraQueueElement struct {
	SentenceDistance int
	Variable         string
	values           []AnaphoraQueueElementValue
}

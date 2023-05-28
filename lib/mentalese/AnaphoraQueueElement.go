package mentalese

type AnaphoraQueueElementValue struct {
	Sort  string
	Id    string
	Score int
}

type AnaphoraQueueElement struct {
	Variable string
	Values   []AnaphoraQueueElementValue
}

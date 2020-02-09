package central

type AnaphoraQueue []EntityReference

//func (queue *AnaphoraQueue) FindEntityReferences(entityType string) []EntityReference {
//	foundReferences := []EntityReference{}
//
//	for _, entityReference := range queue {
//		if entityReference.EntityType == entityType {
//			foundReferences = append(foundReferences, entityReference)
//		}
//	}
//
//	return foundReferences
//}


func (queue *AnaphoraQueue) String() string {
	list := ""
	sep := ""
	for _, ref := range *queue {
		list += sep + ref.String()
		sep = " "
	}
	return "[" + list + "]"
}
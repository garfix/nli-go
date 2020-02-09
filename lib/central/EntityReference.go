package central

type EntityReference struct {
	EntityType string
	Id string
}

func CreateEntityReference(id string, entityType string) EntityReference {
	return EntityReference{
		EntityType: entityType,
		Id:         id,
	}
}

func (ref EntityReference) Equals(otherRef EntityReference) bool {
	return ref.EntityType == otherRef.EntityType && ref.Id == otherRef.Id
}

func (ref EntityReference) String() string {
	return ref.EntityType + ":" + ref.Id
}

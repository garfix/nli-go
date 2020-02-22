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

type EntityReferenceGroup []EntityReference

func (group EntityReferenceGroup) Equals(otherGroup EntityReferenceGroup) bool {
	eq := true
	if len(group) != len(otherGroup) {
		return false
	}
	for i := range group {
		if !group[i].Equals(otherGroup[i]) {
			eq = false
			break
		}
	}
	return eq
}

func (group EntityReferenceGroup) String() string {
	str := ""
	sep := ""
	for _, ref := range group {
		str += sep + ref.String()
		sep = ", "
	}
	return "[" + str + "]"
}

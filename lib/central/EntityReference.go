package central

type EntityReference struct {
	Sort string
	Id   string
}

func CreateEntityReference(id string, sort string) EntityReference {
	return EntityReference{
		Sort: sort,
		Id:   id,
	}
}

func (ref EntityReference) Equals(otherRef EntityReference) bool {
	return ref.Sort == otherRef.Sort && ref.Id == otherRef.Id
}

func (ref EntityReference) String() string {
	return ref.Sort + ":" + ref.Id
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

func (group EntityReferenceGroup) Deduplicate() EntityReferenceGroup {
	newGroup := EntityReferenceGroup{}
	for _, entity := range group {
		found := false
		for _, e := range newGroup {
			if e.Equals(entity) {
				found = true
			}
		}
		if !found {
			newGroup = append(newGroup, entity)
		}
	}
	return newGroup
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

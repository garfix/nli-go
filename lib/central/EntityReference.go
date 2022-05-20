package central

type EntityReference struct {
	Sort     string
	Id       string
	Variable string
}

type EntityReferenceGroup []EntityReference

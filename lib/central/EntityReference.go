package central

type EntityReference struct {
	EntityType string
	// database => key
	Keys map[string]string
}

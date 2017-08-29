package mentalese

type DbStats map[string]RelationStats

type RelationStats struct {
	Size int
	DistinctValues []int
}
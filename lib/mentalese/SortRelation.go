package mentalese

type SortRelation struct {
	superSort string
	subSort string
}

func NewSortRelation(superSort string, subSort string) SortRelation {
	return SortRelation{
		superSort: superSort,
		subSort: subSort,
	}
}

func (s SortRelation) GetSuperSort() string {
	return s.superSort
}

func (s SortRelation) GetSubSort() string {
	return s.subSort
}

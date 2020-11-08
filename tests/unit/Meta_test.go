package tests

import (
	"nli-go/lib/mentalese"
	"testing"
)

func TestMeta(t *testing.T) {
	var tests = []struct {
		subSort     string
		superSort      string
		result     bool
	}{
		{"cow", "mammal", true},
		{"pinguin", "mammal", false},
		{"mammal", "cow", false},
		{"cow", "animal", true},
		{"cow", "entity", true},
		{"cow", "thing", true},
		{"entity", "thing", true},
		{"thing", "entity", true},
		{"thing", "cow", false},
		{"cow", "cow", true},
		// when no sort hierarchy is defined, this still holds
		{"ape", "ape", true},
		// when no predicates are defined, this still holds
		{"", "", true},
	}

	meta := mentalese.NewMeta()

	meta.AddSubSort("mammal", "cow")
	meta.AddSubSort("mammal", "cat")
	meta.AddSubSort("animal", "mammal")
	meta.AddSubSort("living_thing", "animal")
	meta.AddSubSort("entity", "living_thing")
	meta.AddSubSort("thing", "entity")
	meta.AddSubSort("entity", "thing")

	for _, test := range tests {

		result := meta.MatchesSort(test.subSort, test.superSort)

		if result != test.result {
			t.Errorf("%v isa %v failed", test.subSort, test.superSort)
		}
	}
}

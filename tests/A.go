package tests

import (
	"testing"
)

func TestA(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{}

	for _, test := range tests {
		result := test.input
		if result != test.want {
			t.Errorf("%s: got %s, want %s", test.input, result, test.want)
		}
	}
}

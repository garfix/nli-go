package tests

import (
	"testing"
)

func TestMatcher(t *testing.T) {
	var tests = []struct {
		input string
		want string
	} {

	}

	for _, test := range tests {
		result := test.input
		if result != test.want {
			t.Errorf("%s: got %s, want %s", test.input, result, test.want)
		}
	}
}
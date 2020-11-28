package mentalese

import "sort"

type Sorts map[string]string

func (sorts Sorts) String() string {

	keys := []string{}
	for k := range sorts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	s := ""
	sep := ""
	for _, key := range keys {
		s += sep + key + ": " + sorts[key]
		sep = "; "
	}

	return s
}

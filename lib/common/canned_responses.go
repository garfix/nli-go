package common

import "fmt"

const NameNotFound = "Name not found"
const WhichOne = "Which one?"

func GetString(template string, argument string) string {

	t := "unknown"

	if template == "name_not_found" {
		t = "Name not found: %s"
	}

	return fmt.Sprintf(t, argument)
}
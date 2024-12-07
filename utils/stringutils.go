package utils

import (
	"regexp"
	"strings"
)

// To test if a string contains only alphanumeric characters
// Used to validate username ,comapny name,groupname etc
func IsValidName(test string) bool {

	s := strings.TrimSpace(test)
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)

}

func IsNotEmptySring(test string) bool {
	s := strings.TrimSpace(test)
	return s != ""
}

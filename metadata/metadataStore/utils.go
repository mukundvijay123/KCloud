package metadatastore

import (
	"regexp"
	"strings"
)

func isValidName(test string) bool {

	s := strings.TrimSpace(test)
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)

}

func isValidPasswd(test string) bool {
	s := strings.TrimSpace(test)
	return s != ""
}

package uid

import "regexp"

var pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func Valid(value string) bool {
	return pattern.MatchString(value)
}

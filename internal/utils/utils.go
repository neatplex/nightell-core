package utils

import (
	"regexp"
	"strconv"
)

func StringToID(id string) uint64 {
	r, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return r
	}
	return 0
}

func ValidateUsername(username string) bool {
	return regexp.MustCompile("^[a-z0-9_]+$").MatchString(username)
}

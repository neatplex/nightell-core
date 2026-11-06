package utils

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var usernamePattern = regexp.MustCompile("^[a-z0-9_.]+$")

func StringToID(id string, placeholder uint64) uint64 {
	r, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return placeholder
	}
	return r
}

func StringToInt(number string, max, placeholder int) int {
	r, err := strconv.Atoi(number)
	if err != nil {
		return placeholder
	}
	if r > max {
		return max
	}
	return r
}

// FileExist checks if the given file path exists or not.
func FileExist(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) || stat.IsDir() {
		return false
	}
	return true
}

// ValidateUsername accepts lowercase letters, digits, underscores and dots,
// with no leading, trailing or consecutive dots (like Instagram).
func ValidateUsername(username string) bool {
	return usernamePattern.MatchString(username) &&
		!strings.HasPrefix(username, ".") &&
		!strings.HasSuffix(username, ".") &&
		!strings.Contains(username, "..")
}

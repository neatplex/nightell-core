package user

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/models"
	"github.com/neatplex/nightell-core/internal/utils"
)

const (
	usernameMaxLength   = 20
	usernameFallback    = "user"
	usernameMaxAttempts = 10000
)

var (
	usernameInvalidChars = regexp.MustCompile("[^a-z0-9_.]+")
	usernameDots         = regexp.MustCompile(`\.{2,}`)
)

// usernameBase turns the local part of an email into a valid username,
// so the email itself (PII) is never used as the username.
func usernameBase(email string) string {
	base := strings.ToLower(email)
	if i := strings.LastIndex(base, "@"); i >= 0 {
		base = base[:i]
	}
	if i := strings.Index(base, "+"); i >= 0 {
		base = base[:i]
	}
	base = usernameInvalidChars.ReplaceAllString(base, "_")
	base = usernameDots.ReplaceAllString(base, ".")
	base = strings.Trim(base, "._")
	if len(base) > usernameMaxLength {
		base = strings.TrimRight(base[:usernameMaxLength], "._")
	}
	if base == "" {
		base = usernameFallback
	}
	return base
}

// generateUsername returns usernameBase(email), or it with a number appended
// when it is already taken.
func (s *Service) generateUsername(email string) (string, error) {
	base := usernameBase(email)
	candidate := base
	for i := 1; i <= usernameMaxAttempts; i++ {
		user, err := s.FindBy("username", candidate)
		if err != nil {
			return "", errors.WithStack(err)
		}
		if user == nil {
			return candidate, nil
		}

		suffix := strconv.Itoa(i)
		prefix := base
		if len(prefix)+len(suffix) > usernameMaxLength {
			prefix = prefix[:usernameMaxLength-len(suffix)]
		}
		candidate = prefix + suffix
	}
	return "", errors.Errorf("no free username for base %q", base)
}

// FixUsernames replaces invalid usernames (like the emails older versions
// used as usernames) with generated ones and returns how many it changed.
func (s *Service) FixUsernames() (int, error) {
	var users []*models.User
	if r := s.database.Handler().Order("id").Find(&users); r.Error != nil {
		return 0, errors.WithStack(r.Error)
	}

	fixed := 0
	for _, user := range users {
		if utils.ValidateUsername(user.Username) {
			continue
		}
		username, err := s.generateUsername(user.Email)
		if err != nil {
			return fixed, errors.Wrapf(err, "userId: %v", user.Id)
		}
		r := s.database.Handler().Model(user).Update("username", username)
		if r.Error != nil {
			return fixed, errors.Wrapf(r.Error, "userId: %v", user.Id)
		}
		fixed++
	}
	return fixed, nil
}

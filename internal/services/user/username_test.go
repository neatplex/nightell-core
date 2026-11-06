package user

import (
	"testing"

	"github.com/neatplex/nightell-core/internal/utils"
)

func TestUsernameBase(t *testing.T) {
	cases := map[string]string{
		"john@example.com":                       "john",
		"John.Doe@example.com":                   "john.doe",
		"john+nightell@gmail.com":                "john",
		"jane-doe@example.com":                   "jane_doe",
		"a..b@example.com":                       "a.b",
		".x.@example.com":                        "x",
		"-@example.com":                          "user",
		"verylongemailaddresslocalpart@mail.com": "verylongemailaddress",
		"abcdefghijklmnopqrs.tuv@mail.com":       "abcdefghijklmnopqrs",
		"élise@example.com":                      "lise",
	}
	for email, want := range cases {
		got := usernameBase(email)
		if got != want {
			t.Errorf("usernameBase(%q) = %q, want %q", email, got, want)
		}
		if !utils.ValidateUsername(got) {
			t.Errorf("usernameBase(%q) = %q is not a valid username", email, got)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	valid := []string{"john", "john.doe", "john_doe", "_john", "j0hn_.x"}
	invalid := []string{"", "John", "john@example.com", "john-doe", ".john", "john.", "john..doe", "jo hn"}
	for _, u := range valid {
		if !utils.ValidateUsername(u) {
			t.Errorf("ValidateUsername(%q) = false, want true", u)
		}
	}
	for _, u := range invalid {
		if utils.ValidateUsername(u) {
			t.Errorf("ValidateUsername(%q) = true, want false", u)
		}
	}
}

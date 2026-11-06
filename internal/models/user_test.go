package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserJSONHidesEmail(t *testing.T) {
	user := &User{Id: 7, Username: "john", Email: "john@example.com"}

	for name, v := range map[string]interface{}{
		"user": user,
		"post": &Post{Id: 1, User: user},
	} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "john@example.com") || strings.Contains(string(b), `"email"`) {
			t.Errorf("%s JSON exposes email: %s", name, b)
		}
	}
}

func TestPrivateUserJSONIncludesEmail(t *testing.T) {
	b, err := json.Marshal((&User{Id: 7, Username: "john", Email: "john@example.com"}).Private())
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]interface{}
	if err = json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["email"] != "john@example.com" || got["username"] != "john" || got["id"] != float64(7) {
		t.Errorf("unexpected private user JSON: %s", b)
	}
	if (*User)(nil).Private() != nil {
		t.Error("nil user should stay nil")
	}
}

package user

import "github.com/cockroachdb/errors"

var ErrUsernameAlreadyExist = errors.New("username already exist")
var ErrEmailAlreadyExist = errors.New("email already exist")

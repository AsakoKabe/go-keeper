package errs

import "fmt"

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrLoginAlreadyExist = fmt.Errorf("login already exist")
var ErrDataNotFound = fmt.Errorf("data not found")

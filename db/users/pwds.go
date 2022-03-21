package users

import (
	"errors"
	"github.com/trustelem/zxcvbn"
)

var SimplePasswordError = errors.New("password is too simple - try a more complex one")

// CheckPassword computes a zxcvbn password strength of a given
// string, and return Safe == true if the score is more than or
// equal to three
func CheckPassword(pwd string) (bool, error) {
	if len(pwd) < 8 {
		return false, SimplePasswordError
	}

	res := zxcvbn.PasswordStrength(pwd, nil)
	if res.Score >= 3 {
		return true, nil
	}
	return false, SimplePasswordError
}

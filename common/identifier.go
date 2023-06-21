package common

import (
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultIdentifierSpecialCharStr = "_"
)

// CheckIdentifier checks if the given string is a valid identifier
func CheckIdentifier(s string, checkLeadChar bool, specialCharStr string) error {
	if checkLeadChar {
		first := s[constant.ZeroInt]
		if !(first >= 'a' && first <= 'z' || first >= 'A' && first <= 'Z' || first == '_') {
			return errors.Errorf("identifier(%s) must start with upper case, lower case or under bar, %c is not valid", s, first)
		}
	}

	for _, c := range s {
		var hasSpecialChar bool

		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
			continue
		default:
			for _, sc := range specialCharStr {
				if c == sc {
					hasSpecialChar = true
					break
				}
			}

			if !hasSpecialChar {
				return errors.Errorf("identifier(%s) must only contain upper case, lower case, digit or under bar, %c is not valid", s, c)
			}
		}
	}

	return nil
}

// CheckIdentifierWithDefault checks if the given string is a valid identifier, if not, return the default value
func CheckIdentifierWithDefault(s string) error {
	return CheckIdentifier(s, true, DefaultIdentifierSpecialCharStr)
}
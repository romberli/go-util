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
			return errors.Errorf("identifier must start with alphabet or under bar, %s is not valid", s)
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
				return errors.Errorf("identifier must only contain alphabet, digit or under bar, %s is not valid", s)
			}
		}
	}

	return nil
}

// CheckIdentifierWithDefault checks if the given string is a valid identifier, if not, return the default value
func CheckIdentifierWithDefault(s string) error {
	return CheckIdentifier(s, true, DefaultIdentifierSpecialCharStr)
}

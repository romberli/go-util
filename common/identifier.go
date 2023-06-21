package common

import (
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

// CheckIdentifier checks if the given string is a valid identifier
func CheckIdentifier(s string) error {
	first := s[constant.ZeroInt]
	if !(first >= 'a' && first <= 'z' || first >= 'A' && first <= 'Z' || first == '_') {
		return errors.Errorf("identifier %s must start with upper case, lower case or under bar, %c is not valid", s, first)
	}

	for _, c := range s {
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_':
			continue
		default:
			return errors.Errorf("identifier %s must only contain upper case, lower case, digit or under bar, %c is not valid", s, c)
		}
	}

	return nil
}

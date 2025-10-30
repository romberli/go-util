package parser

import (
	"strings"

	"github.com/romberli/go-util/constant"
)

func GetFullFuncName(funcName string, args ...string) string {
	fullName := funcName
	if len(args) > constant.ZeroInt {
		fullName += constant.LeftParenthesisString
	}
	for _, arg := range args {
		fullName += arg + constant.CommaString
	}

	if len(args) > constant.ZeroInt {
		fullName = strings.TrimSuffix(fullName, constant.CommaString)
		fullName += constant.RightParenthesisString
	}

	return fullName
}

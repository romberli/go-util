package parser

import (
	"encoding/json"

	"github.com/romberli/go-util/constant"
)

const (
	ExpressionTypeUnknown ExpressionType = 0
	ExpressionTypeNull    ExpressionType = 1
	ExpressionTypeString  ExpressionType = 2
	ExpressionTypeFunc    ExpressionType = 3

	ExpressionTypeUnknownString = "UNKNOWN"
	ExpressionTypeNullString    = "NULL"
	ExpressionTypeStringString  = "STRING"
	ExpressionTypeFuncString    = "FUNC"
)

type ExpressionType int

func (et *ExpressionType) String() string {
	switch *et {
	case ExpressionTypeNull:
		return ExpressionTypeNullString
	case ExpressionTypeString:
		return ExpressionTypeStringString
	case ExpressionTypeFunc:
		return ExpressionTypeFuncString
	default:
		return ExpressionTypeUnknownString
	}
}

func (et *ExpressionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (et *ExpressionType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch s {
	case ExpressionTypeNullString:
		*et = ExpressionTypeNull
	case ExpressionTypeStringString:
		*et = ExpressionTypeString
	case ExpressionTypeFuncString:
		*et = ExpressionTypeFunc
	default:
		*et = ExpressionTypeUnknown
	}

	return nil
}

type Expression struct {
	ExpressionType  ExpressionType
	ExpressionValue string
}

func NewExpression(et ExpressionType, ev string) *Expression {
	return &Expression{
		et,
		ev,
	}
}

func (e *Expression) String() string {
	switch e.ExpressionType {
	case ExpressionTypeNull:
		return ExpressionTypeNullString
	case ExpressionTypeString:
		return constant.SingleQuoteString + e.ExpressionValue + constant.SingleQuoteString
	case ExpressionTypeFunc:
		return e.ExpressionValue
	default:
		return constant.EmptyString
	}
}

func (e *Expression) Equal(other *Expression) bool {
	if e == nil && other == nil {
		return true
	}

	if e != nil && other != nil {
		return e.ExpressionType == other.ExpressionType && e.ExpressionValue == other.ExpressionValue
	}

	return false
}

func (e *Expression) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

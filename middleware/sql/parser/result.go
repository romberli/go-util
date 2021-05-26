package parser

import (
	"encoding/json"

	"github.com/romberli/go-util/constant"
)

type Result struct {
	SQLType        string            `json:"sql_type"`
	DBNames        []string          `json:"db_names"`
	TableNames     []string          `json:"table_names"`
	TableComments  map[string]string `json:"table_comments"`
	ColumnNames    []string          `json:"column_names"`
	ColumnComments map[string]string `json:"column_comments"`
}

// NewResult returns a new *Result
func NewResult(sqlType string, dbNames []string, tableNames []string, tableComments map[string]string, columnNames []string, columnComments map[string]string) *Result {
	return &Result{
		SQLType:        sqlType,
		DBNames:        dbNames,
		TableNames:     tableNames,
		TableComments:  tableComments,
		ColumnNames:    columnNames,
		ColumnComments: columnComments,
	}
}

// NewEmptyResult returns an empty *Result
func NewEmptyResult() *Result {
	return &Result{
		SQLType:        constant.EmptyString,
		DBNames:        []string{},
		TableNames:     []string{},
		TableComments:  make(map[string]string),
		ColumnNames:    []string{},
		ColumnComments: make(map[string]string),
	}
}

// Marshal marshals result to json bytes
func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

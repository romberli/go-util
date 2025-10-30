package parser

import (
	"github.com/romberli/go-multierror"
)

type ColumnDefinition struct {
	TableSchema      string `json:"tableSchema"`
	TableName        string `json:"tableName"`
	ColumnName       string `json:"columnName"`
	DataType         string `json:"dataType"`
	ColumnType       string `json:"columnType"`
	DefaultValue     string `json:"defaultValue"`
	OnUpdateValue    string `json:"onUpdateValue"`
	CharacterSetName string `json:"characterSetName"`
	CollationName    string `json:"collationName"`
	IsAutoIncrement  bool   `json:"isAutoIncrement"`
	NotNull          bool   `json:"notNull"`
	OrdinalPosition  int    `json:"ordinalPosition"`
	ColumnComment    string `json:"columnComment"`

	Errors error `json:"errors,omitempty"`
}

// NewColumnDefinition returns a new *ColumnDefinition
func NewColumnDefinition(tableSchema, tableName, columnName string) *ColumnDefinition {
	return &ColumnDefinition{
		TableSchema: tableSchema,
		TableName:   tableName,
		ColumnName:  columnName,
	}
}

// NewEmptyColumnDefinition returns a new empty *ColumnDefinition
func NewEmptyColumnDefinition() *ColumnDefinition {
	return &ColumnDefinition{}
}

// Equal checks whether two ColumnDefinition objects are equal
func (cd *ColumnDefinition) Equal(other *ColumnDefinition) bool {
	if cd.ColumnName == other.ColumnName &&
		cd.DataType == other.DataType &&
		cd.DefaultValue == other.DefaultValue &&
		cd.OnUpdateValue == other.OnUpdateValue &&
		cd.CharacterSetName == other.CharacterSetName &&
		cd.CollationName == other.CollationName &&
		cd.IsAutoIncrement == other.IsAutoIncrement &&
		cd.NotNull == other.NotNull &&
		cd.OrdinalPosition == other.OrdinalPosition &&
		cd.ColumnComment == other.ColumnComment {
		return true
	}

	return false
}

// AddError adds error to ColumnDefinition
func (cd *ColumnDefinition) AddError(err error) {
	if err == nil {
		return
	}
	if cd.Errors == nil {
		cd.Errors = &multierror.Error{}
	}
	cd.Errors = multierror.Append(cd.Errors, err)
}

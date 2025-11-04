package parser

import (
	"encoding/json"

	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"

	"github.com/romberli/go-util/constant"
)

const (
	ColumnColumnString        = "COLUMN"
	ColumnCharacterSetString  = "CHARACTER SET"
	ColumnCollateString       = "COLLATE"
	ColumnNotNullString       = "NOT NULL"
	ColumnAutoIncrementString = "AUTO_INCREMENT"
	ColumnDefaultString       = "DEFAULT"
	ColumnOnUpdateString      = "ON UPDATE"
	ColumnCommentString       = "COMMENT"
	ColumnFirstString         = "FIRST"
	ColumnAfterString         = "AFTER"

	ColumnDiffTypeUnknown ColumnDiffType = 0
	ColumnDiffTypeAdd     ColumnDiffType = 1
	ColumnDiffTypeChange  ColumnDiffType = 2
	ColumnDiffTypeDrop    ColumnDiffType = 3
)

type ColumnDiffType int

func (cdt *ColumnDiffType) String() string {
	switch *cdt {
	case ColumnDiffTypeAdd:
		return AddKeyword
	case ColumnDiffTypeChange:
		return ChangeKeyWord
	case ColumnDiffTypeDrop:
		return DropKeyWord
	default:
		return UnknownKeyWord
	}
}

func (cdt *ColumnDiffType) MarshalJSON() ([]byte, error) {
	return json.Marshal(cdt.String())
}

func (cdt *ColumnDiffType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch s {
	case AddKeyword:
		*cdt = ColumnDiffTypeAdd
	case ChangeKeyWord:
		*cdt = ColumnDiffTypeChange
	case DropKeyWord:
		*cdt = ColumnDiffTypeDrop
	default:
		*cdt = ColumnDiffTypeUnknown
	}

	return nil
}

type ColumnDefinition struct {
	TableSchema      string      `json:"tableSchema"`
	TableName        string      `json:"tableName"`
	ColumnName       string      `json:"columnName"`
	DataType         string      `json:"dataType"`
	ColumnType       string      `json:"columnType"`
	DefaultValue     *Expression `json:"defaultValue"`
	OnUpdateValue    *Expression `json:"onUpdateValue"`
	CharacterSetName string      `json:"characterSetName"`
	CollationName    string      `json:"collationName"`
	IsAutoIncrement  bool        `json:"isAutoIncrement"`
	NotNull          bool        `json:"notNull"`
	OrdinalPosition  int         `json:"ordinalPosition"`
	ColumnComment    string      `json:"columnComment"`
	After            string      `json:"after"`
	IsFirst          bool        `json:"isFirst"`

	Errors *multierror.Error `json:"errors,omitempty"`
}

// NewColumnDefinition returns a new *ColumnDefinition
func NewColumnDefinition(tableSchema, tableName, columnName string) *ColumnDefinition {
	return &ColumnDefinition{
		TableSchema: tableSchema,
		TableName:   tableName,
		ColumnName:  columnName,
		Errors:      &multierror.Error{},
	}
}

// NewEmptyColumnDefinition returns a new empty *ColumnDefinition
func NewEmptyColumnDefinition() *ColumnDefinition {
	return &ColumnDefinition{
		Errors: &multierror.Error{},
	}
}

// String returns the string of ColumnDefinition
func (cd *ColumnDefinition) String() string {
	if cd == nil {
		return constant.EmptyString
	}
	s := constant.BackTickString + cd.ColumnName + constant.BackTickString
	s += constant.SpaceString + cd.ColumnType
	if cd.CharacterSetName != constant.EmptyString {
		s += constant.SpaceString + ColumnCharacterSetString + constant.SpaceString + cd.CharacterSetName
	}
	if cd.CollationName != constant.EmptyString {
		s += constant.SpaceString + ColumnCollateString + constant.SpaceString + cd.CollationName
	}
	if cd.NotNull && !cd.IsAutoIncrement {
		s += constant.SpaceString + ColumnNotNullString
	}
	if cd.IsAutoIncrement {
		s += constant.SpaceString + ColumnAutoIncrementString
	}
	if cd.DefaultValue != nil {
		s += constant.SpaceString + ColumnDefaultString + constant.SpaceString + cd.DefaultValue.String()
	}
	if cd.OnUpdateValue != nil {
		s += constant.SpaceString + ColumnOnUpdateString + constant.SpaceString + cd.OnUpdateValue.String()
	}
	if cd.ColumnComment != constant.EmptyString {
		s += constant.SpaceString + ColumnCommentString + constant.SpaceString + constant.SingleQuoteString + cd.ColumnComment + constant.SingleQuoteString
	}
	if cd.IsFirst {
		s += constant.SpaceString + ColumnFirstString
	}
	if cd.After != constant.EmptyString {
		s += constant.SpaceString + ColumnAfterString + constant.SpaceString + constant.BackTickString + cd.After + constant.BackTickString
	}

	return s
}

// Clone returns a new *ColumnDefinition
func (cd *ColumnDefinition) Clone() *ColumnDefinition {
	return &ColumnDefinition{
		TableSchema:      cd.TableSchema,
		TableName:        cd.TableName,
		ColumnName:       cd.ColumnName,
		DataType:         cd.DataType,
		ColumnType:       cd.ColumnType,
		DefaultValue:     cd.DefaultValue,
		OnUpdateValue:    cd.OnUpdateValue,
		CharacterSetName: cd.CharacterSetName,
		CollationName:    cd.CollationName,
		IsAutoIncrement:  cd.IsAutoIncrement,
		NotNull:          cd.NotNull,
		OrdinalPosition:  cd.OrdinalPosition,
		ColumnComment:    cd.ColumnComment,
		After:            cd.After,
		IsFirst:          cd.IsFirst,
		Errors:           cd.Errors,
	}
}

// Error returns the error of ColumnDefinition
func (cd *ColumnDefinition) Error() error {
	return cd.Errors.ErrorOrNil()
}

// Equal checks whether two ColumnDefinition objects are equal
func (cd *ColumnDefinition) Equal(other *ColumnDefinition) bool {
	if other == nil {
		return false
	}

	if cd.ColumnName == other.ColumnName &&
		cd.DataType == other.DataType &&
		cd.ColumnType == other.ColumnType &&
		cd.DefaultValue.Equal(other.DefaultValue) &&
		cd.OnUpdateValue.Equal(other.OnUpdateValue) &&
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

// Diff returns the difference between two column definitions
func (cd *ColumnDefinition) Diff(source *ColumnDefinition) *ColumnDiff {
	if source == nil {
		return NewColumnDiff(ColumnDiffTypeAdd, nil, cd)
	}

	if cd.Equal(source) {
		return nil
	}

	return NewColumnDiff(ColumnDiffTypeChange, source, cd)
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

type ColumnDiff struct {
	DiffType ColumnDiffType    `json:"diffType"`
	Source   *ColumnDefinition `json:"source"`
	Target   *ColumnDefinition `json:"target"`
}

func NewColumnDiff(diffType ColumnDiffType, source, target *ColumnDefinition) *ColumnDiff {
	return &ColumnDiff{
		DiffType: diffType,
		Source:   source,
		Target:   target,
	}
}

func NewEmptyColumnDiff() *ColumnDiff {
	return &ColumnDiff{}
}

// GetMigrationSQL gets the sql of column migration
func (cd *ColumnDiff) GetMigrationSQL() string {
	switch cd.DiffType {
	case ColumnDiffTypeAdd:
		return AddKeyword + constant.SpaceString + ColumnColumnString + constant.SpaceString + cd.Target.String()
	case ColumnDiffTypeChange:
		return ChangeKeyWord + constant.SpaceString + ColumnColumnString + constant.SpaceString +
			constant.BackTickString + cd.Source.ColumnName + constant.BackTickString +
			constant.SpaceString + cd.Target.String()
	case ColumnDiffTypeDrop:
		return DropKeyWord + constant.SpaceString + ColumnColumnString + constant.SpaceString +
			constant.BackTickString + cd.Source.ColumnName + constant.BackTickString
	default:
		return UnknownKeyWord
	}
}

// MarshalJSON returns the json of ColumnDiff
func (cd *ColumnDiff) MarshalJSON() ([]byte, error) {
	type Alias struct {
		DiffType ColumnDiffType `json:"diffType"`
		Source   string         `json:"source"`
		Target   string         `json:"target"`
	}

	aux := &Alias{
		DiffType: cd.DiffType,
		Source:   cd.Source.String(),
		Target:   cd.Target.String(),
	}

	jsonBytes, err := json.Marshal(aux)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return jsonBytes, nil
}

package parser

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/romberli/go-multierror"

	"github.com/romberli/go-util/constant"
)

const (
	IndexPrimaryKeyName   = "PRIMARY"
	IndexDescendingString = "DESC"
	IndexIndexString      = "INDEX"
	IndexKeyString        = "KEY"
	IndexUniqueString     = "UNIQUE"
	IndexVisibleString    = "INVISIBLE"

	IndexDiffTypeUnknown IndexDiffType = 0
	IndexDiffTypeAdd     IndexDiffType = 1
	IndexDiffTypeDrop    IndexDiffType = 2
)

type IndexDiffType int

func (idt *IndexDiffType) String() string {
	switch *idt {
	case IndexDiffTypeAdd:
		return AddKeyword
	case IndexDiffTypeDrop:
		return DropKeyWord
	default:
		return UnknownKeyWord
	}
}

func (idt *IndexDiffType) MarshalJSON() ([]byte, error) {
	return json.Marshal(idt.String())
}

func (idt *IndexDiffType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch s {
	case AddKeyword:
		*idt = IndexDiffTypeAdd
	case DropKeyWord:
		*idt = IndexDiffTypeDrop
	default:
		*idt = IndexDiffTypeUnknown
	}

	return nil
}

// Caution: for now, index definition does not include information below:
// - expressions on the indexed columns will be ignored
// - all indexes are assumed to be btree indexes
type IndexDefinition struct {
	TableSchema string `json:"tableSchema"`
	TableName   string `json:"tableName"`
	IndexName   string `json:"indexName"`

	IsPrimary bool `json:"isPrimary"`
	IsUnique  bool `json:"isUnique"`
	IsVisible bool `json:"isVisible"`

	Columns []*IndexSpec `json:"columns"`

	Errors *multierror.Error `json:"errors,omitempty"`
}

// NewIndexDefinition returns a new *IndexDefinition
func NewIndexDefinition(tableSchema, tableName, indexName string) *IndexDefinition {
	return &IndexDefinition{
		TableSchema: tableSchema,
		TableName:   tableName,
		IndexName:   indexName,
		IsVisible:   true,
		Errors:      &multierror.Error{},
	}
}

// NewEmptyIndexDefinition returns a new empty *IndexDefinition
func NewEmptyIndexDefinition() *IndexDefinition {
	return &IndexDefinition{
		Errors: &multierror.Error{},
	}
}

// String returns the string of IndexDefinition
func (id *IndexDefinition) String() string {
	if id == nil {
		return constant.EmptyString
	}
	var s string
	if id.IsPrimary {
		s = IndexPrimaryKeyName + constant.SpaceString + IndexKeyString + constant.SpaceString + constant.LeftParenthesisString
	} else {
		if id.IsUnique {
			s += IndexUniqueString + constant.SpaceString
		}
		s += IndexIndexString + constant.SpaceString + constant.BackTickString + id.IndexName + constant.BackTickString +
			constant.SpaceString + constant.LeftParenthesisString

	}

	for _, column := range id.Columns {
		s += column.String() + constant.CommaString + constant.SpaceString
	}
	s = strings.TrimSuffix(strings.TrimSpace(s), constant.CommaString)
	s += constant.RightParenthesisString
	if !id.IsVisible {
		s += constant.SpaceString + IndexVisibleString
	}

	return s
}

// Clone returns a new *IndexDefinition
func (id *IndexDefinition) Clone() *IndexDefinition {
	newColumns := make([]*IndexSpec, len(id.Columns))
	for i := constant.ZeroInt; i < len(id.Columns); i++ {
		newColumns[i] = id.Columns[i].Clone()
	}

	return &IndexDefinition{
		TableSchema: id.TableSchema,
		TableName:   id.TableName,
		IndexName:   id.IndexName,
		IsPrimary:   id.IsPrimary,
		IsUnique:    id.IsUnique,
		IsVisible:   id.IsVisible,
		Columns:     newColumns,
		Errors:      id.Errors,
	}
}

// Error returns the error of IndexDefinition
func (id *IndexDefinition) Error() error {
	return id.Errors.ErrorOrNil()
}

// Equal checks whether two IndexDefinition objects are equal
func (id *IndexDefinition) Equal(other *IndexDefinition) bool {
	if id.IndexName == other.IndexName &&
		id.IsPrimary == other.IsPrimary &&
		id.IsUnique == other.IsUnique &&
		id.IsVisible == other.IsVisible &&
		len(id.Columns) == len(other.Columns) {
		for i := constant.ZeroInt; i < len(id.Columns); i++ {
			sourceColumn := other.Columns[i]
			targetColumn := id.Columns[i]
			if !targetColumn.Equal(sourceColumn) {
				return false
			}
		}

		return true
	}

	return false
}

// Diff returns the difference between two index definitions
func (id *IndexDefinition) Diff(source *IndexDefinition) []*IndexDiff {
	if source == nil {
		diff := NewIndexDiff(IndexDiffTypeAdd, nil, id)
		return []*IndexDiff{diff}
	}

	if id.Equal(source) {
		return nil
	}

	drop := NewIndexDiff(IndexDiffTypeDrop, source, nil)
	add := NewIndexDiff(IndexDiffTypeAdd, nil, id)

	return []*IndexDiff{drop, add}
}

// AddIndexSpec adds an index spec to the index definition
func (id *IndexDefinition) AddIndexSpec(is *IndexSpec) {
	id.Columns = append(id.Columns, is)
}

// AddError adds error to IndexDefinition
func (id *IndexDefinition) AddError(err error) {
	if err == nil {
		return
	}
	if id.Errors == nil {
		id.Errors = &multierror.Error{}
	}
	id.Errors = multierror.Append(id.Errors, err)
}

// HandleOption handles the option of the index
func (id *IndexDefinition) HandleOption(option *ast.IndexOption) {
	if option != nil {
		switch option.Visibility {
		case ast.IndexVisibilityVisible:
			id.IsVisible = true
		case ast.IndexVisibilityInvisible:
			id.IsVisible = false
		default:
			err := errors.Errorf("got wrong index visibility value. indexName: %s, visibility: %d",
				id.IndexName, option.Visibility)
			id.AddError(err)
		}
	}
}

type IndexSpec struct {
	Column     *ColumnDefinition `json:"column"`
	Descending bool              `json:"descending"`
	Length     int               `json:"length"`
}

// NewIndexSpec returns a new *IndexSpec
func NewIndexSpec(cd *ColumnDefinition, descending bool, length int) *IndexSpec {
	return &IndexSpec{
		Column:     cd,
		Descending: descending,
		Length:     length,
	}
}

// String returns the string of IndexSpec
func (is *IndexSpec) String() string {
	s := constant.BackTickString + is.Column.ColumnName + constant.BackTickString
	if is.Length > constant.ZeroInt {
		s += constant.LeftParenthesisString + strconv.Itoa(is.Length) + constant.RightParenthesisString
	}
	if is.Descending {
		s += constant.SpaceString + IndexDescendingString
	}

	return s
}

// Clone() returns a new *IndexSpec
func (is *IndexSpec) Clone() *IndexSpec {
	return &IndexSpec{
		Column:     is.Column.Clone(),
		Descending: is.Descending,
		Length:     is.Length,
	}
}

// Equal checks whether two IndexSpec objects are equal
func (is *IndexSpec) Equal(other *IndexSpec) bool {
	if is.Column.Equal(other.Column) &&
		is.Descending == other.Descending &&
		is.Length == other.Length {
		return true
	}

	return false
}

type IndexDiff struct {
	DiffType IndexDiffType    `json:"diffType"`
	Source   *IndexDefinition `json:"source"`
	Target   *IndexDefinition `json:"target"`
}

func NewIndexDiff(diffType IndexDiffType, source, target *IndexDefinition) *IndexDiff {
	return &IndexDiff{
		DiffType: diffType,
		Source:   source,
		Target:   target,
	}
}

// GetMigrationSQL returns the migration sql of IndexDiff
func (id *IndexDiff) GetMigrationSQL() string {
	switch id.DiffType {
	case IndexDiffTypeAdd:
		return AddKeyword + constant.SpaceString + id.Target.String()
	case IndexDiffTypeDrop:
		return DropKeyWord + constant.SpaceString + IndexIndexString + constant.SpaceString +
			constant.BackTickString + id.Source.IndexName + constant.BackTickString
	default:
		return UnknownKeyWord
	}
}

// MarshalJSON returns the json format of IndexDiff
func (id *IndexDiff) MarshalJSON() ([]byte, error) {
	type Alias struct {
		DiffType IndexDiffType `json:"diffType"`
		Source   string        `json:"source"`
		Target   string        `json:"target"`
	}

	aux := &Alias{
		DiffType: id.DiffType,
		Source:   id.Source.String(),
		Target:   id.Target.String(),
	}

	jsonBytes, err := json.Marshal(aux)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return jsonBytes, nil
}

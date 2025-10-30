package parser

import (
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/romberli/go-multierror"

	"github.com/romberli/go-util/constant"
)

const (
	PrimaryKeyName = "PRIMARY"
)

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

// Equal checks whether two IndexSpec objects are equal
func (is *IndexSpec) Equal(other *IndexSpec) bool {
	if is.Column.Equal(other.Column) &&
		is.Descending == other.Descending &&
		is.Length == other.Length {
		return true
	}

	return false
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

	Errors error `json:"errors,omitempty"`
}

// NewIndexDefinition returns a new *IndexDefinition
func NewIndexDefinition(tableSchema, tableName, indexName string) *IndexDefinition {
	return &IndexDefinition{
		TableSchema: tableSchema,
		TableName:   tableName,
		IndexName:   indexName,
		IsVisible:   true,
	}
}

// NewEmptyIndexDefinition returns a new empty *IndexDefinition
func NewEmptyIndexDefinition() *IndexDefinition {
	return &IndexDefinition{
		Errors: &multierror.Error{},
	}
}

// AddIndexSpec adds an index spec to the index definition
func (id *IndexDefinition) AddIndexSpec(is *IndexSpec) {
	id.Columns = append(id.Columns, is)
}

// Equal checks whether two IndexDefinition objects are equal
func (id *IndexDefinition) Equal(other *IndexDefinition) bool {
	if id.IndexName == other.IndexName &&
		id.IsPrimary == other.IsPrimary &&
		id.IsUnique == other.IsUnique &&
		id.IsVisible == other.IsVisible &&
		len(id.Columns) == len(other.Columns) {
		for i := constant.ZeroInt; i < len(id.Columns); i++ {
			if !id.Columns[i].Equal(other.Columns[i]) {
				return false
			}
		}
		return true
	}

	return false
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

package parser

import (
	"encoding/json"
	"strings"

	"github.com/romberli/go-multierror"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	TableTableString     = "TABLE"
	TableRenameString    = "RENAME TO"
	TableEngineString    = "ENGINE"
	TableConvertString   = "CONVERT TO"
	TableCharsetString   = "CHARACTER SET"
	TableCollateString   = "COLLATE"
	TableRowFormatString = "ROW_FORMAT"
	TableCommentString   = "COMMENT"
	TableDropPrefix      = "--//## DANGER!: "

	TableDefaultRowFormat = "DEFAULT"

	TableDiffTypeUnknown TableDiffType = 0
	TableDiffTypeCreate  TableDiffType = 1
	TableDiffTypeAlter   TableDiffType = 2
	TableDiffTypeDrop    TableDiffType = 3
)

type TableDiffType int

func (tdt *TableDiffType) String() string {
	switch *tdt {
	case TableDiffTypeCreate:
		return CreateKeyWord
	case TableDiffTypeAlter:
		return AlterKeyWord
	case TableDiffTypeDrop:
		return DropKeyWord
	default:
		return UnknownKeyWord
	}
}

func (tdt *TableDiffType) MarshalJSON() ([]byte, error) {
	return json.Marshal(tdt.String())
}

func (tdt *TableDiffType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch s {
	case CreateKeyWord:
		*tdt = TableDiffTypeCreate
	case AlterKeyWord:
		*tdt = TableDiffTypeAlter
	case DropKeyWord:
		*tdt = TableDiffTypeDrop
	default:
		*tdt = TableDiffTypeUnknown
	}

	return nil
}

// Caution: for now, table definition does not include information below:
// - only base table, no view
// - no partition info
// - no foreign key info
// - some table options are not included, for example: auto_increment, compression, encryption, etc...
//
// - also see index.go for more limitations
type TableFullDefinition struct {
	CreateTableSQL string                       `json:"createTableSQL"`
	Table          *TableDefinition             `json:"table"`
	Columns        []*ColumnDefinition          `json:"columns"`
	ColumnMap      map[string]*ColumnDefinition `json:"-"`
	Indexes        map[string]*IndexDefinition  `json:"indexes"`
}

// NewEmptyTableFullDefinition returns a new empty *TableFullDefinition
func NewEmptyTableFullDefinition() *TableFullDefinition {
	return &TableFullDefinition{
		Table:     NewEmptyTableDefinition(),
		ColumnMap: make(map[string]*ColumnDefinition),
		Indexes:   make(map[string]*IndexDefinition),
	}
}

// GetColumnDefinition gets the column definition by column name
func (td *TableFullDefinition) GetColumnDefinition(columnName string) *ColumnDefinition {
	return td.ColumnMap[columnName]
}

// Error returns error
func (td *TableFullDefinition) Error() error {
	merr := &multierror.Error{}
	// column
	for _, column := range td.Columns {
		merr = multierror.Append(merr, column.Errors)
	}
	// index
	for _, index := range td.Indexes {
		merr = multierror.Append(merr, index.Errors)
	}

	return merr.ErrorOrNil()
}

// Equal checks whether two TableFullDefinition objects are equal
func (td *TableFullDefinition) Equal(other *TableFullDefinition) bool {
	if td.Table.Equal(other.Table) &&
		len(td.Columns) == len(other.Columns) &&
		len(td.Indexes) == len(other.Indexes) {
		for _, cd := range td.Columns {
			if !cd.Equal(td.ColumnMap[cd.ColumnName]) {
				return false
			}
		}
		for _, id := range td.Indexes {
			if !id.Equal(td.Indexes[id.IndexName]) {
				return false
			}
		}

		return true
	}

	return false
}

// Diff returns the difference between two table definitions
func (td *TableFullDefinition) Diff(source *TableFullDefinition) *TableDefinitionDiff {
	var (
		columnDiffList []*ColumnDiff
		indexDiffList  []*IndexDiff
	)
	// table
	tableDiff := td.Table.Diff(source.Table)
	var extraColumns []string
	// column
	for _, sourceColumn := range source.Columns {
		_, ok := td.ColumnMap[sourceColumn.ColumnName]
		if !ok {
			columnDiffList = append(columnDiffList, NewColumnDiff(ColumnDiffTypeDrop, sourceColumn, nil))
			source.MaintainOrdinalPosition(ColumnDiffTypeDrop, sourceColumn.OrdinalPosition)
		}
	}
	for _, targetColumn := range td.Columns {
		_, ok := source.ColumnMap[targetColumn.ColumnName]
		if !ok {
			columnDiffList = append(columnDiffList, NewColumnDiff(ColumnDiffTypeAdd, nil, targetColumn))
			source.MaintainOrdinalPosition(ColumnDiffTypeAdd, targetColumn.OrdinalPosition)
			extraColumns = append(extraColumns, targetColumn.ColumnName)
		}
	}
	for _, targetColumn := range td.Columns {
		sourceColumn := source.ColumnMap[targetColumn.ColumnName]
		if !common.ElementInSlice(extraColumns, targetColumn.ColumnName) {
			columnDiff := targetColumn.Diff(sourceColumn)
			if columnDiff != nil {
				columnDiffList = append(columnDiffList, columnDiff)
			}
		}
	}
	// index
	for _, sourceIndex := range source.Indexes {
		_, ok := td.Indexes[sourceIndex.IndexName]
		if !ok {
			indexDiffList = append(indexDiffList, NewIndexDiff(IndexDiffTypeDrop, sourceIndex, nil))
		}
	}
	for _, targetIndex := range td.Indexes {
		sourceIndex, ok := source.Indexes[targetIndex.IndexName]
		if ok {
			indexDiffs := targetIndex.Diff(sourceIndex)
			if len(indexDiffs) > constant.ZeroInt {
				indexDiffList = append(indexDiffList, indexDiffs...)
			}
		} else {
			indexDiffList = append(indexDiffList, NewIndexDiff(IndexDiffTypeAdd, nil, targetIndex))
		}
	}

	return NewTableDefinitionDiff(source.Table.GetFullTableName(), td.Table.GetFullTableName(), tableDiff, columnDiffList, indexDiffList)
}

// AddColumn adds a column to the table definition
func (td *TableFullDefinition) AddColumn(cd *ColumnDefinition) {
	td.Columns = append(td.Columns, cd)
	td.ColumnMap[cd.ColumnName] = cd
}

// AddIndex adds an index to the table definition
func (td *TableFullDefinition) AddIndex(id *IndexDefinition) {
	td.Indexes[id.IndexName] = id
}

// MaintainOrdinalPosition maintains the ordinal position of the columns
func (td *TableFullDefinition) MaintainOrdinalPosition(diffType ColumnDiffType, ordinalPosition int) {
	switch diffType {
	case ColumnDiffTypeAdd:
		for _, column := range td.Columns {
			if column.OrdinalPosition >= ordinalPosition {
				column.OrdinalPosition++
			}
		}
	case ColumnDiffTypeDrop:
		for _, column := range td.Columns {
			if column.OrdinalPosition > ordinalPosition {
				column.OrdinalPosition--
			}
		}
	default:
		return
	}
}

type TableDefinitionDiff struct {
	Source     string        `json:"source"`
	Target     string        `json:"target"`
	TableDiff  *TableDiff    `json:"tableDiff"`
	ColumnDiff []*ColumnDiff `json:"columnDiff"`
	IndexDiff  []*IndexDiff  `json:"indexDiff"`
}

// NewTableDefinitionDiff returns a new *TableDefinitionDiff
func NewTableDefinitionDiff(source, target string, tableDiff *TableDiff,
	columnDiff []*ColumnDiff, indexDiff []*IndexDiff) *TableDefinitionDiff {
	return &TableDefinitionDiff{
		Source:     source,
		Target:     target,
		TableDiff:  tableDiff,
		ColumnDiff: columnDiff,
		IndexDiff:  indexDiff,
	}
}

// NewEmptyTableDefinitionDiff returns a new empty *TableDefinitionDiff
func NewEmptyTableDefinitionDiff() *TableDefinitionDiff {
	return &TableDefinitionDiff{}
}

// GetTableMigrationSQL returns the migration sql of TableDefinitionDiff
func (tdd *TableDefinitionDiff) GetTableMigrationSQL() string {
	var sql string
	// table
	if tdd.TableDiff != nil {
		sql += tdd.TableDiff.GetTableMigrationSQL()
		if tdd.TableDiff.DiffType == TableDiffTypeCreate || tdd.TableDiff.DiffType == TableDiffTypeDrop {
			return sql
		}
		sql += constant.CommaString + constant.SpaceString
	}
	// column
	if len(tdd.ColumnDiff) > constant.ZeroInt || len(tdd.IndexDiff) > constant.ZeroInt {
		if tdd.TableDiff == nil {
			sql = AlterKeyWord + constant.SpaceString + TableTableString + constant.SpaceString +
				tdd.Target + constant.SpaceString
		}

		// drop index
		for _, id := range tdd.IndexDiff {
			if id.DiffType == IndexDiffTypeDrop {
				sql += id.GetMigrationSQL() + constant.CommaString + constant.SpaceString
			}
		}
		// drop column
		for _, cd := range tdd.ColumnDiff {
			if cd.DiffType == ColumnDiffTypeDrop {
				sql += cd.GetMigrationSQL() + constant.CommaString + constant.SpaceString
			}
		}
		// add column
		for _, cd := range tdd.ColumnDiff {
			if cd.DiffType == ColumnDiffTypeAdd {
				sql += cd.GetMigrationSQL() + constant.CommaString + constant.SpaceString
			}
		}
		// change column
		for _, cd := range tdd.ColumnDiff {
			if cd.DiffType == ColumnDiffTypeChange {
				sql += cd.GetMigrationSQL() + constant.CommaString + constant.SpaceString
			}
		}
		// add index
		for _, id := range tdd.IndexDiff {
			if id.DiffType == IndexDiffTypeAdd {
				sql += id.GetMigrationSQL() + constant.CommaString + constant.SpaceString
			}
		}

		sql = strings.TrimSuffix(strings.TrimSpace(sql), constant.CommaString) + constant.SemicolonString
	}

	return sql
}

type TableDiff struct {
	DiffType TableDiffType    `json:"diffType"`
	Source   *TableDefinition `json:"source,omitempty"`
	Target   *TableDefinition `json:"target,omitempty"`
}

// NewTableDiff returns a new *TableDiff
func NewTableDiff(diffType TableDiffType, source, target *TableDefinition) *TableDiff {
	return &TableDiff{
		DiffType: diffType,
		Source:   source,
		Target:   target,
	}
}

// NewEmptyTableDiff returns a new empty *TableDiff
func NewEmptyTableDiff() *TableDiff {
	return &TableDiff{
		Source: NewEmptyTableDefinition(),
		Target: NewEmptyTableDefinition(),
	}
}

// GetTableMigrationSQL returns the migration sql of TableDiff
func (td *TableDiff) GetTableMigrationSQL() string {
	if td == nil {
		return constant.EmptyString
	}

	switch td.DiffType {
	case TableDiffTypeCreate:
		return strings.TrimSuffix(strings.TrimSpace(td.Target.CreateTableSQL), constant.SemicolonString) + constant.SemicolonString
	case TableDiffTypeAlter:
		sql := AlterKeyWord + constant.SpaceString + TableTableString + constant.SpaceString +
			td.Source.GetFullTableName() + constant.SpaceString
		if td.Source.TableSchema != td.Target.TableSchema || td.Source.TableName != td.Target.TableName {
			sql += TableRenameString + td.Target.GetFullTableName() + constant.CommaString + constant.SpaceString
		}
		if td.Source.Charset != td.Target.Charset {
			sql += TableCharsetString + constant.SpaceString +
				td.Target.Charset + constant.CommaString + constant.SpaceString
		}
		if td.Source.Collation != td.Target.Collation {
			if td.Source.Charset == td.Target.Charset {
				sql += TableCharsetString + constant.SpaceString +
					td.Target.Charset + constant.CommaString + constant.SpaceString
			}
			sql = strings.TrimSuffix(strings.TrimSpace(sql), constant.CommaString) + constant.SpaceString +
				TableCollateString + constant.SpaceString + td.Target.Collation + constant.CommaString + constant.SpaceString
		}
		if td.Source.RowFormat != td.Target.RowFormat {
			sql += TableRowFormatString + constant.EqualString + td.Target.RowFormat + constant.CommaString + constant.SpaceString
		}
		if td.Source.TableComment != td.Target.TableComment {
			sql += TableCommentString + constant.EqualString + constant.SingleQuoteString + td.Target.TableComment + constant.SingleQuoteString
		}

		return strings.TrimSuffix(strings.TrimSpace(sql), constant.CommaString)
	case TableDiffTypeDrop:
		return TableDropPrefix + td.DiffType.String() + constant.SpaceString + TableTableString + constant.SpaceString +
			td.Source.GetFullTableName() + constant.SemicolonString
	default:
		return constant.EmptyString
	}
}

type TableDefinition struct {
	CreateTableSQL string `json:"-"`
	TableSchema    string `json:"tableSchema,omitempty"`
	TableName      string `json:"tableName,omitempty"`
	TableEngine    string `json:"tableEngine,omitempty"`
	Charset        string `json:"charset,omitempty"`
	Collation      string `json:"collation,omitempty"`
	RowFormat      string `json:"rowFormat,omitempty"`
	TableComment   string `json:"tableComment,omitempty"`
}

// NewTableDefinition returns a new *TableDefinition
func NewTableDefinition(sql, tableName, tableEngine, charset, collation, rowFormat, tableComment string) *TableDefinition {
	return &TableDefinition{
		CreateTableSQL: sql,
		TableName:      tableName,
		TableEngine:    tableEngine,
		Charset:        charset,
		Collation:      collation,
		RowFormat:      rowFormat,
		TableComment:   tableComment,
	}
}

// NewEmptyTableDefinition returns a new empty *TableDefinition
func NewEmptyTableDefinition() *TableDefinition {
	return &TableDefinition{RowFormat: TableDefaultRowFormat}
}

// GetFullTableName gets the full table name
func (td *TableDefinition) GetFullTableName() string {
	var tableSchema string
	if td.TableSchema != constant.EmptyString {
		tableSchema = constant.BackTickString + td.TableSchema + constant.BackTickString + constant.DotString
	}

	return tableSchema + td.GetTableName()
}

// GetTableName gets the table name
func (td *TableDefinition) GetTableName() string {
	return constant.BackTickString + td.TableName + constant.BackTickString
}

// Clone returns a new *TableDefinition
func (td *TableDefinition) Clone() *TableDefinition {
	return &TableDefinition{
		CreateTableSQL: td.CreateTableSQL,
		TableSchema:    td.TableSchema,
		TableName:      td.TableName,
		TableEngine:    td.TableEngine,
		Charset:        td.Charset,
		Collation:      td.Collation,
		RowFormat:      td.RowFormat,
		TableComment:   td.TableComment,
	}
}

// String returns the string of TableDefinition
func (td *TableDefinition) String() string {
	var s string
	if td.TableEngine != constant.EmptyString {
		s += TableEngineString + constant.EqualString + td.TableEngine + constant.SpaceString
	}
	if td.Charset != constant.EmptyString {
		s += TableCharsetString + constant.SpaceString + td.Charset + constant.SpaceString
	}
	if td.Collation != constant.EmptyString {
		s += TableCollateString + constant.SpaceString + td.Collation + constant.SpaceString
	}
	if td.RowFormat != constant.EmptyString {
		s += TableRowFormatString + constant.EqualString + td.RowFormat + constant.SpaceString
	}
	if td.TableComment != constant.EmptyString {
		s += TableCommentString + constant.EqualString + constant.SingleQuoteString + td.TableComment + constant.SingleQuoteString
	}

	return strings.TrimSpace(s)
}

// Equal checks whether two TableDefinition objects are equal
func (td *TableDefinition) Equal(other *TableDefinition) bool {
	if td.TableName == other.TableName &&
		td.TableEngine == other.TableEngine &&
		td.Charset == other.Charset &&
		td.Collation == other.Collation &&
		td.RowFormat == other.RowFormat &&
		td.TableComment == other.TableComment {
		return true
	}

	return false
}

// Diff returns the difference between two table definitions
func (td *TableDefinition) Diff(source *TableDefinition) *TableDiff {
	if td == nil && source != nil {
		return NewTableDiff(TableDiffTypeDrop, source, nil)
	}
	if td != nil && source == nil {
		return NewTableDiff(TableDiffTypeCreate, nil, td)
	}

	if (td == nil && source == nil) || td.Equal(source) {
		return nil
	}

	diff := NewEmptyTableDiff()
	diff.DiffType = TableDiffTypeAlter

	diff.Source.TableSchema = source.TableSchema
	diff.Target.TableSchema = td.TableSchema
	diff.Source.TableName = source.TableName
	diff.Target.TableName = td.TableName
	if td.TableEngine != source.TableEngine {
		diff.Source.TableEngine = source.TableEngine
		diff.Target.TableEngine = td.TableEngine
	}
	if td.Charset != source.Charset {
		diff.Source.Charset = source.Charset
		diff.Target.Charset = td.Charset
	}
	if td.Collation != source.Collation {
		diff.Source.Charset = source.Charset
		diff.Target.Charset = td.Charset
		diff.Source.Collation = source.Collation
		diff.Target.Collation = td.Collation
	}
	if td.RowFormat != source.RowFormat {
		diff.Source.RowFormat = source.RowFormat
		diff.Target.RowFormat = td.RowFormat
	}
	if td.TableComment != source.TableComment {
		diff.Source.TableComment = source.TableComment
		diff.Target.TableComment = td.TableComment
	}

	return diff
}

package parser

// Caution: for now, table definition does not include information below:
// - only base table, no view
// - no partition info
// - no foreign key info
// - some table options are not included, for example: auto_increment, compression, encryption, etc...
//
// - also see index.go for more limitations
type TableDefinition struct {
	TableSchema  string `json:"tableSchema"`
	TableName    string `json:"tableName"`
	TableEngine  string `json:"tableEngine"`
	Charset      string `json:"charset"`
	Collation    string `json:"collation"`
	TableComment string `json:"tableComment"`
	RowFormat    string `json:"rowFormat"`

	Columns   []*ColumnDefinition          `json:"columns"`
	ColumnMap map[string]*ColumnDefinition `json:"-"`
	Indexes   map[string]*IndexDefinition  `json:"indexes"`
}

// NewEmptyTableDefinition returns a new empty *TableDefinition
func NewEmptyTableDefinition() *TableDefinition {
	return &TableDefinition{
		ColumnMap: make(map[string]*ColumnDefinition),
		Indexes:   make(map[string]*IndexDefinition),
	}
}

// GetColumnDefinition gets the column definition by column name
func (td *TableDefinition) GetColumnDefinition(columnName string) *ColumnDefinition {
	return td.ColumnMap[columnName]
}

// AddColumn adds a column to the table definition
func (td *TableDefinition) AddColumn(cd *ColumnDefinition) {
	td.Columns = append(td.Columns, cd)
	td.ColumnMap[cd.ColumnName] = cd
}

// AddIndex adds an index to the table definition
func (td *TableDefinition) AddIndex(id *IndexDefinition) {
	td.Indexes[id.IndexName] = id
}

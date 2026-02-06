package parser

import (
	"fmt"

	"github.com/pingcap/tidb/pkg/parser/mysql"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	SQLTypeCreateUserStmt  = "CreateUserStmt"
	SQLTypeAlterUserStmt   = "AlterUserStmt"
	SQLTypeCreateTableStmt = "CreateTableStmt"
	SQLTypeAlterTableStmt  = "AlterTableStmt"
	SQLTypeDropTableStmt   = "DropTableStmt"
	SQLTypeCreateIndexStmt = "CreateIndexStmt"
	SQLTypeDropIndexStmt   = "DropIndexStmt"
	SQLTypeSelectStmt      = "SelectStmt"
	SQLTypeInsertStmt      = "InsertStmt"
	SQLTypeReplaceStmt     = "ReplaceStmt"
	SQLTypeUpdateStmt      = "UpdateStmt"
	SQLTypeDeleteStmt      = "DeleteStmt"
	SQLTypeGrantStmt       = "GrantStmt"
)

type User struct {
	User string `json:"user"`
	Host string `json:"host"`
	Pass string `json:"pass"`
}

// NewUser returns a new *User
func NewUser(user, host, pass string) *User {
	return &User{
		User: user,
		Host: host,
		Pass: pass,
	}
}

// NewEmptyUser returns an empty *User
func NewEmptyUser() *User {
	return &User{}
}

// String returns the string of the user
func (u *User) String() string {
	return fmt.Sprintf("%s@'%s'", u.User, u.Host)
}

type Result struct {
	SQLType        string                       `json:"sql_type"`
	TableDBListMap map[string][]string          `json:"table_db_list_map"`
	DBNames        []string                     `json:"db_names"`
	TableNames     []string                     `json:"table_names"`
	TableComments  map[string]string            `json:"table_comments"`
	ColumnNames    []string                     `json:"column_names"`
	ColumnTypes    map[string]string            `json:"column_types"`
	ColumnComments map[string]string            `json:"column_comments"`
	User           *User                        `json:"user"`
	Privileges     map[mysql.PrivilegeType]bool `json:"privileges"`
}

// NewResult returns a new *Result
func NewResult(sqlType string, TableDBListMap map[string][]string, dbNames []string, tableNames []string,
	tableComments map[string]string, columnNames []string, columnTypes map[string]string,
	columnComments map[string]string, user *User, privileges map[mysql.PrivilegeType]bool) *Result {
	return &Result{
		SQLType:        sqlType,
		TableDBListMap: TableDBListMap,
		DBNames:        dbNames,
		TableNames:     tableNames,
		TableComments:  tableComments,
		ColumnNames:    columnNames,
		ColumnTypes:    columnTypes,
		ColumnComments: columnComments,
		User:           user,
		Privileges:     privileges,
	}
}

// NewEmptyResult returns an empty *Result
func NewEmptyResult() *Result {
	return &Result{
		SQLType:        constant.EmptyString,
		TableDBListMap: make(map[string][]string),
		DBNames:        []string{},
		TableNames:     []string{},
		TableComments:  make(map[string]string),
		ColumnNames:    []string{},
		ColumnTypes:    make(map[string]string),
		ColumnComments: make(map[string]string),
		User:           NewEmptyUser(),
		Privileges:     make(map[mysql.PrivilegeType]bool),
	}
}

// GetSQLType returns the sql type
func (r *Result) GetSQLType() string {
	return r.SQLType
}

// GetTableDBListMap returns table db list map
func (r *Result) GetTableDBListMap() map[string][]string {
	return r.TableDBListMap
}

// GetDBNames returns the db names
func (r *Result) GetDBNames() []string {
	return r.DBNames
}

// GetTableNames returns the table names
func (r *Result) GetTableNames() []string {
	return r.TableNames
}

// GetTableComments returns the table comments
func (r *Result) GetTableComments() map[string]string {
	return r.TableComments
}

// GetColumnNames returns the column names
func (r *Result) GetColumnNames() []string {
	return r.ColumnNames
}

// GetColumnTypes returns the column types
func (r *Result) GetColumnTypes() map[string]string {
	return r.ColumnTypes
}

// GetColumnComments returns the column comments
func (r *Result) GetColumnComments() map[string]string {
	return r.ColumnComments
}

// SetSQLType sets the sql type
func (r *Result) SetSQLType(sqlType string) {
	r.SQLType = sqlType
}

// AddTableDBListMap adds db name to the result
func (r *Result) AddTableDBListMap(tableName string, dbName string) {
	dbList, ok := r.TableDBListMap[tableName]
	if !ok {
		r.TableDBListMap[tableName] = []string{dbName}
	}
	if ok && !common.ElementInSlice(dbList, dbName) {
		r.TableDBListMap[tableName] = append(dbList, dbName)
	}
}

// AddDBName adds db name to the result
func (r *Result) AddDBName(dbName string) {
	if !common.ElementInSlice(r.DBNames, dbName) {
		r.DBNames = append(r.DBNames, dbName)
	}
}

// AddTableName adds table name to the result
func (r *Result) AddTableName(tableName string) {
	if !common.ElementInSlice(r.TableNames, tableName) {
		r.TableNames = append(r.TableNames, tableName)
	}
}

// SetTableComment sets table comment of corresponding table
func (r *Result) SetTableComment(tableName string, tableComment string) {
	r.TableComments[tableName] = tableComment
}

// AddColumn adds column name to the result
func (r *Result) AddColumn(columnName string) {
	if !common.ElementInSlice(r.ColumnNames, columnName) {
		r.ColumnNames = append(r.ColumnNames, columnName)
	}
}

// SetColumnType sets column type of corresponding column
func (r *Result) SetColumnType(columnName string, columnType string) {
	r.ColumnTypes[columnName] = columnType
}

// SetColumnComment sets column comment of corresponding column
func (r *Result) SetColumnComment(columnName string, columnComment string) {
	r.ColumnComments[columnName] = columnComment
}

// SetUser sets user of the result
func (r *Result) SetUser(user *User) {
	r.User = user
}

// AddPrivilege adds privilege to the result
func (r *Result) AddPrivilege(privilege mysql.PrivilegeType, withGrant bool) {
	r.Privileges[privilege] = withGrant
}

package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	SQLTypeInsert = "INSERT"
	SQLTypeUpdate = "UPDATE"
	SQLTypeDelete = "DELETE"

	andString = "AND"
)

type Message struct {
	SQLType   string                   `json:"type"`
	IsDDL     bool                     `json:"isDdl"`
	DBName    string                   `json:"database"`
	TableName string                   `json:"table"`
	PKNames   []string                 `json:"pkNames"`
	Columns   map[string]string        `json:"mysqlType"`
	Data      []map[string]interface{} `json:"data"`
	Old       []map[string]interface{} `json:"old"`
}

// NewMessage returns a new *Message
func NewMessage(sqlType string, isDDL bool, dbName string, tableName string, pkNames []string,
	columns map[string]string, data, old []map[string]interface{}) *Message {
	return newMessage(sqlType, isDDL, dbName, tableName, pkNames, columns, data, old)
}

// newMessage returns a new *Message
func newMessage(sqlType string, isDDL bool, dbName string, tableName string, pkNames []string,
	columns map[string]string, data, old []map[string]interface{}) *Message {
	return &Message{
		SQLType:   sqlType,
		IsDDL:     isDDL,
		DBName:    dbName,
		TableName: tableName,
		PKNames:   pkNames,
		Columns:   columns,
		Data:      data,
		Old:       old,
	}
}

// NewEmptyMessage returns a new empty *Message
func NewEmptyMessage() *Message {
	return &Message{}
}

// GetSQLType returns the SQLType
func (m *Message) GetSQLType() string {
	return m.SQLType
}

// GetIsDDL returns the IsDDL
func (m *Message) GetIsDDL() bool {
	return m.IsDDL
}

// GetDBName returns the DBName
func (m *Message) GetDBName() string {
	return m.DBName
}

// GetTableName returns the TableName
func (m *Message) GetTableName() string {
	return m.TableName

}

// GetPKNames returns the PKNames
func (m *Message) GetPKNames() []string {
	return m.PKNames
}

// GetColumns returns the Columns
func (m *Message) GetColumns() map[string]string {
	return m.Columns
}

// GetData returns the Data
func (m *Message) GetData() []map[string]interface{} {
	return m.Data
}

// GetOld returns the Old
func (m *Message) GetOld() []map[string]interface{} {
	return m.Old
}

// GetColumnNames gets the column names
func (m *Message) GetColumnNames() []string {
	var columnNames []string
	for name := range m.GetColumns() {
		columnNames = append(columnNames, name)
	}

	return columnNames
}

// ConvertToSQL returns a map of the sql statement and the values
// if ignoreDDL is true, and sql type is ddl, it will only return nil, nil,
// if ignoreDDL is false, and sql type is ddl, it will return nil, error
// if useReplace is true, message with insert or update type will be converted to statements like "replace into ... values ..."
// if useReplace is false, message with insert type will be converted to statements like "insert into ... values ... on duplicate key update ..."
// if useReplace is false, message with update type will be converted to statements like "update ... set ... where ..."
func (m *Message) ConvertToSQL(ignoreDDL bool, useReplace bool) ([]map[string][]interface{}, error) {
	if len(m.GetData()) == constant.ZeroInt {
		return nil, errors.New("data should not be empty")
	}

	if m.GetIsDDL() {
		if ignoreDDL {
			return nil, nil
		}

		return nil, errors.New("does not support ddl statement, either ignore this or avoid to convert ddl statement")
	}

	switch m.GetSQLType() {
	case SQLTypeInsert:
		if useReplace {
			return m.convertToReplaceSQL()
		}
		return m.convertToInsertSQL()
	case SQLTypeUpdate:
		if useReplace {
			return m.convertToReplaceSQL()
		}
		return m.convertToUpdateSQL()
	case SQLTypeDelete:
		return m.convertToDeleteSQL()
	default:
		return nil, errors.Errorf("sql type must be one of [INSERT, UPDATE, DELETE], %s is not supported", m.GetSQLType())
	}
}

func (m *Message) convertToReplaceSQL() ([]map[string][]interface{}, error) {
	var (
		columnNamesStr string
		valuesStr      string
		values         []interface{}
	)

	// use column slice to determine the column order
	columnNames := m.GetColumnNames()

	for _, columnName := range columnNames {
		columnNamesStr += columnName + constant.CommaString
	}

	columnNamesStr = strings.TrimSuffix(columnNamesStr, constant.CommaString)

	value := constant.LeftParenthesisString
	value += strings.Repeat(constant.QuestionMarkString+constant.CommaString, len(columnNames))
	value = strings.TrimSuffix(value, constant.CommaString)
	value += constant.RightParenthesisString + constant.CommaString
	valuesStr = strings.Repeat(value, len(m.GetData()))
	valuesStr = strings.TrimSuffix(valuesStr, constant.CommaString)

	sql := `REPLACE INTO %s.%s(%s) VALUES %s ;`
	sql = fmt.Sprintf(sql, m.GetDBName(), m.GetTableName(), columnNamesStr, valuesStr)

	for _, data := range m.GetData() {
		for _, columnName := range columnNames {
			values = append(values, data[columnName])
		}
	}

	return []map[string][]interface{}{{sql: values}}, nil
}

// convertToInsertSQL converts message to a insert sql statement
func (m *Message) convertToInsertSQL() ([]map[string][]interface{}, error) {
	var (
		columnNamesStr string
		valuesStr      string
		duplicateStr   string
		values         []interface{}
	)

	// use column slice to determine the column order
	columnNames := m.GetColumnNames()

	for _, columnName := range columnNames {
		columnNamesStr += columnName + constant.CommaString
		if !common.StringInSlice(m.GetPKNames(), columnName) {
			duplicateStr += fmt.Sprintf("%s=VALUES(%s),", columnName, columnName)
		}
	}

	columnNamesStr = strings.TrimSuffix(columnNamesStr, constant.CommaString)
	duplicateStr = strings.TrimSuffix(duplicateStr, constant.CommaString)

	value := constant.LeftParenthesisString
	value += strings.Repeat(constant.QuestionMarkString+constant.CommaString, len(columnNames))
	value = strings.TrimSuffix(value, constant.CommaString)
	value += constant.RightParenthesisString + constant.CommaString
	valuesStr = strings.Repeat(value, len(m.GetData()))
	valuesStr = strings.TrimSuffix(valuesStr, constant.CommaString)

	sql := `INSERT INTO %s.%s(%s) VALUES %s ON DUPLICATE KEY UPDATE %s ;`
	sql = fmt.Sprintf(sql, m.GetDBName(), m.GetTableName(), columnNamesStr, valuesStr, duplicateStr)

	for _, data := range m.GetData() {
		for _, columnName := range columnNames {
			values = append(values, data[columnName])
		}
	}

	return []map[string][]interface{}{{sql: values}}, nil
}

// convertToUpdateSQL converts message to a update sql statement
func (m *Message) convertToUpdateSQL() ([]map[string][]interface{}, error) {
	lenData := len(m.GetData())
	lenOld := len(m.GetOld())
	if len(m.GetData()) != len(m.GetOld()) {
		return nil, errors.Errorf("the lengths of the data and old are not the same. table: %s, data: %d, old: %d", m.GetTableName(), lenData, lenOld)
	}

	if len(m.GetPKNames()) == constant.ZeroInt {
		return nil, errors.Errorf("table does not have a primary key. table: %s", m.GetTableName())
	}

	statements := make([]map[string][]interface{}, lenData)

	for i, old := range m.GetOld() {
		var values []interface{}

		sql := `UPDATE %s.%s SET %s WHERE%s;`
		setStr := constant.EmptyString
		whereStr := constant.EmptyString

		data := m.GetData()[i]
		for columnName := range old {
			values = append(values, data[columnName])
			setStr += columnName + constant.EqualString + constant.QuestionMarkString + constant.CommaString
		}

		setStr = strings.TrimSuffix(setStr, constant.CommaString)

		for _, pkName := range m.GetPKNames() {
			pkValue, ok := old[pkName]
			if !ok {
				pkValue = data[pkName]
			}

			values = append(values, pkValue)
			whereStr += fmt.Sprintf(" %s=? AND", pkName)
		}

		whereStr = strings.TrimSuffix(whereStr, andString)
		sql = fmt.Sprintf(sql, m.GetDBName(), m.GetTableName(), setStr, whereStr)

		statements[i] = map[string][]interface{}{sql: values}
	}

	return statements, nil
}

func (m *Message) convertToDeleteSQL() ([]map[string][]interface{}, error) {
	// todo: implement
	return nil, errors.New("does not support delete statement")
}

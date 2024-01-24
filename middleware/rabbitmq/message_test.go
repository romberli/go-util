package rabbitmq

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
)

const (
	testMessageIsDDL            = false
	testMessageDBName           = "test_db"
	testMessageTableName        = "test_table"
	testMessageInsertJSONString = `
        {
            "data": [
                {
                    "pk1": 88888888888888888888,
                    "pk2": "a",
                    "col1": "test_new_col1_a",
                    "col2": "test_new_col2_a"
                },
                {
                    "pk1": 2,
                    "pk2": "b",
                    "col1": "test_new_col1_b",
                    "col2": "test_new_col2_b"
                }
            ],
            "database": "test_db",
            "es": 1750061771504201728,
            "id": 1383,
            "isDdl": false,
            "mysqlType": {
                "pk1": "int",
                "pk2": "varchar(100)",
                "col1": "varchar(100)",
                "col2": "varchar(100)"
            },
            "old": null,
            "pkNames": [
                "pk1",
                "pk2"
            ],
            "sql": "",
            "sqlType": {
                "pk1": 12,
                "pk2": 12,
                "col1": 12,
                "col2": 12
            },
            "table": "test_table",
            "ts": 1656315780000,
            "type": "INSERT"
        }
    `
	testMessageUpdateJSONString = `
        {
            "data": [
                {
                    "pk1": 1,
                    "pk2": "a",
                    "col1": "test_update_col1_a",
                    "col2": "test_update_col2_a"
                },
                {
                    "pk1": 2,
                    "pk2": "b",
                    "col1": "test_update_col1_b",
                    "col2": "test_update_col2_b"
                }
            ],
            "database": "test_db",
            "es": 1656315780000,
            "id": 1383,
            "isDdl": false,
            "mysqlType": {
                "pk1": "int",
                "pk2": "varchar(100)",
                "col1": "varchar(100)",
                "col2": "varchar(100)"
            },
            "old": [
                {
                    "col1": "test_old_col1_a"
                },
                {
                    "col2": "test_old_col2_b"
                }
            ],
            "pkNames": [
                "pk1",
                "pk2"
            ],
            "sql": "",
            "sqlType": {
                "pk1": 12,
                "pk2": 12,
                "col1": 12,
                "col2": 12
            },
            "table": "test_table",
            "ts": 1656315780000,
            "type": "UPDATE"
        }
    `
	testMessageDeleteJSONString = `
        {
            "data": null,
            "database": "test_db",
            "es": 1656315780000,
            "id": 1383,
            "isDdl": false,
            "mysqlType": {
                "pk1": "int",
                "pk2": "varchar(100)",
                "col1": "varchar(100)",
                "col2": "varchar(100)"
            },
            "old": [
                {
                    "pk1": 1,
                    "pk2": "a",
                    "col1": "test_old_col1_a",
                    "col2": "test_old_col2_a"
                },
                {
                    "pk1": 2,
                    "pk2": "b",
                    "col1": "test_old_col1_b",
                    "col2": "test_old_col2_b"
                }
            ],
            "pkNames": [
                "pk1",
                "pk2"
            ],
            "sql": "",
            "sqlType": {
                "pk1": 12,
                "pk2": 12,
                "col1": 12,
                "col2": 12
            },
            "table": "test_table",
            "ts": 1656315780000,
            "type": "DELETE"
        }
    `
	testMessageReplaceSQL = `REPLACE INTO test_db.test_table(pk1,pk2,col1,col2) VALUES (?,?,?,?),(?,?,?,?) ;`
	testMessageInsertSQL  = `INSERT INTO test_db.test_table(pk1,pk2,col1,col2) VALUES (?,?,?,?),(?,?,?,?) ON DUPLICATE KEY UPDATE col1=VALUES(col1),col2=VALUES(col2) ;`
	testMessageUpdateSQL  = `UPDATE test_db.test_table SET col1=? WHERE pk1=? AND pk2=? ;`
)

var (
	testMessagePKNames = []string{"pk1", "pk2"}
	testMessageColumns = map[string]string{"pk1": "int", "pk2": "varchar(100)", "col1": "varchar(100)", "col2": "varchar(100)"}
	testMessageNewData = []map[string]interface{}{
		map[string]interface{}{"pk1": 1, "pk2": "a", "col1": "test_new_col1_a", "col2": "test_new_col2_a"},
		map[string]interface{}{"pk1": 2, "pk2": "b", "col1": "test_new_col1_b", "col2": "test_new_col2_b"},
	}
	testMessageUpdateData = []map[string]interface{}{
		map[string]interface{}{"pk1": 1, "pk2": "a", "col1": "test_update_col1_a", "col2": "test_update_col2_a"},
		map[string]interface{}{"pk1": 2, "pk2": "b", "col1": "test_update_col1_b", "col2": "test_update_col2_b"},
	}
	testMessageOld = []map[string]interface{}{
		map[string]interface{}{"col1": "test_old_col1_a"},
		map[string]interface{}{"col1": "test_old_col2_b"},
	}
)

func TestMessage_All(t *testing.T) {
	TestMessage_GetColumnNames(t)
	TestMessage_Split(t)
	TestMessage_ConvertToSQL(t)
}

func TestMessage_GetColumnNames(t *testing.T) {
	asst := assert.New(t)

	testMessage := NewEmptyMessage()
	err := json.Unmarshal([]byte(testMessageInsertJSONString), &testMessage)
	asst.Nil(err, common.CombineMessageWithError("test GetColumnNames() failed", err))
	b, err := json.Marshal(testMessage)
	asst.Nil(err, common.CombineMessageWithError("test GetColumnNames() failed", err))
	t.Logf("testMessage:\t%s", b)

	columnNames := testMessage.GetColumnNames()
	asst.Equal(len(testMessageColumns), len(columnNames), "test GetColumnNames() failed")
	for _, columnName := range columnNames {
		asst.Contains(testMessageColumns, columnName, "test GetColumnNames() failed")
	}
}

func TestMessage_Split(t *testing.T) {
	asst := assert.New(t)

	testMessage := NewEmptyMessage()
	err := json.Unmarshal([]byte(testMessageInsertJSONString), &testMessage)
	asst.Nil(err, common.CombineMessageWithError("test Split() failed", err))
	b, err := json.Marshal(testMessage)
	asst.Nil(err, common.CombineMessageWithError("test Split() failed", err))
	t.Logf("testMessage:\t%s", b)

	messages := testMessage.Split()
	asst.Equal(len(testMessage.GetData()), len(messages), "test Split() failed")
}

func TestMessage_ConvertToSQL(t *testing.T) {
	TestMessage_convertToInsertSQL(t)
	TestMessage_convertToUpdateSQL(t)
	TestMessage_convertToDeleteSQL(t)
}

func TestMessage_convertToInsertSQL(t *testing.T) {
	asst := assert.New(t)

	testMessage := NewEmptyMessage()
	err := json.Unmarshal([]byte(testMessageInsertJSONString), &testMessage)
	asst.Nil(err, common.CombineMessageWithError("test convertToInsertSQL() failed", err))
	b, err := json.Marshal(testMessage)
	asst.Nil(err, common.CombineMessageWithError("test convertToInsertSQL() failed", err))
	t.Logf("testMessage:\t%s", b)

	// useReplace is true
	statements, err := testMessage.ConvertToSQL(true, true)
	asst.Nil(err, common.CombineMessageWithError("test convertToInsertSQL() failed", err))
	asst.Equal(1, len(statements), "test convertToInsertSQL() failed")
	for _, statement := range statements {
		for _, values := range statement {
			asst.Equal(8, len(values), "test convertToInsertSQL() failed")
		}
	}
	t.Logf("expected:\t%s", testMessageReplaceSQL)
	for _, statement := range statements {
		for sql, values := range statement {
			t.Logf("actual:\t%s", sql)
			t.Logf("values:\t%v", values)
		}
	}

	// useReplace is false
	statements, err = testMessage.ConvertToSQL(true, false)
	asst.Nil(err, common.CombineMessageWithError("test convertToInsertSQL() failed", err))
	asst.Equal(1, len(statements), "test convertToInsertSQL() failed")
	for _, statement := range statements {
		for _, values := range statement {
			asst.Equal(8, len(values), "test convertToInsertSQL() failed")
		}
	}
	t.Logf("expected:\t%s", testMessageInsertSQL)
	for _, statement := range statements {
		for sql, values := range statement {
			t.Logf("actual:\t%s", sql)
			t.Logf("values:\t%v", values)
		}
	}
}

func TestMessage_convertToUpdateSQL(t *testing.T) {
	asst := assert.New(t)

	testMessage := NewEmptyMessage()
	err := json.Unmarshal([]byte(testMessageUpdateJSONString), &testMessage)
	asst.Nil(err, common.CombineMessageWithError("test convertToUpdateSQL() failed", err))
	b, err := json.Marshal(testMessage)
	asst.Nil(err, common.CombineMessageWithError("test convertToUpdateSQL() failed", err))
	t.Logf("testMessage:\t%s", b)

	// useReplace is true
	statements, err := testMessage.ConvertToSQL(true, true)
	asst.Nil(err, common.CombineMessageWithError("test convertToUpdateSQL() failed", err))
	asst.Equal(1, len(statements), "test convertToUpdateSQL() failed")
	for _, statement := range statements {
		for _, values := range statement {
			asst.Equal(8, len(values), "test convertToUpdateSQL() failed")
		}
	}
	t.Logf("expected:\t%s", testMessageReplaceSQL)
	for _, statement := range statements {
		for sql, values := range statement {
			t.Logf("actual:\t%s", sql)
			t.Logf("values:\t%v", values)
		}
	}

	// useReplace is false
	statements, err = testMessage.ConvertToSQL(true, false)
	asst.Nil(err, common.CombineMessageWithError("test convertToUpdateSQL() failed", err))
	asst.Equal(2, len(statements), "test convertToUpdateSQL() failed")
	for _, statement := range statements {
		for _, values := range statement {
			asst.Equal(3, len(values), "test convertToUpdateSQL() failed")
		}
	}
	t.Logf("expected:\t%s", testMessageUpdateSQL)
	for _, statement := range statements {
		for sql, values := range statement {
			t.Logf("actual:\t%s", sql)
			t.Logf("values:\t%v", values)
		}
	}
}

func TestMessage_convertToDeleteSQL(t *testing.T) {

}

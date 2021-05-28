package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_All(t *testing.T) {
	TestParser_GetFingerprint(t)
	TestParser_GetSQLID(t)
	TestParser_Parse(t)
	TestParser_Split(t)
	TestParser_MergeDDLStatements(t)
}

func TestParser_GetFingerprint(t *testing.T) {
	asst := assert.New(t)

	sql := `select col1 from t01 where id = 1; select col2 from t02 where id in (select * from t03) and name = ';';select * from t04 where col1='abc'`
	p := NewParserWithDefault()

	fp := p.GetFingerprint(sql)
	asst.NotEmpty(fp, "test GetFingerprint() failed")
	t.Log(fp)
}

func TestParser_GetSQLID(t *testing.T) {
	asst := assert.New(t)

	sql := `select col1 from t01 where id = 1; select col2 from t02 where id in (select * from t03) and name = ';';select * from t04 where col1='abc'`
	p := NewParserWithDefault()

	id := p.GetSQLID(sql)
	asst.NotEmpty(id, "test GetSQLID() failed")
	t.Log(id)
}

func TestParser_Parse(t *testing.T) {
	asst := assert.New(t)

	sql := `CREATE TABLE ` + "`t01`" + `(
	 id bigint(20) comment '主键ID',
	 col1 varchar(64) NOT NULL,
	 col2 varchar(64)  NOT NULL,
	 col3 varchar(64) NOT NULL,
	 col4 int,
	 col5 mediumtext,
	 created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	 last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	 PRIMARY KEY (id),
	 KEY idx_col1_col2_col3 (col1, col2, col3)
	 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	p := NewParserWithDefault()

	result, err := p.Parse(sql)
	asst.Nil(err, "test Parse() failed")
	asst.Equal("t01", result.TableNames[0])

	// print result
	jsonBytes, err := result.Marshal()
	asst.Nil(err, "test Parse() failed")
	t.Log(string(jsonBytes))
}

func TestParser_Split(t *testing.T) {
	asst := assert.New(t)

	sql := `select col1 from t01; select col2 from t02 where id in (select * from t03) and name = ';';select * from t04`
	p := NewParserWithDefault()

	sqlList, err := p.Split(sql)
	asst.Nil(err, "test Split() failed")
	asst.Equal(3, len(sqlList))
}

func TestParser_MergeDDLStatements(t *testing.T) {
	asst := assert.New(t)

	sqls := []string{
		`create index idx01_col1 on t01(col1);`,
		`alter table t01 modify column col2 varchar(100);`,
		`alter table t01 add column col3 int(11) comment 'this is column3' after col2;`,
		`alter table t02 modify column col4 varchar(100);`,
		`alter table t02 change col5 col5 int(11) after col4;`,
		`alter table t03 add column col6 int(11); alter table t04 add column col8 int(11);alter table t03 add column col7 varchar(100);`,
		`alter table t04 modify column col9 varchar(100);`,
	}

	p := NewParserWithDefault()

	result, err := p.MergeDDLStatements(sqls...)
	asst.Nil(err, "test MergeDDLStatements() failed")
	for _, sql := range result {
		t.Log(sql)
	}
}

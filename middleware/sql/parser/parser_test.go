package parser

import (
	"encoding/json"
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

	sql := `select col1 from t01 where id = 1; select col2 from t02 where id in (select * from t03) and name = ';';select * from t04 where col1='abc';select * from t_meta_db_info where create_time<'2021-01-01';`
	p := NewParserWithDefault()

	fp := p.GetFingerprint(sql)
	asst.NotEmpty(fp, "test GetFingerprint() failed")
	t.Log(fp)
}

func TestParser_GetSQLID(t *testing.T) {
	asst := assert.New(t)

	sql := `select col1 from t01 where id = 1; select col2 from t02 where id in (select * from t03) and name = ';';select * from t04 where col1='abc';`
	p := NewParserWithDefault()

	id := p.GetSQLID(sql)
	asst.NotEmpty(id, "test GetSQLID() failed")
	t.Log(id)

	sql = `select * from t_meta_db_info where create_time<'2021-01-01'`
	id = p.GetSQLID(sql)
	asst.NotEmpty(id, "test GetSQLID() failed")
	t.Log(id)

	sql = `select              sleep(1)`
	fingerprint := p.GetFingerprint(sql)
	id = p.GetSQLID(sql)

	t.Log(fingerprint, id)
}

func TestParser_Parse(t *testing.T) {
	asst := assert.New(t)

	sql := `CREATE TABLE ` + "t01" + `(
	 id bigint(20) auto_increment comment '主键ID',
	 col1 varchar(64) character set gbk NOT NULL,
	 col2 varchar(64) collate utf8mb4_bin NOT NULL default 'abc',
	 col3 varchar(64) NOT NULL,
	 col4 int unsigned NOT NULL Default 123,
	 col5 decimal(10,2), 
	 col6 mediumtext,
	 col7 mediumblob,
	 created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	 last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	 PRIMARY KEY (id),
	 UNIQUE KEY idx01_col1 (col1) visible,
	 KEY IDX02_COL1_COL2_COL3 (col1(10) asc, col2(20) desc, col3),
	 KEY Idx03_col2(col2) invisible
	 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	// sql = `
	// 	select *
	// 	from t01
	// 			 inner join db1.t01 dt01 on t01.id = dt01.id
	// 			 inner join t02 on t01.id = t02.id
	// 			 inner join db2.t02 dt02 on dt01.id = dt02.id
	// `
	// sql = "GRANT SELECT ON *.* TO `mysql.sys`@`localhost`"
	p := NewParserWithDefault()
	p.SetParseTableDefinition(true)

	result, err := p.Parse(sql)
	asst.Nil(err, "test Parse() failed")
	// asst.Equal("t01", result.TableNames[0])

	// print result
	jsonBytes, err := json.Marshal(result)
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

func TestParser_GetTableDefinition(t *testing.T) {
	asst := assert.New(t)

	sql := `CREATE TABLE ` + "t01" + `(
	 id bigint(20) auto_increment comment '主键ID',
	 col1 varchar(64) character set gbk NOT NULL,
	 col2 varchar(64) collate utf8mb4_bin NOT NULL default 'abc',
	 col3 varchar(64) NOT NULL,
	 col4 int unsigned NOT NULL Default 123 comment 'this is col4',
	 col5 decimal(10,2), 
	 col6 mediumtext,
	 col7 mediumblob,
	 created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	 last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	 PRIMARY KEY (id),
	 UNIQUE KEY idx01_col1 (col1) visible,
	 KEY IDX02_COL1_COL2_COL3 (col1(10) asc, col2(20) desc, col3),
	 KEY Idx03_col2(col2) invisible
	 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	p := NewParserWithDefault()
	p.SetParseTableDefinition(true)

	td, err := p.GetTableDefinition(sql)
	asst.Nil(err, "test GetTableDefinition() failed")
	jsonBytes, err := json.Marshal(td)
	asst.Nil(err, "test GetTableDefinition() failed")
	t.Log(string(jsonBytes))
}

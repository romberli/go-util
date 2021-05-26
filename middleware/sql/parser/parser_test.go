package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_All(t *testing.T) {
	TestParser_Parse(t)
	TestParser_Split(t)
}

func TestParser_Parse(t *testing.T) {
	asst := assert.New(t)

	sql := `CREATE TABLE ` + "`t01`" + `(
	 id bigint(20) comment '主键ID',
	 col1 varchar(64) NOT NULL,
	 col2 varchar(64)  NOT NULL,
	 col3 varchar(64) NOT NULL,
	 col4 mediumtext,
	 col5 mediumtext,
	 created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
	 last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
	 PRIMARY KEY (id),
	 KEY idx_col1_col2_col3 (col1, col2, col3)
	 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ;`
	p := NewParserWithDefault()

	result, warns, err := p.Parse(sql)
	asst.Nil(warns, "test Parse() failed")
	asst.Nil(err, "test Parse() failed")
	asst.Equal("t01", result.TableNames[0])
}

func TestParser_Split(t *testing.T) {
	asst := assert.New(t)

	sql := `select col1 from t01; select col2 from t02 where id in (select * from t03) and name = ";";select * from t04`
	p := NewParserWithDefault()

	sqlList, warns, err := p.Split(sql)
	asst.Nil(warns, "test Split() failed")
	asst.Nil(err, "test Split() failed")
	asst.Equal(3, len(sqlList))
}

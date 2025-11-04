package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableFullDefinition_Diff(t *testing.T) {
	asst := assert.New(t)

	sourceSQL := `
		CREATE TABLE t001 (
		  id bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
		  col1 varchar(64) CHARACTER SET gbk COLLATE gbk_chinese_ci NOT NULL,
		  col2 varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'abc',
		  col3 varchar(64) NOT NULL,
		  col4 int unsigned NOT NULL DEFAULT '123' COMMENT 'this is col4',
		  col5 decimal(10,2) DEFAULT NULL,
		  col6 mediumtext,
		  col7 mediumblob,
		  created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
		  last_updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
		  PRIMARY KEY (id),
		  UNIQUE KEY idx01_col1 (col1),
		  KEY IDX02_COL1_COL2_COL3 (col1(10),col2(20) DESC,col3),
		  KEY Idx03_col2 (col2) /*!80000 INVISIBLE */
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
	`

	sourceTD, err := testParserGetTableDefinition(sourceSQL)
	asst.Nil(err, "test Diff() failed")

	targetSQL := `
		CREATE TABLE t002 (
		  id bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
		  col1 varchar(64) CHARACTER SET gbk COLLATE gbk_chinese_ci NOT NULL,
		  col21 varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT 'abc',
		  col3 varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
		  col4 int unsigned NOT NULL DEFAULT '123' COMMENT 'this is col4',
		  col6 mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
		  col5 decimal(10,2) DEFAULT NULL,
		  col7 mediumblob,
		  col8 int unsigned NOT NULL DEFAULT '123' COMMENT 'this is col8',
		  created_at datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
		  last_updated_at datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '最后更新时间',
		  PRIMARY KEY (id),
		  UNIQUE KEY idx01_col1 (col4),
		  KEY IDX02_COL1_COL2_COL3 (col1(10),col21(20) DESC,col3),
		  KEY Idx03_col21 (col21) /*!80000 INVISIBLE */
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
	`

	targetTD, err := testParserGetTableDefinition(targetSQL)
	asst.Nil(err, "test Diff() failed")

	diff := targetTD.Diff(sourceTD)
	jsonBytes, err := json.Marshal(diff)
	asst.Nil(err, "test Diff() failed")
	t.Log(string(jsonBytes))

	sql := diff.GetTableMigrationSQL()
	t.Log(sql)
}

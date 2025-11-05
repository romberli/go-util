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
	sourceSQL = `
		CREATE TABLE t_meta_insight_info (
		  insight_id bigint NOT NULL COMMENT 'insight ID',
		  insight_name varchar(100) NOT NULL COMMENT 'insight名称',
		  host_ip varchar(100) NOT NULL COMMENT 'ip地址',
		  service_port_num int NOT NULL COMMENT '服务端口号',
		  db_port_num int NOT NULL COMMENT '数据库端口号',
		  rest_port_num int NOT NULL COMMENT 'agent端口号',
		  zk_port_num int DEFAULT NULL COMMENT 'zk端口号',
		  mds_port_num int DEFAULT NULL COMMENT 'mds端口号',
		  domain_name varchar(100) DEFAULT NULL COMMENT '域名',
		  drsp_name varchar(100) DEFAULT NULL COMMENT 'drsp名称',
		  manage_user varchar(100) NOT NULL COMMENT '管理用户',
		  cpu_type tinyint NOT NULL COMMENT 'cpu类型: 1-hygon, 2-arm, 3-x86',
		  network_zone tinyint NOT NULL COMMENT '网络区域: 1-总行业务网, 2-总行办公网, 3-分行业务网, 4-分行办公网',
		  remote_insight_id bigint DEFAULT NULL COMMENT '异地insight ID',
		  az_id bigint NOT NULL COMMENT '可用区ID',
		  deployment_city varchar(100) NOT NULL COMMENT '部署城市',
		  db_proxy_cipher varchar(500) DEFAULT NULL COMMENT 'db proxy密文',
		  cn_install_path varchar(100) DEFAULT NULL COMMENT 'CN安装目录',
		  dn_install_path varchar(100) DEFAULT NULL COMMENT 'DN安装目录',
		  dn_data_path varchar(100) DEFAULT NULL COMMENT 'DN数据目录',
		  dn_log_path varchar(100) DEFAULT NULL COMMENT 'DN日志目录',
		  gtm_install_path varchar(100) DEFAULT NULL COMMENT 'GTM安装目录',
		  arch_types json DEFAULT NULL COMMENT '部署架构列表',
		  hosts json NOT NULL COMMENT '主机列表',
		  tags json NOT NULL COMMENT '标签列表',
		  parameter_template json DEFAULT NULL COMMENT '参数模板',
		  del_flag tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
		  create_time datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
		  last_update_time datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '最后更新时间',
		  PRIMARY KEY (insight_id),
		  UNIQUE KEY idx01_host_ip_service_port_num (host_ip,service_port_num),
		  UNIQUE KEY idx02_host_ip_db_port_num (host_ip,db_port_num),
		  UNIQUE KEY idx06_insight_name (insight_name),
		  KEY idx03_cpu_type_network_zone_deployment_city (cpu_type,network_zone,deployment_city),
		  KEY idx04_hosts ((cast(hosts as char(100) array))),
		  KEY idx05_tags ((cast(tags as char(100) array)))
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='insight信息表';
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

	// diff := targetTD.Diff(sourceTD)
	// jsonBytes, err := json.Marshal(diff)
	// asst.Nil(err, "test Diff() failed")
	// t.Log(string(jsonBytes))
	diff := sourceTD.Diff(targetTD)
	jsonBytes, err := json.Marshal(diff)
	asst.Nil(err, "test Diff() failed")
	t.Log(string(jsonBytes))

	sql := diff.GetTableMigrationSQL()
	t.Log(sql)
}

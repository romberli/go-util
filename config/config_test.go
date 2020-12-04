package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

type Mysqld struct {
	InnodbBufferPoolSize int64  `ini:"innodb_buffer_pool_size"`
	ReportHost           string `ini:"report_host"`
	ReportPort           int    `ini:"report_port"`
	SkipNameResolve      string `ini:"skip-name-resolve"`
}

func TestConvertToString(t *testing.T) {
	var (
		err     error
		mysqld  *Mysqld
		tagType string
		title   string
		s       string
	)

	asst := assert.New(t)

	mysqld = &Mysqld{
		InnodbBufferPoolSize: 100,
		ReportHost:           "192.168.137.11",
		ReportPort:           3306,
		SkipNameResolve:      constant.DefaultRandomString,
	}
	tagType = "ini"
	title = "[mysqld]"

	t.Log("==========test ConvertToString started.==========")
	s, err = ConvertToStringWithTitle(mysqld, title, tagType)
	asst.Nil(err, "test mysqld failed")
	t.Logf("mysqld:\n%s", s)
	t.Log("==========test ConvertToString completed.==========")
}

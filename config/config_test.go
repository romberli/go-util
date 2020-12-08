package config

import (
	"strings"
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
		err       error
		mysqld    *Mysqld
		mysqldStr string
		tagType   string
		title     string
		s         string
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
	mysqldStr = "[mysqld]\ninnodb_buffer_pool_size = 100\nreport_host = 192.168.137.11\nreport_port = 3306\nskip-name-resolve\n"

	t.Log("==========test ConvertToStringWithTitle started.==========")
	s, err = ConvertToStringWithTitle(mysqld, title, tagType)
	asst.Nil(err, "test mysqld failed")
	asst.Equal(strings.TrimSpace(mysqldStr), strings.TrimSpace(s), "test converting mysqld failed")
	t.Logf("mysqld:\n%s", s)
	t.Logf("mysqldStr:\n%s", mysqldStr)
	t.Log("==========test ConvertToStringWithTitle completed.==========")
}

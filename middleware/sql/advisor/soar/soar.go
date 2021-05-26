package soar

import (
	"github.com/XiaoMi/soar/advisor"
)

type Advisor struct {
	sqlText           string
	QueryForAuditList []*advisor.Query4Audit
}

func NewAdvisor(sql string) *Advisor {
	return newAdvisor(sql)
}

func newAdvisor(sql string) *Advisor {
	return nil
}

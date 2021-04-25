package prometheus

import (
	"fmt"
	"testing"
	"time"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	defaultAddr = "192.168.137.11:80/prometheus"
	defaultUser = "admin"
	defaultPass = "admin"
)

var conn = initConn()

func initConn() *Conn {
	config := NewConfigWithBasicAuth(defaultAddr, defaultUser, defaultPass)
	c, err := NewConnWithConfig(config)
	if err != nil {
		log.Error(fmt.Sprintf("initAppRepo() failed.\n%s", err.Error()))
		return nil
	}

	return c
}

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)

	// query := "1"
	query := `rate(mysql_global_status_queries)[1m]`
	// query := `mysql_global_status_queries`
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	step := time.Minute
	// r := apiv1.Range{
	// 	Start: start,
	// 	End:   end,
	// 	Step:  step,
	// }
	result, err := conn.Execute(query, start, end, step)
	asst.Nil(err, "test Execute() failed")
	s, err := result.GetString(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "test Execute() failed")
	ts, err := result.GetString(constant.ZeroInt, 1)
	asst.Nil(err, "test Execute() failed")
	t.Log(s, ts)
}

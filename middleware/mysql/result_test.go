package mysql

import (
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

type EnvInfo struct {
	ID             int       `middleware:"id"`
	EnvName        string    `middleware:"env_name"`
	DelFlag        bool      `middleware:"del_flag"`
	CreateTime     time.Time `middleware:"create_time"`
	LastUpdateTime time.Time `middleware:"last_update_time"`
}

func TestResult(t *testing.T) {
	var (
		err    error
		pool   *Pool
		conn   middleware.PoolConn
		result middleware.Result
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)

	addr := "192.168.137.11:3306"
	dbName := "das"
	dbUser := "root"
	dbPass := "root"

	// create pool
	pool, err = NewMySQLPoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get connection from the pool
	conn, err = pool.Get()
	asst.Nil(err, "get connection from pool failed")

	// map to struct
	sql := `select id, env_name, del_flag, create_time, last_update_time from t_meta_env_info;`
	result, err = conn.Execute(sql)
	asst.Nil(err, "execute sql failed.")
	envInfoList := make([]interface{}, result.RowNumber())
	for i := range envInfoList {
		envInfoList[i] = &EnvInfo{}
	}
	err = result.MapToStruct(envInfoList, constant.DefaultTagType)
	asst.Nil(err, "map to struct failed")
}

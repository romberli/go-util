package mysql

import (
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

type EnvInfo struct {
	UnknownField   *Pool
	ID             int       `middleware:"id"`
	EnvName        string    `middleware:"env_name"`
	DelFlag        bool      `middleware:"del_flag"`
	CreateTime     time.Time `middleware:"create_time"`
	LastUpdateTime time.Time `middleware:"last_update_time"`
}

type TestStruct struct {
	OK bool `middleware:"ok"`
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
	pool, err = NewPoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get connection from the pool
	conn, err = pool.Get()
	asst.Nil(err, "get connection from pool failed")

	sql := `select id, env_name, del_flag, create_time, last_update_time from t_meta_env_info;`
	result, err = conn.Execute(sql)
	asst.Nil(err, "execute sql failed")
	// map to int slice
	idList := make([]int, result.RowNumber())
	err = result.MapToIntSlice(idList, constant.ZeroInt)
	asst.Nil(err, "map to int slice failed")
	// map to string slice
	envNameList := make([]string, result.RowNumber())
	err = result.MapToStringSlice(envNameList, constant.OneInt)
	asst.Nil(err, "map to string slice failed")
	// map to float64 slice
	delFlagList := make([]float64, result.RowNumber())
	err = result.MapToFloatSlice(delFlagList, constant.TwoInt)
	asst.Nil(err, "map to float64 slice failed")

	// map to struct
	envInfoList := make([]*EnvInfo, result.RowNumber())
	for i := range envInfoList {
		envInfoList[i] = &EnvInfo{}
	}
	value, err := result.GetMap(0, 0)
	asst.NotNil(err, "get map value failed")
	asst.Nil(value, "get map value failed")
	err = result.MapToStructSlice(envInfoList, constant.DefaultMiddlewareTag)
	asst.Nil(err, "map to struct failed")

	sql = `select 'true' as ok;`
	result, err = conn.Execute(sql)
	asst.Nil(err, "execute sql failed")
	testStructList := make([]*TestStruct, result.RowNumber())
	for i := range testStructList {
		testStructList[i] = &TestStruct{}
	}
	valueB, err := result.GetBool(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "execute sql failed")
	asst.True(valueB, "execute sql failed")
	err = result.MapToStructSlice(testStructList, constant.DefaultMiddlewareTag)
	asst.Nil(err, "map to struct failed")
}

func TestResult_Tmp(t *testing.T) {
	asst := assert.New(t)

	type ServerInfo struct {
		ServerID    int      `middleware:"server_id"`
		HostTags    string   `middleware:"host_tags"`
		HostTagList []string `middleware:"host_tag_list"`
	}

	addr := "192.168.137.11:3306"
	dbName := "gp"
	dbUser := "root"
	dbPass := "root"

	sql := `select server_id, host_tags, host_tags as host_tag_list from t_meta_server_info;`
	serverInfo := &ServerInfo{}

	conn, err := NewConn(addr, dbName, dbUser, dbPass)
	asst.Nil(err, common.CombineMessageWithError("create connection failed", err))
	result, err := conn.Execute(sql)
	asst.Nil(err, common.CombineMessageWithError("execute sql failed", err))
	err = result.MapToStructByRowIndex(serverInfo, 0, constant.DefaultMiddlewareTag)
	asst.Nil(err, common.CombineMessageWithError("map to struct failed", err))
	t.Logf(serverInfo.HostTags)
	t.Logf("%v", serverInfo.HostTagList)
}

package constant

import (
	"os"
)

const (
	DefaultFileMode         os.FileMode = 0644
	DefaultExecFileMode     os.FileMode = 0755
	DefaultKillSignal                   = 15
	DefaultNormalExitCode               = 0
	DefaultAbnormalExitCode             = 1
	ZeroInt                             = 0
	MinUInt                 uint        = 0
	MaxUInt                             = ^uint(0)
	MaxInt                              = int(^uint(0) >> 1)
	MinInt                              = ^MaxInt
	MinPort                             = 1
	MaxPort                             = 65535
	EmptyString                         = ""
	CurrentDir                          = "./"
	DefaultRandomString                 = "sadfio3mj23gsk9lj8ou"
	DefaultRandomInt                    = 2053167498
	DefaultRandomFloat      float64     = DefaultRandomInt
	TrueString                          = "true"
	FalseString                         = "false"
	CRLFString                          = "\n"
	CommaString                         = ","
	DefaultRandomTimeString             = "9999-07-02 09:55:32.346082"
	DefaultTimeLayout                   = "2006-01-02 15:04:05.999999"
	DefaultMarshalTag                   = "json"
	DefaultMiddlewareTag                = "middleware"
	DefaultListenIP                     = "0.0.0.0"
	DefaultLocalHostName                = "localhost"
	DefaultLocalHostIP                  = "127.0.0.1"
	DefaultMySQLPort                    = 3306
	DefaultRootUserName                 = "root"
	DefaultRootUserPass                 = "root"
)

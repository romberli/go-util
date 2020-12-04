package constant

import (
	"os"
)

const (
	ZeroInt                         = 0
	MinUInt             uint        = 0
	MaxUInt                         = ^uint(0)
	MaxInt                          = int(^uint(0) >> 1)
	MinInt                          = ^MaxInt
	MinPort                         = 1
	MaxPort                         = 65535
	EmptyString                     = ""
	CurrentDir                      = "./"
	DefaultRandomString             = "sadfio3mj23gsk9lj8ou"
	DefaultRandomInt                = 345920654907418
	TrueString                      = "true"
	FalseString                     = "false"
	CRLFString                      = "\n"
	DefaultFileMode     os.FileMode = 0644
	DefaultExecFileMode os.FileMode = 0755
)

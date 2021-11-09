package constant

import (
	"os"
	"time"
)

const (
	EmptyString             = ""
	SpaceString             = " "
	CRLFString              = "\n"
	CommaString             = ","
	AsteriskString          = "*"
	DotString               = "."
	VerticalBarString       = "|"
	ColonString             = ":"
	SemicolonString         = ";"
	LeftParenthesisString   = "("
	RightParenthesisString  = ")"
	LeftBracketString       = "["
	RightBracketString      = "]"
	LeftAngleBracketString  = "<"
	RightAngleBracketString = ">"
	LeftBraceString         = "{"
	RightBraceString        = "}"
	SlashString             = "/"
	BackSlashString         = "\\"
	DashString              = "-"
	UnderBarString          = "_"
	PlusString              = "+"
	MinusString             = "-"
	MultiplicationString    = "*"
	DivisionString          = "/"
	EqualString             = "="
	UnequalString           = "!="
	LargerString            = ">"
	LargerEqualString       = ">="
	SmallerString           = "<"
	SmallerEqualString      = "<="

	NullString                          = "null"
	NoneString                          = "none"
	NanString                           = "nan"
	TransportProtocolTCP                = "tcp"
	TransportProtocolUDP                = "udp"
	GOOSLinux                           = "linux"
	GOOSDarwin                          = "darwin"
	GOOSWindows                         = "windows"
	DefaultFileMode         os.FileMode = 0644
	DefaultExecFileMode     os.FileMode = 0755
	DefaultKillSignal                   = 15
	DefaultNormalExitCode               = 0
	DefaultAbnormalExitCode             = 1
	MaxPercentage                       = 100
	ZeroInt                             = 0
	MinUInt                 uint        = 0
	MaxUInt                             = ^uint(0)
	MaxInt                              = int(^uint(0) >> 1)
	MinInt                              = ^MaxInt
	MinPort                             = 1
	MaxPort                             = 65535

	RootDir                         = "/"
	CurrentDir                      = "./"
	DefaultRandomString             = "sadfio3mj23gsk9lj8ou"
	DefaultRandomInt                = 2053167498
	DefaultRandomFloat      float64 = DefaultRandomInt
	TrueString                      = "true"
	FalseString                     = "false"
	DefaultRandomTimeString         = "9999-07-02 09:55:32.346082"
	TimeLayoutSecond                = "2006-01-02 15:04:05"
	TimeLayoutMicrosecond           = "2006-01-02 15:04:05.999999"
	DefaultTimeLayout               = TimeLayoutMicrosecond
	Day                             = 24 * time.Hour
	Week                            = 7 * Day
	DefaultMarshalTag               = DefaultJSONTag
	DefaultJSONTag                  = "json"
	DefaultMiddlewareTag            = "middleware"
	DefaultListenIP                 = "0.0.0.0"
	DefaultLocalHostName            = "localhost"
	DefaultLocalHostIP              = "127.0.0.1"
	DefaultMySQLPort                = 3306
	DefaultRootUserName             = "root"
	DefaultRootUserPass             = "root"
)

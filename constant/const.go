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
	QuestionMarkString      = "?"
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
	SingleQuoteString       = "'"
	DoubleQuoteString       = "\""
	FourSpaceString         = "    "
	TabString               = "\t"
	DefaultIndentString     = FourSpaceString
	NullString              = "null"
	NoneString              = "none"
	NanString               = "nan"
	EmptyListString         = "[]"
	EmptyMapString          = "{}"

	LogWithStackString = "%+v"

	TransportProtocolTCP                = "tcp"
	TransportProtocolUDP                = "udp"
	GOOSLinux                           = "linux"
	GOOSDarwin                          = "darwin"
	GOOSWindows                         = "windows"
	DefaultFileMode         os.FileMode = 0644
	DefaultExecFileMode     os.FileMode = 0755
	DefaultAllFileMode      os.FileMode = 0777
	DefaultFileModeStr                  = "0644"
	DefaultExecFileModeStr              = "0755"
	DefaultAllFileModeStr               = "0777"
	DefaultKillSignal                   = 15
	DefaultNormalExitCode               = 0
	DefaultAbnormalExitCode             = 1
	MaxPercentage                       = 100
	ZeroInt                             = 0
	OneInt                              = 1
	TwoInt                              = 2
	ThreeInt                            = 3
	FourInt                             = 4
	FiveInt                             = 5
	SixInt                              = 6
	SevenInt                            = 7
	EightInt                            = 8
	NineInt                             = 9
	TenInt                              = 10
	HundredInt                          = 100
	ThousandInt                         = 1000
	MillionInt                          = 1000000
	BillionInt                          = 1000000000
	TrillionInt                         = 1000000000000
	QuadrillionInt                      = 1000000000000000
	KiloInt                             = 1024
	MegaInt                             = KiloInt * KiloInt
	GigaInt                             = KiloInt * MegaInt
	TeraInt                             = KiloInt * GigaInt
	PetaInt                             = KiloInt * TeraInt
	ExaInt                              = KiloInt * PetaInt

	MinUInt uint = 0
	MaxUInt      = ^uint(0)
	MaxInt       = int(^uint(0) >> 1)
	MinInt       = ^MaxInt
	MinPort      = 1
	MaxPort      = 65535

	RootDir                         = "/"
	CurrentDir                      = "./"
	DefaultTmpDir                   = "/tmp"
	DefaultRandomString             = "sadfio3mj23gsk9lj8ou"
	DefaultRandomInt                = 2053167498
	DefaultRandomFloat      float64 = DefaultRandomInt
	TrueString                      = "true"
	FalseString                     = "false"
	DefaultRandomTimeString         = "9999-07-02 09:55:32.346082"
	TimeLayoutSecond                = "2006-01-02 15:04:05"
	TimeLayoutMicrosecond           = "2006-01-02 15:04:05.999999"
	TimeLayoutSecondDash            = "20060102-150405"
	TimeLayoutSecondInt             = "20060102150405"
	DefaultTimeLayout               = TimeLayoutMicrosecond
	Day                             = 24 * time.Hour
	Week                            = 7 * Day
	Month                           = 30 * Day
	Year                            = 365 * Day
	Century                         = 100 * Year

	DefaultMarshalTag    = DefaultJSONTag
	DefaultJSONTag       = "json"
	DefaultMiddlewareTag = "middleware"
	DefaultListenIP      = "0.0.0.0"
	DefaultLocalHostName = "localhost"
	DefaultLocalHostIP   = "127.0.0.1"

	DefaultRootUserName = "root"
	DefaultRootUserPass = "root"

	DefaultSSHPort = 22

	DefaultMySQLPort = 3306
	DefaultMySQLAddr = "127.0.0.1:3306"

	DefaultRedisPort = 6379
	DefaultRedisAddr = "127.0.0.1:6379"

	DefaultRabbitmqPort  = 5672
	DefaultRabbitmqAddr  = "127.0.0.1:5672"
	DefaultGuestUserName = "guest"
	DefaultGuestUserPass = "guest"
	DefaultVhost         = "/"

	HTTPSchemePrefix       = "http://"
	HTTPSSchemePrefix      = "https://"
	DefaultHTTPPort        = 80
	DefaultTextContentType = "text/plain"
	DefaultJSONContentType = "application/json"

	AArch64Arch = "aarch64"
	X64Arch     = "x86_64"

	UTF8MB4Charset = "utf8mb4"
	GB18030Charset = "gb18030"
	GBKCharset     = "gbk"
)

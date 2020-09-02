package constant

const (
	MinUInt uint = 0
	MaxUInt      = ^uint(0)
	MaxInt       = int(^uint(0) >> 1)
	MinInt       = ^MaxInt
	MinPort      = 1
	MaxPort      = 65535
)

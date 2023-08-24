package types

type Primitive interface {
	~bool | ~string | Number
}

type Number interface {
	Int | UnsignedInt | Float
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

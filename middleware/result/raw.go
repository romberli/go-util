package result

type Raw interface {
	// GetRaw returns the raw data of the result
	GetRaw() interface{}
}

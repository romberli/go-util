package result

type Bool interface {
	// GetBool returns bool type value of given row and column number
	GetBool(row, column int) (bool, error)
	// GetBoolByName returns bool type value of given row number and column name
	GetBoolByName(row int, name string) (bool, error)
}

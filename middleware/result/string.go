package result

type String interface {
	// GetString returns string type value of given row and column number
	GetString(row, column int) (string, error)
	// GetStringByName returns string type value of given row number and column name
	GetStringByName(row int, name string) (string, error)
}

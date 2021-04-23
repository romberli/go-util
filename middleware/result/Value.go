package result

type Value interface {
	// RowNumber returns how many rows in the result
	RowNumber() int
	// ColumnNumber return how many columns in the result
	ColumnNumber() int
	// GetValue returns interface{} type value of given row and column number
	GetValue(row, column int) (interface{}, error)
	// ColumnExists check if column exists in the result
	ColumnExists(name string) bool
	// NameIndex returns number of given column
	NameIndex(name string) (int, error)
	// GetValueByName returns interface{} type value of given row number and column name
	GetValueByName(row int, name string) (interface{}, error)
	// IsNull checks if value of given row and column number is nil
	IsNull(row, column int) (bool, error)
	// IsNullByName checks if value of given row number and column name is nil
	IsNullByName(row int, name string) (bool, error)
}

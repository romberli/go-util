package result

type Unmarshaler interface {
	// MapToIntSlice maps each row to a int slice,
	// first argument must be a slice of int,
	// only the specified column of the row will be mapped,
	// the column must be able to convert to int, and will map to the int in the slice
	MapToIntSlice(in []int, column int) error
	// MapToStringSlice maps each row to a string slice,
	// first argument must be a slice of string,
	// only the specified column of the row will be mapped,
	// the column must be able to convert to string, and will map to the string in the slice
	MapToStringSlice(in []string, column int) error
	// MapToFloatSlice maps each row to a float slice,
	// first argument must be a slice of float64,
	// only the specified column of the row will be mapped,
	// the column must be able to convert to int, and will map to the float in the slice
	MapToFloatSlice(in []float64, column int) error
	// MapToStructSlice maps each row to a struct of the first argument,
	// first argument must be a slice of pointers to structs,
	// each row in the result maps to a struct in the slice,
	// each column in the row maps to a field of the struct,
	// tag argument is the tag of the field, it represents the column name,
	// if there is no such tag in the field, this field will be ignored,
	// so set tag to each field that need to be mapped,
	// using "middleware" as the tag is recommended.
	MapToStructSlice(in interface{}, tag string) error
	// MapToStructByRowIndex maps row of given index result to the struct
	// first argument must be a pointer to struct,
	// each column in the row maps to a field of the struct,
	// tag argument is the tag of the field, it represents the column name,
	// if there is no such tag in the field, this field will be ignored,
	// so set tag to each field that need to be mapped,
	// using "middleware" as the tag is recommended.
	MapToStructByRowIndex(in interface{}, row int, tag string) error
}

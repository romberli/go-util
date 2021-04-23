package result

type Unmarshaler interface {
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

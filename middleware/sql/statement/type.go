package statement

import (
	"strings"
)

type Type int

const (
	Unknown Type = iota
	Select
	Insert
	Update
	Delete
	Create
	Alter
	Drop

	selectPrefix = "select"
	insertPrefix = "insert"
	updatePrefix = "update"
	deletePrefix = "delete"
	createPrefix = "create"
	alterPrefix  = "alter"
	dropPrefix   = "drop"
)

func GetType(sql string) Type {
	sql = strings.ToLower(strings.TrimSpace(sql))

	switch {
	case strings.HasPrefix(sql, selectPrefix):
		return Select
	case strings.HasPrefix(sql, insertPrefix):
		return Insert
	case strings.HasPrefix(sql, updatePrefix):
		return Update
	case strings.HasPrefix(sql, deletePrefix):
		return Delete
	case strings.HasPrefix(sql, createPrefix):
		return Create
	case strings.HasPrefix(sql, alterPrefix):
		return Alter
	case strings.HasPrefix(sql, dropPrefix):
		return Drop
	default:
		return Unknown
	}
}

package sqls

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
)

const (
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
	if strings.HasPrefix(sql, selectPrefix) {
		return Select
	}
	if strings.HasPrefix(sql, insertPrefix) {
		return Insert
	}
	if strings.HasPrefix(sql, updatePrefix) {
		return Update
	}
	if strings.HasPrefix(sql, deletePrefix) {
		return Delete
	}
	if strings.HasPrefix(sql, createPrefix) {
		return Create
	}
	if strings.HasPrefix(sql, alterPrefix) {
		return Alter
	}
	if strings.HasPrefix(sql, dropPrefix) {
		return Drop
	}

	return Unknown
}

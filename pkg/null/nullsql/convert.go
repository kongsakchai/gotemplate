package nullsql

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/pkg/null"
)

func ToNullSQL[T any](src null.N[T]) sql.Null[T] {
	return sql.Null[T]{
		V:     src.V,
		Valid: src.Valid,
	}
}

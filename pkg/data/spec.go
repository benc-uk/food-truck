package data

import "database/sql"

type Database interface {
	QuerySimple(query string) (*sql.Rows, error)
}

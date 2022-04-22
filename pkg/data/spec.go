package data

import (
	"database/sql"
)

// The Database spec defines the interface for a database
type Database interface {
	QueryTrucks(query string) ([]TruckRow, error)
}

// "TruckRow" is data returned from the database
type TruckRow struct {
	ID          sql.NullString
	Name        sql.NullString
	Description sql.NullString
	Lat         sql.NullFloat64
	Long        sql.NullFloat64
	Address     sql.NullString
}

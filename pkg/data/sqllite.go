package data

import (
	"database/sql"
	"fmt"
	"log"

	// Import sqlite
	_ "github.com/mattn/go-sqlite3"
)

// A concrete implementation of data.Database that uses SQLite
type SQLiteDB struct {
	*sql.DB
}

func NewDatabase(dbFilePath string) *SQLiteDB {
	// Note force rw mode here, otherwise it creates an empty DB if file not found
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=rw", dbFilePath))

	if err != nil {
		log.Fatalln(err)
		log.Fatalln("### Failed to open database, can not start")
	}

	return &SQLiteDB{
		db,
	}
}

// QueryTrucks returns a list of trucks from the database
// TODO: This is a leaky abstraction, we could do A LOT better here! Needs to be refactored
func (db *SQLiteDB) QueryTrucks(q string) ([]TruckRow, error) {
	rows, err := db.Query(q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := []TruckRow{}

	defer rows.Close()
	for rows.Next() {
		t := TruckRow{}
		err := rows.Scan(&t.ID, &t.Name, &t.Lat, &t.Long, &t.Address, &t.Description)
		if err != nil {
			return nil, err
		}
		data = append(data, t)
	}

	return data, nil
}

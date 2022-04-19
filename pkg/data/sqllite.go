package data

import (
	"database/sql"
	"fmt"
	"log"

	// Import sqlite
	_ "github.com/mattn/go-sqlite3"
)

// An implementation of data.Database using in memory KV store
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

func (db *SQLiteDB) QuerySimple(q string) (*sql.Rows, error) {
	rows, err := db.Query(q)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return rows, nil
}

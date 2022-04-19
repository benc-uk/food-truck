package trucks

import (
	"database/sql"
	"os"
	"testing"
)

//
// Mock DB implementation
//
type MockDB struct {
}

func (db *MockDB) QuerySimple(q string) (*sql.Rows, error) {
	return nil, nil
}

//
// Initialize service for testing with mock DB
//
var service Service

func TestMain(m *testing.M) {
	service = NewService(&MockDB{})

	// Now run the rest of the tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

//
// Test cases for the Truck Service
//
func TestFindNear(t *testing.T) {

	// TODO: Tests removed due to lack of time to mock the database

	// trucks, err := service.FindNear(23.8, 34.0, 500)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if len(trucks) <= 0 {
	// 	t.Error("Expected results but got zero")
	// }
}

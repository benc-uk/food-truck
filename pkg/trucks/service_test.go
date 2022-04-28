package trucks

//
// Unit tests for Trucks Service
//

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/benc-uk/food-truck/pkg/data"
)

// A concrete implementation of data.Database that is a mock
type MockDB struct {
	mockdata []data.TruckRow
}

//
// NewMockDB returns a new mock database for testing
//
func NewMockDB() *MockDB {
	return &MockDB{
		mockdata: []data.TruckRow{
			{
				ID:          sql.NullString{"1", true},
				Name:        sql.NullString{"Truck 1", true},
				Description: sql.NullString{"Truck 1 description", true},
				Lat:         sql.NullFloat64{23.8, true},
				Long:        sql.NullFloat64{34.0, true},
				Address:     sql.NullString{"Truck 1 address", true},
			},
			{
				ID:          sql.NullString{"2", true},
				Name:        sql.NullString{"Truck 2", true},
				Description: sql.NullString{"Truck 2 description", true},
				Lat:         sql.NullFloat64{30.1, true},
				Long:        sql.NullFloat64{34.0, true},
				Address:     sql.NullString{"Truck 2 address", true},
			},
		},
	}
}

func (db *MockDB) QueryTrucks(q string) ([]data.TruckRow, error) {
	return db.mockdata, nil
}

//
// Initialize service for testing with mock DB
//
var service Service

func TestMain(m *testing.M) {
	service = NewService(NewMockDB())

	// Now run the rest of the tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

//
// Test cases for the Truck Service
//
func TestFindNear(t *testing.T) {
	trucks, err := service.FindNear(23.8, 34.0, 500)
	if err != nil {
		log.Fatal(err)
		t.Error(err)
	}

	if len(trucks) <= 0 {
		t.Error("Expected results but got zero")
	}
}

func TestFindNearBad(t *testing.T) {
	_, err := service.FindNear(123, 456, -50)
	if err == nil {
		t.Error("Expected error not returned")
	}
}

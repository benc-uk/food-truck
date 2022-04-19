package trucks

import (
	"fmt"
	"log"

	"github.com/benc-uk/food-truck/pkg/data"
)

// TruckService is implementation of the Service interface backed by a database
type TruckService struct {
	db data.Database
}

// NewService returns a new instance of TruckService
func NewService(db data.Database) *TruckService {
	return &TruckService{
		// Another embedding
		db: db,
	}
}

// FindNear returns a list of trucks near a given lat & long and within a radius
func (s *TruckService) FindNear(lat float64, long float64, radius int) ([]Truck, error) {
	log.Printf("Finding trucks near %f, %f", lat, long)

	// TODO: This is a HACK to approximate the radius in meters in long/lat
	const oneMLat = 0.000009
	const oneMLong = 0.000014

	latMinStr := fmt.Sprintf("%f", lat-(oneMLat*float64(radius)))
	latMaxStr := fmt.Sprintf("%f", lat+(oneMLat*float64(radius)))
	longMinStr := fmt.Sprintf("%f", long-(oneMLong*float64(radius)))
	longMaxStr := fmt.Sprintf("%f", long+(oneMLong*float64(radius)))

	// TODO: Phase 1: Use Euclidian/Pythagorean distance, i.e. a circle with SQRT(x^2 + y^2)
	// TODO: Phase 2: Migrate to a REAL spatial database, e.g. Azure Cosmos DB or SQL Server
	rows, err := s.db.QuerySimple(`SELECT locationId, value as name, Latitude, Longitude, Address, FoodItems 
	FROM Mobile_Food_Facility_Permit 
	INNER JOIN Applicant ON Applicant.id = Mobile_Food_Facility_Permit.Applicant 
	WHERE Latitude != 0 AND Latitude > ` + latMinStr + ` AND Latitude < ` + latMaxStr +
		` AND Longitude > ` + longMinStr + ` AND Longitude < ` + longMaxStr)

	if err != nil {
		return nil, &Error{Code: ErrorOther, Message: "Failed to fetch trucks"}
	}

	trucks := []Truck{}

	defer rows.Close()
	for rows.Next() {
		t := Truck{}
		err := rows.Scan(&t.ID, &t.Name, &t.Lat, &t.Long, &t.addressNullStr, &t.descNullStr)
		if t.descNullStr.Valid {
			t.Description = t.descNullStr.String
		}
		if t.addressNullStr.Valid {
			t.Address = t.addressNullStr.String
		}
		if err != nil {
			return nil, err
		}
		trucks = append(trucks, t)
	}

	return trucks, nil
}

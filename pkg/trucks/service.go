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
	if radius < 0 {
		return nil, fmt.Errorf("radius must be greater than 0")
	}

	// TODO: This is a HACK to approximate the radius in meters in long/lat
	const oneMLat = 0.000009
	const oneMLong = 0.000014

	latMinStr := fmt.Sprintf("%f", lat-(oneMLat*float64(radius)))
	latMaxStr := fmt.Sprintf("%f", lat+(oneMLat*float64(radius)))
	longMinStr := fmt.Sprintf("%f", long-(oneMLong*float64(radius)))
	longMaxStr := fmt.Sprintf("%f", long+(oneMLong*float64(radius)))

	// TODO: Phase 1: Use Euclidian/Pythagorean distance, i.e. a circle with SQRT(x^2 + y^2)
	// TODO: Phase 2: Migrate to a REAL spatial database, e.g. Azure Cosmos DB or SQL Server
	rows, err := s.db.QueryTrucks(`SELECT locationId, value as name, Latitude, Longitude, Address, FoodItems 
	FROM Mobile_Food_Facility_Permit 
	INNER JOIN Applicant ON Applicant.id = Mobile_Food_Facility_Permit.Applicant 
	WHERE Latitude != 0 AND Latitude > ` + latMinStr + ` AND Latitude < ` + latMaxStr +
		` AND Longitude > ` + longMinStr + ` AND Longitude < ` + longMaxStr)

	if err != nil {
		return nil, &Error{Code: ErrorOther, Message: "Failed to fetch trucks"}
	}

	log.Printf("Found %d trucks", len(rows))

	// Map the rows to trucks
	trucks := []Truck{}
	for _, truckRow := range rows {
		truck := Truck{}
		if truckRow.ID.Valid {
			truck.ID = truckRow.ID.String
		}
		if truckRow.Name.Valid {
			truck.Name = truckRow.Name.String
		}
		if truckRow.Description.Valid {
			truck.Description = truckRow.Description.String
		}
		if truckRow.Lat.Valid {
			truck.Lat = truckRow.Lat.Float64
		}
		if truckRow.Long.Valid {
			truck.Long = truckRow.Long.Float64
		}
		if truckRow.Address.Valid {
			truck.Address = truckRow.Address.String
		}
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

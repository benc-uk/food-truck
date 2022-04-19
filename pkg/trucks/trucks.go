package trucks

import (
	"database/sql"

	"github.com/benc-uk/food-truck/pkg/api"
)

// API is the main truck API
type API struct {
	Service Service
	// Use composition and embedding to extend the API base
	api.Base
}

// Truck is the truck model, used for both requests and responses
// swagger:model truck
type Truck struct {
	// The truck ID
	ID string `json:"id"`

	// The name of the truck
	Name string `json:"name"`

	descNullStr sql.NullString
	// Descripton / long text
	Description string `json:"description"`

	// Latitude location
	Lat float64 `json:"lat"`

	// Logitude location
	Long float64 `json:"long"`

	addressNullStr sql.NullString
	// Address in human readable form
	Address string `json:"address"`
}

// Service provides truck operations
type Service interface {
	FindNear(lat float64, long float64, radius int) ([]Truck, error)
}

const (
	// ErrorOther is a generic error
	ErrorOther int = iota
	// ErrorDuplicate is returned when a duplicate truck is created
	ErrorDuplicate
	// ErrorNotFound is	returned when a truck is not found
	ErrorNotFound
	// ErrorValidation is returned when a truck is not valid or JSON parsing fails
	ErrorValidation
)

// Error is an application specific error type which also provides an error code.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

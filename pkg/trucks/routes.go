package trucks

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const startRadius = 400

// AddRoutes registers all application routes with the API
func (a *API) AddRoutes(router *mux.Router) {
	router.HandleFunc("/trucks/{lat}/{long}", a.findTrucksNear).Methods("GET")
	router.HandleFunc("/", a.redirectToApp).Methods("GET")
	a.Healthy = true
}

// Use when hitting the base or root
func (a *API) redirectToApp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/", http.StatusFound)
}

// swagger:operation GET /trucks/{lat}/{long} findTrucksNear
// Returns a list of trucks near a given lat & long
// ---
// produces:
// - application/json
// parameters:
// - name: lat
//   in: path
//   description: Latitude to search around
// - name: long
//   in: path
//   description: Longitude to search around
// - name: radius
//   in: query
//   description: Radius of search in meters (approx)
// responses:
//   '200':
//     description: truck
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/truck"
//   '400':
//     description: Input validation error
//   '500':
//     description: Other error
func (a *API) findTrucksNear(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	latStr := params["lat"]
	longStr := params["long"]
	radius, err := strconv.Atoi(r.URL.Query().Get("radius"))
	if err != nil {
		radius = 0
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trucks := []Truck{}

	// This handles the two modes of operation, either we're getting a list of trucks near a given lat/long and given radius,
	// or we're searching in in increasing radius until we find at least 5 trucks
	if radius > 0 {
		trucksResult, err := a.Service.FindNear(lat, long, radius)
		if err != nil {
			truckErr, isTruckError := err.(*Error)
			if isTruckError {
				http.Error(w, truckErr.Message, httpErrorCode(truckErr))
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		trucks = append(trucks, trucksResult...)
	} else {
		// HACK: This brute force approach is less than ideal, this should be moved to a spatial DB ASAP
		radius = startRadius
		for len(trucks) < 5 {
			log.Printf("Searching for more trucks with radius %d", radius)
			trucksResult, err := a.Service.FindNear(lat, long, radius)
			if err != nil {
				truckErr, isTruckError := err.(*Error)
				if isTruckError {
					http.Error(w, truckErr.Message, httpErrorCode(truckErr))
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			trucks = append(trucks, trucksResult...)
			// Double the radius for the next search
			radius *= 2
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(trucks) == 0 {
		// Trick to return an empty JSON array, rather than null
		_ = json.NewEncoder(w).Encode([]string{})
		return
	}

	_ = json.NewEncoder(w).Encode(trucks)
}

// httpErrorCode maps truck errors to HTTP status codes
func httpErrorCode(err *Error) int {
	switch err.Code {
	case ErrorDuplicate:
		return http.StatusConflict
	case ErrorNotFound:
		return http.StatusNotFound
	case ErrorValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

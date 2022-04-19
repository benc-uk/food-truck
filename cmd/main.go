package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benc-uk/food-truck/pkg/api"
	"github.com/benc-uk/food-truck/pkg/data"
	"github.com/flowchartsman/swaggerui"

	"github.com/benc-uk/food-truck/pkg/trucks"
	"github.com/gorilla/mux"
)

const serverPort = "8080"

//go:generate ../bin/swagger generate spec -m -o swagger.yaml
//go:embed swagger.yaml
var swaggerSpec []byte

func main() {
	router := mux.NewRouter()

	dbPath := "./data/food-trucks.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	log.Printf("Using database: %s", dbPath)

	frontendDir := "./web/client"
	if len(os.Args) > 2 {
		frontendDir = os.Args[2]
	}
	log.Println("Using frontend dir:", frontendDir)

	truckAPI := &trucks.API{
		Service: trucks.NewService(data.NewDatabase(dbPath)),
		Base: api.Base{
			Healthy: true,
			Version: "1.0.0",
			Name:    "Food Truck API",
		},
	}

	// Bind application routes to the router
	truckAPI.AddRoutes(router)

	// Add logging, health & metrics middleware
	truckAPI.AddLogging(router)
	truckAPI.AddMetrics(router)
	truckAPI.AddHealth(router)
	truckAPI.AddStatus(router)

	// Lastly add the swagger UI
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", swaggerui.Handler(swaggerSpec)))

	// Static content to serve the frontend
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir(frontendDir))))

	server := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           truckAPI.CORSHandler([]string{"*"}, router),
		Addr:              ":" + serverPort,
	}

	log.Printf("HTTP server for '%s' starting on port %s", truckAPI.Name, serverPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

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

	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

var version = "0.0.0"           // App version number, set at build time with -ldflags "-X main.version=1.2.3"
var buildInfo = "No build info" // Build details, set at build time with -ldflags "-X main.buildInfo=blah"

//go:generate ../bin/swagger generate spec -m -o swagger.yaml
//go:embed swagger.yaml
var swaggerSpec []byte

func main() {
	router := mux.NewRouter()

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/food-trucks.db"
	}

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "./web/client"
	}

	azureMapsKey := os.Getenv("AZURE_MAPS_KEY")
	if azureMapsKey == "" {
		log.Fatalln("FATAL! Missing AZURE_MAPS_KEY environment variable")
	}

	log.Printf("Using database: %s", dbPath)
	log.Println("Using frontend dir:", frontendDir)

	truckAPI := &trucks.API{
		Service: trucks.NewService(data.NewDatabase(dbPath)),
		Base: api.Base{
			Healthy: true,
			Version: version,
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

	// Add the swagger UI
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", swaggerui.Handler(swaggerSpec)))

	// Static content handler, to serve the frontend
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir(frontendDir))))

	// Very simple config API to return the AZURE_MAPS_KEY env var
	router.Path("/config").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"azureMapsKey": "` + azureMapsKey + `"}`))
	})

	server := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		// CORS needs to be added here, and can't be done as middleware
		Handler: truckAPI.CORSHandler([]string{"*"}, router),
		Addr:    ":" + serverPort,
	}

	log.Printf("API version: %s - Build details: %s", truckAPI.Version, buildInfo)
	log.Printf("HTTP server for '%s' starting on port %s", truckAPI.Name, serverPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

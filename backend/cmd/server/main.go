package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/tgschnetzer/nascar-pool/internal/cache"
	"github.com/tgschnetzer/nascar-pool/internal/database"
	"github.com/tgschnetzer/nascar-pool/internal/handlers"
)

func main() {
	// Load .env file - check both current dir and parent (project root)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Connect to Redis (optional - used for caching)
	if err := cache.InitRedis(); err != nil {
		log.Printf("Redis connection failed: %v", err)
	}

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Participants
	api.HandleFunc("/participants", handlers.GetParticipants).Methods("GET")

	// Drivers
	api.HandleFunc("/drivers", handlers.GetDrivers).Methods("GET")

	// Races
	api.HandleFunc("/races", handlers.GetRaces).Methods("GET")
	api.HandleFunc("/races", handlers.CreateRace).Methods("POST")
	api.HandleFunc("/races/{id}", handlers.GetRace).Methods("GET")
	api.HandleFunc("/races/{id}", handlers.UpdateRace).Methods("PUT")
	api.HandleFunc("/races/{id}/generate-teams", handlers.GenerateTeams).Methods("POST")
	api.HandleFunc("/races/{id}/teams", handlers.GetRaceTeams).Methods("GET")
	api.HandleFunc("/races/{id}/results", handlers.EnterRaceResults).Methods("POST")

	// Standings
	api.HandleFunc("/standings", handlers.GetStandings).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
